package config

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/adapter"
)

func NewAdapters(keys []adapter.Key) adapter.ProcessedConfigsMap {
	adapters := make(adapter.ProcessedConfigsMap, len(keys))
	for _, key := range keys {
		// explicitly initialize with empty maps. nil maps are serialized to `null` in json, empty maps are serialized to `{}`
		adapters[key] = map[string]string{}
	}
	return adapters
}

type AdaptersBuilder struct {
	ConfigurationFetcher ConfigurationFetcher
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . ConfigurationFetcher

type ConfigurationFetcher interface {
	// Fetch is used get one profile per adapter key, if present
	Fetch(ctx context.Context, appID int64, adapterKeys []adapter.Key) (adapter.RawConfigsMap, error)
}

func (b *AdaptersBuilder) Build(ctx context.Context, appID int64, adapterKeys []adapter.Key) (adapter.ProcessedConfigsMap, error) {
	profiles, err := b.ConfigurationFetcher.Fetch(ctx, appID, adapterKeys)
	if err != nil {
		return nil, err
	}

	adapters := NewAdapters(adapterKeys)

	for key, profile := range profiles {
		extra := profile.AccountExtra
		appData := profile.AppData
		switch key {
		case adapter.ApplovinKey:
			adapters[key]["app_key"] = appData["app_key"]
		case adapter.BidmachineKey:
			adapters[key]["seller_id"] = extra["seller_id"]
			adapters[key]["endpoint"] = extra["endpoint"]
			adapters[key]["mediation_config"] = extra["mediation_config"]
		case adapter.BigoAdsKey:
			adapters[key]["app_id"] = appData["app_id"]
		case adapter.DTExchangeKey:
			adapters[key]["app_id"] = appData["app_id"]
		case adapter.MintegralKey:
			adapters[key]["app_id"] = appData["app_id"]
			adapters[key]["app_key"] = extra["app_key"]
		case adapter.UnityAdsKey:
			adapters[key]["game_id"] = appData["game_id"]
		default:
			adapters[key] = appData
		}
	}

	return adapters, nil
}
