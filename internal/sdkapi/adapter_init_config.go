package sdkapi

import (
	"fmt"

	"github.com/bidon-io/bidon-backend/internal/adapter"
)

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
		config = &BidmachineInitConfig{
			Placements: make(map[string]string),
		}
	case adapter.BigoAdsKey:
		config = new(BigoAdsInitConfig)
	case adapter.ChartboostKey:
		config = new(ChartboostInitConfig)
	case adapter.DTExchangeKey:
		config = new(DTExchangeInitConfig)
	case adapter.GAMKey:
		config = new(GAMInitConfig)
	case adapter.MetaKey:
		config = new(MetaInitConfig)
	case adapter.MintegralKey:
		config = new(MintegralInitConfig)
	case adapter.MolocoKey:
		config = new(MolocoInitConfig)
	case adapter.UnityAdsKey:
		config = new(UnityAdsInitConfig)
	case adapter.VKAdsKey:
		config = new(VKAdsInitConfig)
	case adapter.VungleKey:
		config = new(VungleInitConfig)
	case adapter.MobileFuseKey:
		config = new(MobileFuseInitConfig)
	case adapter.InmobiKey:
		config = new(InmobiInitConfig)
	case adapter.IronSourceKey:
		config = new(IronSourceInitConfig)
	case adapter.TaurusXKey:
		config = &TaurusXInitConfig{
			Placements: make([]string, 0),
		}
	case adapter.AmazonKey:
		config = new(AmazonInitConfig)
	case adapter.YandexKey:
		config = new(YandexInitConfig)
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
	SellerID        string            `json:"seller_id,omitempty"`
	Endpoint        string            `json:"endpoint,omitempty"`
	MediationConfig []string          `json:"mediation_config,omitempty"`
	Placements      map[string]string `json:"placements,omitempty"`
	Order           int               `json:"order"`
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

type ChartboostInitConfig struct {
	AppID        string `json:"app_id,omitempty"`
	AppSignature string `json:"app_signature,omitempty"`
	Order        int    `json:"order"`
}

func (a *ChartboostInitConfig) Key() adapter.Key {
	return adapter.ChartboostKey
}

func (a *ChartboostInitConfig) SetDefaultOrder() {
	a.Order = 2
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
	AppID string `json:"app_id,omitempty"`
	Order int    `json:"order"`
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

type MolocoInitConfig struct {
	AppKey string `json:"app_key,omitempty"`
	Order  int    `json:"order"`
}

func (a *MolocoInitConfig) Key() adapter.Key {
	return adapter.MolocoKey
}

func (a *MolocoInitConfig) SetDefaultOrder() {
	a.Order = 0
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

type VKAdsInitConfig struct {
	AppID string `json:"app_id,omitempty"`
	Order int    `json:"order"`
}

func (a *VKAdsInitConfig) Key() adapter.Key {
	return adapter.VKAdsKey
}

func (a *VKAdsInitConfig) SetDefaultOrder() {
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

type YandexInitConfig struct {
	MetricaID string `json:"metrica_id,omitempty"`
	Order     int    `json:"order"`
}

func (a *YandexInitConfig) Key() adapter.Key {
	return adapter.YandexKey
}

func (a *YandexInitConfig) SetDefaultOrder() {
	a.Order = 2
}

type IronSourceInitConfig struct {
	AppKey string `json:"app_key,omitempty"`
	Order  int    `json:"order"`
}

func (a *IronSourceInitConfig) Key() adapter.Key {
	return adapter.IronSourceKey
}

func (a *IronSourceInitConfig) SetDefaultOrder() {
	a.Order = 2
}

type TaurusXInitConfig struct {
	AppID      string   `json:"app_id,omitempty"`
	Placements []string `json:"placements,omitempty"`
	Order      int      `json:"order"`
}

func (a *TaurusXInitConfig) Key() adapter.Key {
	return adapter.TaurusXKey
}

func (a *TaurusXInitConfig) SetDefaultOrder() {
	a.Order = 0
}
