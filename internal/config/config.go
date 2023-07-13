package config

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/adapter"
)

type Adapters map[adapter.Key]map[string]any

func NewAdapters(keys []adapter.Key) Adapters {
	adapters := make(Adapters, len(keys))
	for _, key := range keys {
		// explicitly initialize with empty maps. nil maps are serialized to `null` in json, empty maps are serialized to `{}`
		adapters[key] = map[string]any{}
	}
	return adapters
}

type AdaptersBuilder struct {
	AppDemandProfileFetcher AppDemandProfileFetcher
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . AppDemandProfileFetcher

type AppDemandProfileFetcher interface {
	// Fetch is used get one profile per adapter key, if present
	Fetch(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]AppDemandProfile, error)
}

type AppDemandProfile struct {
	AdapterKey   adapter.Key
	AccountExtra map[string]any
}

func (b *AdaptersBuilder) Build(ctx context.Context, appID int64, adapterKeys []adapter.Key) (Adapters, error) {
	profiles, err := b.AppDemandProfileFetcher.Fetch(ctx, appID, adapterKeys)
	if err != nil {
		return nil, err
	}

	adapters := NewAdapters(adapterKeys)

	for _, profile := range profiles {
		key := profile.AdapterKey
		extra := profile.AccountExtra
		switch key {
		case adapter.ApplovinKey:
			adapters[key]["app_key"] = extra["api_key"] // notice the "app" and "api" difference
		case adapter.BidmachineKey:
			adapters[key]["seller_id"] = extra["seller_id"]
			adapters[key]["endpoint"] = extra["endpoint"]
			adapters[key]["mediation_config"] = extra["mediation_config"]
		default:
			adapters[key] = extra
		}
	}

	return adapters, nil
}
