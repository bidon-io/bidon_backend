package store_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/bidon-io/bidon-backend/internal/device"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shopspring/decimal"
)

func TestLineItemsMatcher_Match(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	apps := dbtest.CreateAppsList(t, tx, 2)

	applovinDemandSource := dbtest.CreateDemandSource(t, tx, dbtest.WithDemandSourceOptions(&db.DemandSource{
		APIKey: "applovin",
	}))
	applovinAccount := dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(&db.DemandSourceAccount{
		DemandSourceID: applovinDemandSource.ID,
		DemandSource:   *applovinDemandSource,
	}))

	bidmachineDemandSource := dbtest.CreateDemandSource(t, tx, dbtest.WithDemandSourceOptions(&db.DemandSource{
		APIKey: "bidmachine",
	}))
	bidmachineAccount := dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(&db.DemandSourceAccount{
		DemandSourceID: bidmachineDemandSource.ID,
		DemandSource:   *bidmachineDemandSource,
	}))

	lineItems := []db.LineItem{
		{
			AppID:  apps[0].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.BannerFormat),
				Valid:  true,
			},
			Code:      ptr("applovin-banner-banner"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.1")),
			AccountID: applovinAccount.ID,
		},
		{
			AppID:  apps[0].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.AdaptiveFormat),
				Valid:  true,
			},
			Code:      ptr("applovin-banner-adaptive"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.2")),
			AccountID: applovinAccount.ID,
		},
		{
			AppID:  apps[0].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.LeaderboardFormat),
				Valid:  true,
			},
			Code:      ptr("applovin-banner-leaderboard"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.3")),
			AccountID: applovinAccount.ID,
		},
		{
			AppID:     apps[0].ID,
			AdType:    db.InterstitialAdType,
			Code:      ptr("applovin-interstitial"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.3")),
			AccountID: applovinAccount.ID,
		},
		{
			AppID:     apps[0].ID,
			AdType:    db.InterstitialAdType,
			Code:      ptr("bidmachine-interstitial"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.3")),
			AccountID: bidmachineAccount.ID,
		},
		{
			AppID:  apps[1].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.MRECFormat),
				Valid:  true,
			},
			Code:      ptr("app2-applovin-banner-mrec"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.4")),
			AccountID: applovinAccount.ID,
		},
		{
			AppID:  apps[1].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.MRECFormat),
				Valid:  true,
			},
			Code:      ptr("app2-bidmachine-banner-mrec"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.4")),
			AccountID: bidmachineAccount.ID,
		},
	}
	if err := tx.Create(&lineItems).Error; err != nil {
		t.Fatalf("Error creating line items (%+v): %v", lineItems, err)
	}

	matcher := store.LineItemsMatcher{DB: tx}
	pf := 0.15

	testCases := []struct {
		params *auction.BuildParams
		want   []auction.LineItem
	}{
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.EmptyFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
			},
			want: []auction.LineItem{},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.AdaptiveFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
			},
			want: []auction.LineItem{
				{ID: "applovin", PriceFloor: 0.1, AdUnitID: "applovin-banner-banner"},
				{ID: "applovin", PriceFloor: 0.2, AdUnitID: "applovin-banner-adaptive"},
			},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.AdaptiveFormat,
				DeviceType: device.TabletType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
			},
			want: []auction.LineItem{
				{ID: "applovin", PriceFloor: 0.2, AdUnitID: "applovin-banner-adaptive"},
				{ID: "applovin", PriceFloor: 0.3, AdUnitID: "applovin-banner-leaderboard"},
			},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.AdaptiveFormat,
				DeviceType: device.UnknownType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
			},
			want: []auction.LineItem{
				{ID: "applovin", PriceFloor: 0.2, AdUnitID: "applovin-banner-adaptive"},
			},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.InterstitialType,
				AdFormat:   ad.EmptyFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey},
			},
			want: []auction.LineItem{
				{ID: "applovin", PriceFloor: 0.3, AdUnitID: "applovin-interstitial"},
				{ID: "bidmachine", PriceFloor: 0.3, AdUnitID: "bidmachine-interstitial"},
			},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[1].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.MRECFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
			},
			want: []auction.LineItem{
				{ID: "applovin", PriceFloor: 0.4, AdUnitID: "app2-applovin-banner-mrec"},
			},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.AdaptiveFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
				PriceFloor: &pf,
			},
			want: []auction.LineItem{
				{ID: "applovin", PriceFloor: 0.2, AdUnitID: "applovin-banner-adaptive"},
			},
		},
	}

	for _, tC := range testCases {
		got, err := matcher.Match(context.Background(), tC.params)
		if err != nil {
			t.Errorf("Error matching line items: %v", err)
		}

		less := func(a, b auction.LineItem) bool { return a.AdUnitID < b.AdUnitID }
		if diff := cmp.Diff(tC.want, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("matcher.Match(ctx, %+v) mismatch (-want, +got)\n%s", tC.params, diff)
		}
	}
}
