package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/config"
	"github.com/bidon-io/bidon-backend/internal/config/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/google/go-cmp/cmp"
)

func TestAppDemandProfileFetcher_Fetch(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	user := dbtest.CreateUser(t, tx, 1)
	apps := make([]*db.App, 2)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx, i, user)
	}
	demandSources := dbtest.CreateDemandSourcesList(t, tx, 2)
	accountApplovin := dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(
		&db.DemandSourceAccount{
			UserID:         user.ID,
			DemandSourceID: demandSources[0].ID,
			DemandSource: db.DemandSource{
				APIKey: string(adapter.ApplovinKey),
			},
			Extra: map[string]any{"applovin": "applovin"},
		}))
	accountBidmachine := dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(
		&db.DemandSourceAccount{
			UserID:         user.ID,
			DemandSourceID: demandSources[0].ID,
			DemandSource: db.DemandSource{
				APIKey: string(adapter.BidmachineKey),
			},
			Extra: map[string]any{"bidmachine": "bidmachine"},
		}))
	accountDtexchange := dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(
		&db.DemandSourceAccount{
			UserID:         user.ID,
			DemandSourceID: demandSources[1].ID,
			DemandSource: db.DemandSource{
				APIKey: string(adapter.DTExchangeKey),
			},
			Extra: map[string]any{"dtexchange": "dtexchange"},
		}))
	accountUnity := dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(
		&db.DemandSourceAccount{
			UserID:         user.ID,
			DemandSourceID: demandSources[1].ID,
			DemandSource: db.DemandSource{
				APIKey: string(adapter.UnityAdsKey),
			},
			Extra: map[string]any{"unity": "unity"},
		}))
	profiles := []db.AppDemandProfile{
		{
			AppID:          apps[0].ID,
			AccountID:      accountApplovin.ID,
			DemandSourceID: demandSources[0].ID,
			Account:        *accountApplovin,
		},
		{
			AppID:          apps[0].ID,
			AccountID:      accountBidmachine.ID,
			DemandSourceID: demandSources[1].ID,
			Account:        *accountBidmachine,
		},
		{
			AppID:          apps[1].ID,
			AccountID:      accountDtexchange.ID,
			DemandSourceID: demandSources[0].ID,
			Account:        *accountDtexchange,
		},
		{
			AppID:          apps[1].ID,
			AccountID:      accountUnity.ID,
			DemandSourceID: demandSources[1].ID,
			Account:        *accountUnity,
		},
	}

	// Batch insert does not set AppDemandProfile.AccountID from created associations.
	// But when creating individually, it works, I don't know why
	for _, profile := range profiles {
		if err := tx.Create(&profile).Error; err != nil {
			t.Fatalf("failed to create test data: %v", err)
		}
	}

	testCases := []struct {
		name        string
		appID       int64
		adapterKeys []adapter.Key
		want        []config.AppDemandProfile
	}{
		{
			name:        "All keys, App 1",
			appID:       apps[0].ID,
			adapterKeys: adapter.Keys,
			want: []config.AppDemandProfile{
				{
					AdapterKey:   adapter.ApplovinKey,
					AccountExtra: map[string]any{"applovin": "applovin"},
				},
				{
					AdapterKey:   adapter.BidmachineKey,
					AccountExtra: map[string]any{"bidmachine": "bidmachine"},
				},
			},
		},
		{
			name:        "One key, App 1",
			appID:       apps[0].ID,
			adapterKeys: []adapter.Key{adapter.ApplovinKey},
			want: []config.AppDemandProfile{
				{
					AdapterKey:   adapter.ApplovinKey,
					AccountExtra: map[string]any{"applovin": "applovin"},
				},
			},
		},
		{
			name:        "No keys, App 1",
			appID:       apps[0].ID,
			adapterKeys: []adapter.Key{},
			want:        []config.AppDemandProfile{},
		},
		{
			name:        "One key, App 2",
			appID:       apps[1].ID,
			adapterKeys: []adapter.Key{adapter.DTExchangeKey},
			want: []config.AppDemandProfile{
				{
					AdapterKey:   adapter.DTExchangeKey,
					AccountExtra: map[string]any{"dtexchange": "dtexchange"},
				},
			},
		},
	}

	fetcher := store.AppDemandProfileFetcher{DB: tx}

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
