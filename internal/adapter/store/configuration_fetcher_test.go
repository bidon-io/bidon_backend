package store_test

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/config"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/adapter/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/google/go-cmp/cmp"
)

func TestAppDemandProfileFetcher_FetchCached(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	rdbCache, _ := redismock.NewClusterMock()

	user := dbtest.CreateUser(t, tx)

	apps := make([]db.App, 2)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx, func(app *db.App) {
			app.User = user
		})
	}

	keys := []adapter.Key{
		adapter.ApplovinKey,
		adapter.UnityAdsKey,
		adapter.BidmachineKey,
		adapter.DTExchangeKey,
		adapter.AmazonKey,
		adapter.GAMKey,
		adapter.YandexKey,
	}

	demandSources := make([]db.DemandSource, len(keys))
	for i, key := range keys {
		demandSources[i] = dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
			source.APIKey = string(key)
		})
	}

	accounts := make([]db.DemandSourceAccount, len(demandSources))
	for i, source := range demandSources {
		accounts[i] = dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
			account.User = user
			account.DemandSource = source
			switch adapter.Key(source.APIKey) {
			case adapter.ApplovinKey:
				account.Extra = []byte(`{"applovin": "applovin"}`)
			case adapter.UnityAdsKey:
				account.Extra = []byte(`{"unity": "unity"}`)
			case adapter.BidmachineKey:
				account.Extra = []byte(`{"bidmachine": "bidmachine"}`)
			case adapter.DTExchangeKey:
				account.Extra = []byte(`{"dtexchange": "dtexchange"}`)
			case adapter.AmazonKey:
				account.Extra = []byte(`{"amazon": "amazon", "price_points": [{ "name": "name", "price_point": "price_point", "price": 1.0 }]}`)
			case adapter.GAMKey:
				account.Extra = []byte(`{"network_code": "111"}`)
			case adapter.YandexKey:
				account.Extra = []byte(`{"oauth_token": "yandex"}`)
			default:
				account.Extra = []byte{}
			}
		})
	}

	for i, account := range accounts {
		dbtest.CreateAppDemandProfile(t, tx, func(profile *db.AppDemandProfile) {
			profile.App = apps[i%len(apps)]
			profile.Account = account
			profile.Data = []byte(`{}`)
		})
	}

	testCases := []struct {
		name        string
		appID       int64
		adapterKeys []adapter.Key
		want        adapter.RawConfigsMap
	}{
		{
			name:        "All keys, App 1",
			appID:       apps[0].ID,
			adapterKeys: adapter.Keys,
			want: adapter.RawConfigsMap{
				adapter.ApplovinKey: {
					AccountExtra: map[string]any{"applovin": "applovin"},
					AppData:      map[string]any{},
				},
				adapter.BidmachineKey: {
					AccountExtra: map[string]any{"bidmachine": "bidmachine"},
					AppData:      map[string]any{},
				},
				adapter.AmazonKey: {
					AccountExtra: map[string]any{"amazon": "amazon", "price_points": []any{
						map[string]any{"name": "name", "price_point": "price_point", "price": 1.0},
					}},
					AppData: map[string]any{},
				},
				adapter.YandexKey: {
					AccountExtra: map[string]any{"oauth_token": "yandex"},
					AppData:      map[string]any{},
				},
			},
		},
		{
			name:        "One key, App 1",
			appID:       apps[0].ID,
			adapterKeys: []adapter.Key{adapter.ApplovinKey},
			want: adapter.RawConfigsMap{
				adapter.ApplovinKey: {
					AccountExtra: map[string]any{"applovin": "applovin"},
					AppData:      map[string]any{},
				},
			},
		},
		{
			name:        "No keys, App 1",
			appID:       apps[0].ID,
			adapterKeys: []adapter.Key{},
			want:        adapter.RawConfigsMap{},
		},
		{
			name:        "All keys, App 2",
			appID:       apps[1].ID,
			adapterKeys: adapter.Keys,
			want: adapter.RawConfigsMap{
				adapter.DTExchangeKey: {
					AccountExtra: map[string]any{"dtexchange": "dtexchange"},
					AppData:      map[string]any{},
				},
				adapter.GAMKey: {
					AccountExtra: map[string]any{"network_code": "111"},
					AppData:      map[string]any{},
				},
				adapter.UnityAdsKey: {
					AccountExtra: map[string]any{"unity": "unity"},
					AppData:      map[string]any{},
				},
			},
		},
	}

	configsCache := config.NewRedisCacheOf[adapter.RawConfigsMap](rdbCache, 10*time.Minute, "configs")
	fetcher := store.ConfigurationFetcher{
		DB:    tx,
		Cache: configsCache,
	}

	for _, tC := range testCases {
		got, err := fetcher.FetchCached(context.Background(), tC.appID, tC.adapterKeys)
		if err != nil {
			t.Fatalf("failed to fetch app demand profiles: %v", err)
		}

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("fetcher.Fetch -> %v mismatch (-want +got):\n%s", tC.name, diff)
		}
	}
}
