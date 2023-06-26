package config_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/config"
	"github.com/google/go-cmp/cmp"
)

func TestAdaptersBuilder_Build(t *testing.T) {
	profiles := []config.AppDemandProfile{
		{adapter.ApplovinKey, map[string]any{"api_key": "applovin_api_key", "ext": "some"}},
		{adapter.BidmachineKey, map[string]any{
			"seller_id":        "bidmachine_seller_id",
			"endpoint":         "http://example.com/bidmachine",
			"mediation_config": "{\"config\": true}",
			"ext":              "some",
		}},
		{adapter.DTExchangeKey, map[string]any{"dt_key_1": 1, "dt_key_2": 2}},
		{adapter.UnityAdsKey, map[string]any{"unity_key_1": 1, "unity_key_2": 2}},
	}

	applovinProfile := profiles[0]
	bidmachineProfile := profiles[1]
	dtExchangeProfile := profiles[2]
	unityAdsProfile := profiles[3]

	testCases := []struct {
		name        string
		profiles    []config.AppDemandProfile
		adapterKeys []adapter.Key
		want        config.Adapters
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
			want:        config.Adapters{},
		},
		{
			name:        "All keys match all profiles",
			profiles:    profiles,
			adapterKeys: []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey, adapter.DTExchangeKey, adapter.UnityAdsKey},
			want: config.Adapters{
				adapter.ApplovinKey: map[string]any{
					"app_key": applovinProfile.AccountExtra["api_key"],
				},
				adapter.BidmachineKey: map[string]any{
					"seller_id":        bidmachineProfile.AccountExtra["seller_id"],
					"endpoint":         bidmachineProfile.AccountExtra["endpoint"],
					"mediation_config": bidmachineProfile.AccountExtra["mediation_config"],
				},
				adapter.DTExchangeKey: dtExchangeProfile.AccountExtra,
				adapter.UnityAdsKey:   unityAdsProfile.AccountExtra,
			},
		},
		{
			name:        "Some keys do not have matching profile",
			profiles:    []config.AppDemandProfile{dtExchangeProfile, unityAdsProfile},
			adapterKeys: []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey, adapter.DTExchangeKey, adapter.UnityAdsKey},
			want: config.Adapters{
				adapter.ApplovinKey:   map[string]any{},
				adapter.BidmachineKey: map[string]any{},
				adapter.DTExchangeKey: dtExchangeProfile.AccountExtra,
				adapter.UnityAdsKey:   unityAdsProfile.AccountExtra,
			},
		},
	}

	for _, tC := range testCases {
		fetcher := &config.AppDemandProfileFetcherMock{
			FetchFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]config.AppDemandProfile, error) {
				return tC.profiles, nil
			},
		}
		builder := &config.AdaptersBuilder{AppDemandProfileFetcher: fetcher}

		got, err := builder.Build(context.Background(), 0, tC.adapterKeys)
		if err != nil {
			t.Errorf("builder.Build -> %v: %v", tC.name, err)
		}

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("builder.Build -> %v mismatch (-want,+got)\n%s", tC.name, diff)
		}
	}
}
