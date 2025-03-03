package adapters_builder

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/bidmachine"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/bigoads"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/meta"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/mintegral"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/mobilefuse"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/vkads"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/vungle"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

var biddingAdapters = map[adapter.Key]adapters.Builder{
	adapter.BidmachineKey: bidmachine.Builder,
	adapter.BigoAdsKey:    bigoads.Builder,
	adapter.MetaKey:       meta.Builder,
	adapter.MintegralKey:  mintegral.Builder,
	adapter.MobileFuseKey: mobilefuse.Builder,
	adapter.VKAdsKey:      vkads.Builder,
	adapter.VungleKey:     vungle.Builder,
	// adapter.AdmobKey: admob.Builder,
	// adapter.ApplovinKey: applovin.Builder,
	// adapter.DTExchangeKey: dtexchange.Builder,
	// adapter.UnityAdsKey: unityads.Builder,
}

type AdaptersBuilder struct {
	AdaptersMap map[adapter.Key]adapters.Builder
	Client      *http.Client
	Config      *config.DemandConfig
}

func (b AdaptersBuilder) Build(adapterKey adapter.Key, cfg adapter.ProcessedConfigsMap) (*adapters.Bidder, error) {
	if f, ok := b.AdaptersMap[adapterKey]; ok {
		return f(cfg, b.Client)
	}

	return nil, fmt.Errorf("adapter %s not found", adapterKey)
}

func BuildBiddingAdapters(client *http.Client) AdaptersBuilder {
	return AdaptersBuilder{
		AdaptersMap: biddingAdapters,
		Client:      client,
	}
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . ConfigurationFetcher

type ConfigurationFetcher interface {
	// FetchCached is used get one profile per adapter key, if present
	FetchCached(ctx context.Context, appID int64, adapterKeys []adapter.Key) (adapter.RawConfigsMap, error)
}

type AdaptersConfigBuilder struct {
	ConfigurationFetcher ConfigurationFetcher
	DemandConfig         *config.DemandConfig
}

func NewAdaptersConfigBuilder(fetcher ConfigurationFetcher, config *config.DemandConfig) *AdaptersConfigBuilder {
	return &AdaptersConfigBuilder{
		ConfigurationFetcher: fetcher,
		DemandConfig:         config,
	}
}

func NewAdapters(keys []adapter.Key) adapter.ProcessedConfigsMap {
	adaptersMap := make(adapter.ProcessedConfigsMap, len(keys))
	for _, key := range keys {
		// explicitly initialize with empty maps. nil maps are serialized to `null` in json, empty maps are serialized to `{}`
		adaptersMap[key] = map[string]any{}
	}
	return adaptersMap
}

func (b *AdaptersConfigBuilder) Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, adUnitsMap *auction.AdUnitsMap) (adapter.ProcessedConfigsMap, error) {
	profiles, err := b.ConfigurationFetcher.FetchCached(ctx, appID, adapterKeys)
	if err != nil {
		return nil, err
	}

	adaptersMap := NewAdapters(adapterKeys)

	for key, profile := range profiles {
		extra := profile.AccountExtra
		appData := profile.AppData
		switch key {
		case adapter.AmazonKey:
			adaptersMap[key]["price_points_map"] = extra["price_points_map"]
		case adapter.BidmachineKey:
			adaptersMap[key]["seller_id"] = extra["seller_id"]
			adaptersMap[key]["endpoint"] = extra["endpoint"]
			adaptersMap[key]["mediation_config"] = extra["mediation_config"]
		case adapter.BigoAdsKey:
			adaptersMap[key]["app_id"] = appData["app_id"]
			adaptersMap[key]["seller_id"] = extra["publisher_id"]

			adUnit, _ := adUnitsMap.First(key, schema.RTBBidType)
			if adUnit != nil {
				adaptersMap[key]["tag_id"] = adUnit.Extra["slot_id"]
				adaptersMap[key]["placement_id"] = adUnit.Extra["placement_id"]
			}
		case adapter.MintegralKey:
			adaptersMap[key]["app_id"] = appData["app_id"]
			adaptersMap[key]["seller_id"] = extra["publisher_id"]

			adUnit, _ := adUnitsMap.First(key, schema.RTBBidType)
			if adUnit != nil {
				adaptersMap[key]["tag_id"] = adUnit.Extra["unit_id"]
				adaptersMap[key]["placement_id"] = adUnit.Extra["placement_id"]
			}
		case adapter.VKAdsKey:
			adaptersMap[key]["app_id"] = appData["app_id"]

			adUnit, _ := adUnitsMap.First(key, schema.RTBBidType)
			if adUnit != nil {
				adaptersMap[key]["tag_id"] = adUnit.Extra["slot_id"]
			}
		case adapter.VungleKey:
			adaptersMap[key]["app_id"] = appData["app_id"]
			adaptersMap[key]["seller_id"] = extra["account_id"]

			adUnit, _ := adUnitsMap.First(key, schema.RTBBidType)
			if adUnit != nil {
				adaptersMap[key]["tag_id"] = adUnit.Extra["placement_id"]
			}
		case adapter.MetaKey:
			adaptersMap[key]["app_id"] = appData["app_id"]
			adaptersMap[key]["app_secret"] = b.DemandConfig.MetaAppSecret
			adaptersMap[key]["platform_id"] = b.DemandConfig.MetaPlatformID

			adUnit, _ := adUnitsMap.First(key, schema.RTBBidType)
			if adUnit != nil {
				adaptersMap[key]["tag_id"] = adUnit.Extra["placement_id"]
			}
		case adapter.MobileFuseKey:
			adUnit, _ := adUnitsMap.First(key, schema.RTBBidType)
			if adUnit != nil {
				adaptersMap[key]["tag_id"] = adUnit.Extra["placement_id"]
			}
		default:
			adaptersMap[key] = extra
		}
	}

	return adaptersMap, nil
}
