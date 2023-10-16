package sdkapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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
	FetchAdapterInitConfigs(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]AdapterInitConfig, error)
}

type ConfigResponse struct {
	Init       ConfigResponseInit `json:"init"`
	Placements []any              `json:"placements"`
	Token      string             `json:"token"`
	Segment    Segment            `json:"segment"`
}

type Segment struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
}

type ConfigResponseInit struct {
	TMax     int                               `json:"tmax"`
	Adapters map[adapter.Key]AdapterInitConfig `json:"adapters"`
}

func (h *ConfigHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()

	configEvent := event.NewConfig(&req.raw, req.geoData)
	h.EventLogger.Log(configEvent, func(err error) {
		logError(c, fmt.Errorf("log config event: %v", err))
	})

	segmentParams := &segment.Params{
		Country: req.countryCode(),
		Ext:     req.raw.Segment.Ext,
		AppID:   req.app.ID,
	}

	sgmnt := h.SegmentMatcher.Match(ctx, segmentParams)

	var segmentID string
	if sgmnt.ID != 0 {
		segmentID = strconv.Itoa(int(sgmnt.ID))
	}

	adapterInitConfigs, err := h.AdapterInitConfigsFetcher.FetchAdapterInitConfigs(ctx, req.app.ID, req.raw.Adapters.Keys())
	if err != nil {
		return err
	}
	if len(adapterInitConfigs) == 0 {
		return ErrNoAdaptersFound
	}

	adapters := make(map[adapter.Key]AdapterInitConfig, len(adapterInitConfigs))
	for _, cfg := range adapterInitConfigs {
		adapters[cfg.Key()] = cfg
	}

	resp := &ConfigResponse{
		Init: ConfigResponseInit{
			TMax:     5000,
			Adapters: adapters,
		},
		Placements: []any{},
		Token:      "{}",
		Segment:    Segment{ID: segmentID, UID: sgmnt.UID},
	}

	return c.JSON(http.StatusOK, resp)
}

type AdapterInitConfig interface {
	Key() adapter.Key
}

func NewAdapterInitConfig(key adapter.Key) (AdapterInitConfig, error) {
	switch key {
	case adapter.AdmobKey:
		return new(AdmobInitConfig), nil
	case adapter.ApplovinKey:
		return new(ApplovinInitConfig), nil
	case adapter.BidmachineKey:
		return new(BidmachineInitConfig), nil
	case adapter.BigoAdsKey:
		return new(BigoAdsInitConfig), nil
	case adapter.DTExchangeKey:
		return new(DTExchangeInitConfig), nil
	case adapter.MetaKey:
		return new(MetaInitConfig), nil
	case adapter.MintegralKey:
		return new(MintegralInitConfig), nil
	case adapter.UnityAdsKey:
		return new(UnityAdsInitConfig), nil
	case adapter.VungleKey:
		return new(VungleInitConfig), nil
	case adapter.MobileFuseKey:
		return new(MobileFuseInitConfig), nil
	case adapter.InmobiKey:
		return new(InmobiInitConfig), nil
	case adapter.AmazonKey:
		return new(AmazonInitConfig), nil
	default:
		return nil, fmt.Errorf("AdapterInitConfig for key %q not defined", key)
	}
}

type AdmobInitConfig struct {
	AppID string `json:"app_id,omitempty"`
}

func (a *AdmobInitConfig) Key() adapter.Key {
	return adapter.AdmobKey
}

type ApplovinInitConfig struct {
	// AppKey is deprecated, it must be the same as SDKKey
	AppKey string `json:"app_key,omitempty"`
	SDKKey string `json:"sdk_key,omitempty"`
}

func (a *ApplovinInitConfig) Key() adapter.Key {
	return adapter.ApplovinKey
}

type BidmachineInitConfig struct {
	SellerID        string   `json:"seller_id,omitempty"`
	Endpoint        string   `json:"endpoint,omitempty"`
	MediationConfig []string `json:"mediation_config,omitempty"`
}

func (a *BidmachineInitConfig) Key() adapter.Key {
	return adapter.BidmachineKey
}

type BigoAdsInitConfig struct {
	AppID string `json:"app_id,omitempty"`
}

func (a *BigoAdsInitConfig) Key() adapter.Key {
	return adapter.BigoAdsKey
}

type DTExchangeInitConfig struct {
	AppID string `json:"app_id,omitempty"`
}

func (a *DTExchangeInitConfig) Key() adapter.Key {
	return adapter.DTExchangeKey
}

type MetaInitConfig struct {
	AppID     string `json:"app_id,omitempty"`
	AppSecret string `json:"app_secret,omitempty"`
}

func (a *MetaInitConfig) Key() adapter.Key {
	return adapter.MetaKey
}

type MintegralInitConfig struct {
	AppID  string `json:"app_id,omitempty"`
	AppKey string `json:"app_key,omitempty"`
}

func (a *MintegralInitConfig) Key() adapter.Key {
	return adapter.MintegralKey
}

type UnityAdsInitConfig struct {
	GameID string `json:"game_id,omitempty"`
}

func (a *UnityAdsInitConfig) Key() adapter.Key {
	return adapter.UnityAdsKey
}

type VungleInitConfig struct {
	AppID string `json:"app_id,omitempty"`
}

func (a *VungleInitConfig) Key() adapter.Key {
	return adapter.VungleKey
}

type MobileFuseInitConfig struct {
	PublisherID string `json:"publisher_id,omitempty"`
	AppKey      string `json:"app_key,omitempty"`
}

func (a *MobileFuseInitConfig) Key() adapter.Key {
	return adapter.MobileFuseKey
}

type InmobiInitConfig struct {
	AccountID string `json:"account_id,omitempty"`
	AppKey    string `json:"app_key,omitempty"`
}

func (a *InmobiInitConfig) Key() adapter.Key {
	return adapter.InmobiKey
}

type AmazonSlot struct {
	SlotUUID string `json:"slot_uuid,omitempty"`
	Format   string `json:"format,omitempty"`
}

type AmazonInitConfig struct {
	AppKey string       `json:"app_key,omitempty"`
	Slots  []AmazonSlot `json:"slots"`
}

func (a *AmazonInitConfig) Key() adapter.Key {
	return adapter.AmazonKey
}
