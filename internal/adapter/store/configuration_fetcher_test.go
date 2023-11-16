package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/adapter/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/google/go-cmp/cmp"
)

func TestAppDemandProfileFetcher_Fetch(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

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
			name:        "One key, App 2",
			appID:       apps[1].ID,
			adapterKeys: []adapter.Key{adapter.DTExchangeKey},
			want: adapter.RawConfigsMap{
				adapter.DTExchangeKey: {
					AccountExtra: map[string]any{"dtexchange": "dtexchange"},
					AppData:      map[string]any{},
				},
			},
		},
	}

	fetcher := store.ConfigurationFetcher{DB: tx}

	for _, tC := range testCases {
		got, err := fetcher.Fetch(context.Background(), tC.appID, tC.adapterKeys)
		if err != nil {
			t.Fatalf("failed to fetch app demand profiles: %v", err)
		}

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("fetcher.Fetch -> %v mismatch (-want +got):\n%s", tC.name, diff)
		}
	}
}
