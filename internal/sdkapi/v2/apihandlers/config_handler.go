package apihandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
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

	adapterInitConfigs, err := h.AdapterInitConfigsFetcher.FetchAdapterInitConfigs(ctx, req.app.ID, req.raw.Adapters.Keys(), setAmazonSlots, setOrder)

	if err != nil {
		return err
	}
	if len(adapterInitConfigs) == 0 {
		return sdkapi.ErrNoAdaptersFound
	}

	isIOS := req.raw.Device.OS == "iOS" // For iOS devices we should skip Amazon adapter
	isCOPPA := false
	if req.raw.Regulations != nil {
		isCOPPA = req.raw.Regulations.COPPA
	}
	chardonnayHack := req.app.ID == 735400 || req.app.ID == 735401 || req.app.ID == 735402 || req.app.ID == 735456
	adapters := make(map[adapter.Key]sdkapi.AdapterInitConfig, len(adapterInitConfigs))

	var bidMachinePlacements map[string]string

	for _, cfg := range adapterInitConfigs {
		// Skip Amazon adapter for iOS devices with SDK version < 0.7.3
		if isIOS && sdkapi.VersionLessThan073Constraint.Check(sdkVersion) && cfg.Key() == adapter.AmazonKey {
			continue
		}
		if isCOPPA && adapter.IsDisabledForCOPPA(cfg.Key()) {
			continue
		}
		// TODO: Remove hacks after experiment
		if chardonnayHack && cfg.Key() != adapter.BidmachineKey {
			continue
		}

		if cfg.Key() == adapter.BidmachineKey {
			placements, err := h.ConfigFetcher.FetchBidMachinePlacements(ctx, req.app.ID)
			if err != nil {
				sdkapi.LogError(c, fmt.Errorf("fetch bidmachine placements: %v", err))
			} else {
				bidMachinePlacements = placements
			}
		}

		if cfg.Key() == adapter.BidmachineKey && bidMachinePlacements != nil {
			if bmCfg, ok := cfg.(*sdkapi.BidmachineInitConfig); ok {
				bmCfg.Placements = bidMachinePlacements
			}
		}

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
