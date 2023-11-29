package sdkapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
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
	FetchAdapterInitConfigs(ctx context.Context, appID int64, adapterKeys []adapter.Key, sdkVersion *semver.Version, setOrder bool) ([]AdapterInitConfig, error)
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
		logError(c, fmt.Errorf("log config event: %v", err))
	})

	sdkVersion, err := req.raw.GetSDKVersionSemver()
	if err != nil {
		return ErrInvalidSDKVersion
	}

	setOrder := req.raw.Device.OS == "android"
	adapterInitConfigs, err := h.AdapterInitConfigsFetcher.FetchAdapterInitConfigs(ctx, req.app.ID, req.raw.Adapters.Keys(), sdkVersion, setOrder)
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
			TMax:     10000,
			Adapters: adapters,
		},
		Placements: []any{},
		Token:      "{}",
		Segment:    Segment{ID: sgmnt.StringID(), UID: sgmnt.UID},
	}

	return c.JSON(http.StatusOK, resp)
}

func prepareConfigEvent(req *request[schema.ConfigRequest, *schema.ConfigRequest]) *event.RequestEvent {
	adRequestParams := event.AdRequestParams{
		EventType: "config",
	}
	return event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
}

type AdapterInitConfig interface {
	Key() adapter.Key
	SetDefaultOrder()
}

func NewAdapterInitConfig(key adapter.Key, setOrder bool) (AdapterInitConfig, error) {
	var config AdapterInitConfig
	switch key {
	case adapter.AdmobKey:
		config = new(AdmobInitConfig)
	case adapter.ApplovinKey:
		config = new(ApplovinInitConfig)
	case adapter.BidmachineKey:
		config = new(BidmachineInitConfig)
	case adapter.BigoAdsKey:
		config = new(BigoAdsInitConfig)
	case adapter.DTExchangeKey:
		config = new(DTExchangeInitConfig)
	case adapter.GAMKey:
		config = new(GAMInitConfig)
	case adapter.MetaKey:
		config = new(MetaInitConfig)
	case adapter.MintegralKey:
		config = new(MintegralInitConfig)
	case adapter.UnityAdsKey:
		config = new(UnityAdsInitConfig)
	case adapter.VungleKey:
		config = new(VungleInitConfig)
	case adapter.MobileFuseKey:
		config = new(MobileFuseInitConfig)
	case adapter.InmobiKey:
		config = new(InmobiInitConfig)
	case adapter.AmazonKey:
		config = new(AmazonInitConfig)
	default:
		return nil, fmt.Errorf("AdapterInitConfig for key %q not defined", key)
	}

	if setOrder {
		config.SetDefaultOrder()
	}

	return config, nil
}

type AdmobInitConfig struct {
	AppID string `json:"app_id,omitempty"`
	Order int    `json:"order"`
}

func (a *AdmobInitConfig) Key() adapter.Key {
	return adapter.AdmobKey
}

func (a *AdmobInitConfig) SetDefaultOrder() {
	a.Order = 1
}

type ApplovinInitConfig struct {
	// AppKey is deprecated, it must be the same as SDKKey
	AppKey string `json:"app_key,omitempty"`
	SDKKey string `json:"sdk_key,omitempty"`
	Order  int    `json:"order"`
}

func (a *ApplovinInitConfig) Key() adapter.Key {
	return adapter.ApplovinKey
}

func (a *ApplovinInitConfig) SetDefaultOrder() {
	a.Order = 1
}

type BidmachineInitConfig struct {
	SellerID        string   `json:"seller_id,omitempty"`
	Endpoint        string   `json:"endpoint,omitempty"`
	MediationConfig []string `json:"mediation_config,omitempty"`
	Order           int      `json:"order"`
}

func (a *BidmachineInitConfig) Key() adapter.Key {
	return adapter.BidmachineKey
}

func (a *BidmachineInitConfig) SetDefaultOrder() {
	a.Order = 0
}

type BigoAdsInitConfig struct {
	AppID string `json:"app_id,omitempty"`
	Order int    `json:"order"`
}

func (a *BigoAdsInitConfig) Key() adapter.Key {
	return adapter.BigoAdsKey
}

func (a *BigoAdsInitConfig) SetDefaultOrder() {
	a.Order = 0
}

type DTExchangeInitConfig struct {
	AppID string `json:"app_id,omitempty"`
	Order int    `json:"order"`
}

func (a *DTExchangeInitConfig) Key() adapter.Key {
	return adapter.DTExchangeKey
}

func (a *DTExchangeInitConfig) SetDefaultOrder() {
	a.Order = 0
}

type GAMInitConfig struct {
	NetworkCode string `json:"network_code,omitempty"`
	AppID       string `json:"app_id,omitempty"`
	Order       int    `json:"order"`
}

func (a *GAMInitConfig) Key() adapter.Key {
	return adapter.GAMKey
}

func (a *GAMInitConfig) SetDefaultOrder() {
	a.Order = 1
}

type MetaInitConfig struct {
	AppID     string `json:"app_id,omitempty"`
	AppSecret string `json:"app_secret,omitempty"`
	Order     int    `json:"order"`
}

func (a *MetaInitConfig) Key() adapter.Key {
	return adapter.MetaKey
}

func (a *MetaInitConfig) SetDefaultOrder() {
	a.Order = 0
}

type MintegralInitConfig struct {
	AppID  string `json:"app_id,omitempty"`
	AppKey string `json:"app_key,omitempty"`
	Order  int    `json:"order"`
}

func (a *MintegralInitConfig) Key() adapter.Key {
	return adapter.MintegralKey
}

func (a *MintegralInitConfig) SetDefaultOrder() {
	a.Order = 3
}

type UnityAdsInitConfig struct {
	GameID string `json:"game_id,omitempty"`
	Order  int    `json:"order"`
}

func (a *UnityAdsInitConfig) Key() adapter.Key {
	return adapter.UnityAdsKey
}

func (a *UnityAdsInitConfig) SetDefaultOrder() {
	a.Order = 2
}

type VungleInitConfig struct {
	AppID string `json:"app_id,omitempty"`
	Order int    `json:"order"`
}

func (a *VungleInitConfig) Key() adapter.Key {
	return adapter.VungleKey
}

func (a *VungleInitConfig) SetDefaultOrder() {
	a.Order = 2
}

type MobileFuseInitConfig struct {
	PublisherID string `json:"publisher_id,omitempty"`
	AppKey      string `json:"app_key,omitempty"`
	Order       int    `json:"order"`
}

func (a *MobileFuseInitConfig) Key() adapter.Key {
	return adapter.MobileFuseKey
}

func (a *MobileFuseInitConfig) SetDefaultOrder() {
	a.Order = 3
}

type InmobiInitConfig struct {
	AccountID string `json:"account_id,omitempty"`
	AppKey    string `json:"app_key,omitempty"`
	Order     int    `json:"order"`
}

func (a *InmobiInitConfig) Key() adapter.Key {
	return adapter.InmobiKey
}

func (a *InmobiInitConfig) SetDefaultOrder() {
	a.Order = 3
}

// Deprecated in 0.5.0
type AmazonSlot struct {
	SlotUUID string `json:"slot_uuid,omitempty"`
	Format   string `json:"format,omitempty"`
}

type AmazonInitConfig struct {
	AppKey string       `json:"app_key,omitempty"`
	Slots  []AmazonSlot `json:"slots"` // Deprecated in 0.5.0
	Order  int          `json:"order"`
}

func (a *AmazonInitConfig) Key() adapter.Key {
	return adapter.AmazonKey
}

func (a *AmazonInitConfig) SetDefaultOrder() {
	a.Order = 0
}
