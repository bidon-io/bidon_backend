package config_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/config"
	configmocks "github.com/bidon-io/bidon-backend/internal/config/mocks"
	"github.com/google/go-cmp/cmp"
)

func TestAdaptersBuilder_Build(t *testing.T) {
	profiles := adapter.RawConfigsMap{
		adapter.ApplovinKey: {
			AppData: map[string]any{"api_key": "applovin_api_key", "ext": "some"},
		},
		adapter.BidmachineKey: {
			AccountExtra: map[string]any{
				"seller_id":        "bidmachine_seller_id",
				"endpoint":         "http://example.com/bidmachine",
				"mediation_config": "{\"config\": true}",
				"ext":              "some",
			},
		},
		adapter.DTExchangeKey: {
			AppData: map[string]any{"app_id": "123", "dt_key_1": "1", "dt_key_2": "2"},
		},
		adapter.UnityAdsKey: {
			AppData: map[string]any{"game_id": "234", "unity_key_1": "1", "unity_key_2": "2"},
		},
	}

	applovinProfile := profiles[adapter.ApplovinKey]
	bidmachineProfile := profiles[adapter.BidmachineKey]
	dtExchangeProfile := profiles[adapter.DTExchangeKey]
	unityAdsProfile := profiles[adapter.UnityAdsKey]

	testCases := []struct {
		name        string
		profiles    adapter.RawConfigsMap
		adapterKeys []adapter.Key
		want        adapter.ProcessedConfigsMap
	}{
		{
			name:        "All keys, no profiles",
			profiles:    nil,
			adapterKeys: adapter.Keys,
			want:        config.NewAdapters(adapter.Keys),
		},
		{
			name:        "No keys",
			profiles:    nil,
			adapterKeys: []adapter.Key{},
			want:        adapter.ProcessedConfigsMap{},
		},
		{
			name:        "All keys match all profiles",
			profiles:    profiles,
			adapterKeys: []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey, adapter.DTExchangeKey, adapter.UnityAdsKey},
			want: adapter.ProcessedConfigsMap{
				adapter.ApplovinKey: map[string]any{
					"app_key": applovinProfile.AccountExtra["api_key"],
				},
				adapter.BidmachineKey: map[string]any{
					"seller_id":        bidmachineProfile.AccountExtra["seller_id"],
					"endpoint":         bidmachineProfile.AccountExtra["endpoint"],
					"mediation_config": bidmachineProfile.AccountExtra["mediation_config"],
				},
				adapter.DTExchangeKey: map[string]any{"app_id": "123"},
				adapter.UnityAdsKey:   map[string]any{"game_id": "234"},
			},
		},
		{
			name: "Some keys do not have matching profile",
			profiles: adapter.RawConfigsMap{
				adapter.DTExchangeKey: dtExchangeProfile,
				adapter.UnityAdsKey:   unityAdsProfile,
			},
			adapterKeys: []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey, adapter.DTExchangeKey, adapter.UnityAdsKey},
			want: adapter.ProcessedConfigsMap{
				adapter.ApplovinKey:   map[string]any{},
				adapter.BidmachineKey: map[string]any{},
				adapter.DTExchangeKey: map[string]any{"app_id": "123"},
				adapter.UnityAdsKey:   map[string]any{"game_id": "234"},
			},
		},
	}

	for _, tC := range testCases {
		fetcher := &configmocks.ConfigurationFetcherMock{
			FetchFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key) (adapter.RawConfigsMap, error) {
				return tC.profiles, nil
			},
		}
		builder := &config.AdaptersBuilder{ConfigurationFetcher: fetcher}

		got, err := builder.Build(context.Background(), 0, tC.adapterKeys)
		if err != nil {
			t.Errorf("builder.Build -> %v: %v", tC.name, err)
		}

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("builder.Build -> %v mismatch (-want,+got)\n%s", tC.name, diff)
		}
	}
}
