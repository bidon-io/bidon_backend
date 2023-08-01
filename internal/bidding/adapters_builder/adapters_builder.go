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
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/mintegral"
)

var biddingAdapters = map[adapter.Key]adapters.Builder{
	adapter.BidmachineKey: bidmachine.Builder,
	adapter.BigoAdsKey:    bigoads.Builder,
	adapter.MintegralKey:  mintegral.Builder,
	// adapter.AdmobKey: admob.Builder,
	// adapter.ApplovinKey: applovin.Builder,
	// adapter.DTExchangeKey: dtexchange.Builder,
	// adapter.MetaKey: meta.Builder,
	// adapter.MobileFuseKey: mobilefuse.Builder,
	// adapter.UnityAdsKey: unityads.Builder,
	// adapter.VungleKey: vungle.Builder,
}

type AdaptersBuilder struct {
	AdaptersMap map[adapter.Key]adapters.Builder
	Client      *http.Client
}

func (b AdaptersBuilder) Build(adapterKey adapter.Key, cfg adapter.ProcessedConfigsMap) (adapters.Bidder, error) {
	if f, ok := b.AdaptersMap[adapterKey]; ok {
		return f(cfg, b.Client)
	}
	empty := adapters.Bidder{}

	return empty, fmt.Errorf("adapter %s not found", adapterKey)
}

func BuildBiddingAdapters(client *http.Client) AdaptersBuilder {
	return AdaptersBuilder{
		AdaptersMap: biddingAdapters,
		Client:      client,
	}
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . ConfigurationFetcher LineItemsMatcher

type ConfigurationFetcher interface {
	// Fetch is used get one profile per adapter key, if present
	Fetch(ctx context.Context, appID int64, adapterKeys []adapter.Key) (adapter.RawConfigsMap, error)
}

type LineItemsMatcher interface {
	Match(ctx context.Context, params *auction.BuildParams) ([]auction.LineItem, error)
}

type AdaptersConfigBuilder struct {
	ConfigurationFetcher ConfigurationFetcher
	LineItemsMatcher     LineItemsMatcher
}

func NewAdapters(keys []adapter.Key) adapter.ProcessedConfigsMap {
	adapters := make(adapter.ProcessedConfigsMap, len(keys))
	for _, key := range keys {
		// explicitly initialize with empty maps. nil maps are serialized to `null` in json, empty maps are serialized to `{}`
		adapters[key] = map[string]string{}
	}
	return adapters
}

func (b *AdaptersConfigBuilder) Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp) (adapter.ProcessedConfigsMap, error) {
	profiles, err := b.ConfigurationFetcher.Fetch(ctx, appID, adapterKeys)
	if err != nil {
		return nil, err
	}
	lineItemsMap, err := b.buildLineItemsMap(ctx, appID, adapterKeys, imp)
	if err != nil {
		return nil, err
	}
	adapters := NewAdapters(adapterKeys)

	for key, profile := range profiles {
		extra := profile.AccountExtra
		appData := profile.AppData
		switch key {
		case adapter.BidmachineKey:
			adapters[key]["seller_id"] = extra["seller_id"]
			adapters[key]["endpoint"] = extra["endpoint"]
			adapters[key]["mediation_config"] = extra["mediation_config"]
		case adapter.BigoAdsKey:
			adapters[key]["app_id"] = appData["app_id"]
			adapters[key]["seller_id"] = extra["seller_id"]
			adapters[key]["tag_id"] = ""
			adapters[key]["placement_id"] = ""

			if lineItem, ok := lineItemsMap[key]; ok {
				adapters[key]["tag_id"] = lineItem.AdUnitID
				adapters[key]["placement_id"] = lineItem.AdUnitID
			}
		case adapter.MintegralKey:
			adapters[key]["app_id"] = appData["app_id"]
			adapters[key]["seller_id"] = extra["seller_id"]
			adapters[key]["tag_id"] = ""
			adapters[key]["placement_id"] = ""

			if lineItem, ok := lineItemsMap[key]; ok {
				adapters[key]["tag_id"] = lineItem.AdUnitID
				adapters[key]["placement_id"] = lineItem.PlacementID
			}
		default:
			adapters[key] = extra
		}
	}

	return adapters, nil
}

func (b *AdaptersConfigBuilder) buildLineItemsMap(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp) (map[adapter.Key]auction.LineItem, error) {
	lineItems, err := b.LineItemsMatcher.Match(ctx, &auction.BuildParams{
		Adapters: adapterKeys,
		AppID:    appID,
		AdType:   imp.Type(),
		AdFormat: imp.Format(),
	})
	if err != nil {
		return nil, err
	}

	// adapter key to lineItem map
	lineItemsMap := make(map[adapter.Key]auction.LineItem)
	for _, item := range lineItems {
		if _, exists := lineItemsMap[adapter.Key(item.ID)]; !exists {
			lineItemsMap[adapter.Key(item.ID)] = item
		}
	}

	return lineItemsMap, nil
}
