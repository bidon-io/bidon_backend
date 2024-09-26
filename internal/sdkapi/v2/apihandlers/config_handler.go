package apihandlers

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/sdkapi"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	*BaseHandler[schema.ConfigRequest, *schema.ConfigRequest]
	AdapterInitConfigsFetcher AdapterInitConfigsFetcher
	SegmentMatcher            *segment.Matcher
	EventLogger               *event.Logger
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/config_mocks.go -pkg mocks . AdapterInitConfigsFetcher
type AdapterInitConfigsFetcher interface {
	FetchAdapterInitConfigs(ctx context.Context, appID int64, adapterKeys []adapter.Key, setAmazonSlots bool, setOrder bool) ([]sdkapi.AdapterInitConfig, error)
}

type ConfigResponse struct {
	Init       ConfigResponseInit `json:"init"`
	Placements []any              `json:"placements"`
	Token      string             `json:"token"`
	Segment    Segment            `json:"segment"`
	Bidding    ConfigBidding      `json:"bidding"`
}

type Segment struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
}

type ConfigResponseInit struct {
	TMax     int                                      `json:"tmax"`
	Adapters map[adapter.Key]sdkapi.AdapterInitConfig `json:"adapters"`
}

type ConfigBidding struct {
	TokenTimeoutMS int `json:"token_timeout_ms"`
}

func (h *ConfigHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()

	segmentParams := &segment.Params{
		Country: req.countryCode(),
		Ext:     req.raw.Segment.Ext,
		AppID:   req.app.ID,
	}

	sgmnt := h.SegmentMatcher.Match(ctx, segmentParams)
	req.raw.Segment.ID = sgmnt.StringID()
	req.raw.Segment.UID = sgmnt.UID

	configRequestEvent := prepareConfigEvent(req)
	h.EventLogger.Log(configRequestEvent, func(err error) {
		sdkapi.LogError(c, fmt.Errorf("log config event: %v", err))
	})

	sdkVersion, err := req.raw.GetSDKVersionSemver()
	if err != nil {
		return sdkapi.ErrInvalidSDKVersion
	}

	setOrder := req.raw.Device.OS == "android"                      // Set order for Android devices only
	setAmazonSlots := !sdkapi.Version05Constraint.Check(sdkVersion) // Do not set Amazon slots for SDK version 0.5.x

	// TODO remove this hack after test is completed
	isMergeBlockIOS := req.app.ID == 735361
	if isMergeBlockIOS {
		constraint, _ := semver.NewConstraint("= 0.7.0-next.3 || = 0.7.0-next.4")
		setAmazonSlots = !constraint.Check(sdkVersion) // For all versions except 0.7.0-0-next.3 sets to false it used in ConfigsFetcher
	}

	adapterInitConfigs, err := h.AdapterInitConfigsFetcher.FetchAdapterInitConfigs(ctx, req.app.ID, req.raw.Adapters.Keys(), setAmazonSlots, setOrder)
	if err != nil {
		return err
	}
	if len(adapterInitConfigs) == 0 {
		return sdkapi.ErrNoAdaptersFound
	}

	adapters := make(map[adapter.Key]sdkapi.AdapterInitConfig, len(adapterInitConfigs))
	for _, cfg := range adapterInitConfigs {
		adapters[cfg.Key()] = cfg
	}

	resp := &ConfigResponse{
		Init: ConfigResponseInit{
			TMax:     10000,
			Adapters: adapters,
		},
		Placements: []any{},
		Token:      "{}",
		Segment:    Segment{ID: sgmnt.StringID(), UID: sgmnt.UID},
		Bidding:    ConfigBidding{TokenTimeoutMS: 10000},
	}

	return c.JSON(http.StatusOK, resp)
}

func prepareConfigEvent(req *request[schema.ConfigRequest, *schema.ConfigRequest]) *event.AdEvent {
	adRequestParams := event.AdRequestParams{
		EventType: "config",
	}

	return event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData)
}
