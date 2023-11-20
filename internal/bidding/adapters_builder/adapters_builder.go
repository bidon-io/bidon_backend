package adapters_builder

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/bidmachine"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/bigoads"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/meta"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/mintegral"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/mobilefuse"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/vungle"
)

var biddingAdapters = map[adapter.Key]adapters.Builder{
	adapter.BidmachineKey: bidmachine.Builder,
	adapter.BigoAdsKey:    bigoads.Builder,
	adapter.MetaKey:       meta.Builder,
	adapter.MintegralKey:  mintegral.Builder,
	adapter.MobileFuseKey: mobilefuse.Builder,
	adapter.VungleKey:     vungle.Builder,
	// adapter.AdmobKey: admob.Builder,
	// adapter.ApplovinKey: applovin.Builder,
	// adapter.DTExchangeKey: dtexchange.Builder,
	// adapter.UnityAdsKey: unityads.Builder,
}

type AdaptersBuilder struct {
	AdaptersMap map[adapter.Key]adapters.Builder
	Client      *http.Client
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
	// Fetch is used get one profile per adapter key, if present
	Fetch(ctx context.Context, appID int64, adapterKeys []adapter.Key) (adapter.RawConfigsMap, error)
}

type AdaptersConfigBuilder struct {
	ConfigurationFetcher ConfigurationFetcher
}

func NewAdapters(keys []adapter.Key) adapter.ProcessedConfigsMap {
	adapters := make(adapter.ProcessedConfigsMap, len(keys))
	for _, key := range keys {
		// explicitly initialize with empty maps. nil maps are serialized to `null` in json, empty maps are serialized to `{}`
		adapters[key] = map[string]any{}
	}
	return adapters
}

func (b *AdaptersConfigBuilder) Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp, adUnitsMap *map[adapter.Key][]auction.AdUnit) (adapter.ProcessedConfigsMap, error) {
	profiles, err := b.ConfigurationFetcher.Fetch(ctx, appID, adapterKeys)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	adapters := NewAdapters(adapterKeys)

	for key, profile := range profiles {
		extra := profile.AccountExtra
		appData := profile.AppData
		switch key {
		case adapter.AmazonKey:
			adapters[key]["price_points_map"] = extra["price_points_map"]
		case adapter.BidmachineKey:
			adapters[key]["seller_id"] = extra["seller_id"]
			adapters[key]["endpoint"] = extra["endpoint"]
			adapters[key]["mediation_config"] = extra["mediation_config"]
		case adapter.BigoAdsKey:
			adapters[key]["app_id"] = appData["app_id"]
			adapters[key]["seller_id"] = extra["publisher_id"]

			if adUnits, ok := (*adUnitsMap)[key]; ok {
				adapters[key]["tag_id"] = adUnits[0].Extra["slot_id"]
				adapters[key]["placement_id"] = adUnits[0].Extra["placement_id"]
			}
		case adapter.MintegralKey:
			adapters[key]["app_id"] = appData["app_id"]
			adapters[key]["seller_id"] = extra["publisher_id"]

			if adUnits, ok := (*adUnitsMap)[key]; ok {
				adapters[key]["tag_id"] = adUnits[0].Extra["unit_id"]
				adapters[key]["placement_id"] = adUnits[0].Extra["placement_id"]
			}
		case adapter.VungleKey:
			adapters[key]["app_id"] = appData["app_id"]
			adapters[key]["seller_id"] = extra["account_id"]

			if adUnits, ok := (*adUnitsMap)[key]; ok {
				adapters[key]["tag_id"] = adUnits[0].Extra["placement_id"]
			}
		case adapter.MetaKey:
			adapters[key]["app_id"] = appData["app_id"]
			adapters[key]["app_secret"] = appData["app_secret"]

			if adUnits, ok := (*adUnitsMap)[key]; ok {
				adapters[key]["tag_id"] = adUnits[0].Extra["placement_id"]
			}
		case adapter.MobileFuseKey:
			if adUnits, ok := (*adUnitsMap)[key]; ok {
				adapters[key]["tag_id"] = adUnits[0].Extra["placement_id"]
			}
		default:
			adapters[key] = extra
		}
	}

	return adapters, nil
}
