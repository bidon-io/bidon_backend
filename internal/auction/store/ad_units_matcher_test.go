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
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shopspring/decimal"
)

func TestAdUnitsMatcher_Match(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	apps := make([]db.App, 2)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx)
	}

	applovinDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.ApplovinKey)
	})
	applovinAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = applovinDemandSource
	})

	bidmachineDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.BidmachineKey)
	})
	bidmachineAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = bidmachineDemandSource
	})

	lineItems := []db.LineItem{
		{
			AppID:  apps[0].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.BannerFormat),
				Valid:  true,
			},
			HumanName: "applovin-banner-banner",
			Code:      ptr("applovin-banner-banner"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.1")),
			AccountID: applovinAccount.ID,
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547776,
				Valid: true,
			},
			Extra: map[string]any{
				"placement_id": "applovin-banner-banner",
			},
		},
		{
			AppID:  apps[0].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.AdaptiveFormat),
				Valid:  true,
			},
			HumanName: "applovin-banner-adaptive",
			Code:      ptr("applovin-banner-adaptive"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.2")),
			AccountID: applovinAccount.ID,
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547777,
				Valid: true,
			},
			Extra: map[string]any{
				"placement_id": "applovin-banner-adaptive",
			},
		},
		{
			AppID:  apps[0].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.LeaderboardFormat),
				Valid:  true,
			},
			HumanName: "applovin-banner-leaderboard",
			Code:      ptr("applovin-banner-leaderboard"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.3")),
			AccountID: applovinAccount.ID,
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547778,
				Valid: true,
			},
			Extra: map[string]any{
				"placement_id": "applovin-banner-leaderboard",
			},
		},
		{
			AppID:     apps[0].ID,
			AdType:    db.InterstitialAdType,
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.3")),
			HumanName: "applovin-interstitial",
			Code:      ptr("applovin-interstitial"),
			AccountID: applovinAccount.ID,
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547779,
				Valid: true,
			},
			Extra: map[string]any{
				"placement_id": "applovin-interstitial",
			},
		},
		{
			AppID:     apps[0].ID,
			AdType:    db.InterstitialAdType,
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.0")),
			HumanName: "bidmachine-interstitial",
			Code:      ptr("bidmachine-interstitial"),
			AccountID: bidmachineAccount.ID,
			IsBidding: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547780,
				Valid: true,
			},
		},
		{
			AppID:  apps[1].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.MRECFormat),
				Valid:  true,
			},
			HumanName: "app2-applovin-banner-mrec",
			Code:      ptr("app2-applovin-banner-mrec"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.4")),
			AccountID: applovinAccount.ID,
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547781,
				Valid: true,
			},
			Extra: map[string]any{
				"placement_id": "app2-applovin-banner-mrec",
			},
		},
		{
			AppID:  apps[1].ID,
			AdType: db.BannerAdType,
			Format: sql.NullString{
				String: string(ad.MRECFormat),
				Valid:  true,
			},
			HumanName: "app2-bidmachine-banner-mrec",
			Code:      ptr("app2-bidmachine-banner-mrec"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.4")),
			AccountID: bidmachineAccount.ID,
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547782,
				Valid: true,
			},
		},
	}
	if err := tx.Create(&lineItems).Error; err != nil {
		t.Fatalf("Error creating line items (%+v): %v", lineItems, err)
	}

	matcher := store.AdUnitsMatcher{DB: tx}
	pf := 0.15

	testCases := []struct {
		params *auction.BuildParams
		want   []auction.AdUnit
	}{
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.EmptyFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
			},
			want: []auction.AdUnit{},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.AdaptiveFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
			},
			want: []auction.AdUnit{
				{
					DemandID:   "applovin",
					UID:        "1701972528521547776",
					PriceFloor: ptr(0.1),
					Label:      "applovin-banner-banner",
					BidType:    schema.CPMBidType,
					Extra: map[string]any{
						"placement_id": "applovin-banner-banner",
					},
				},
				{
					DemandID:   "applovin",
					UID:        "1701972528521547777",
					PriceFloor: ptr(0.2),
					Label:      "applovin-banner-adaptive",
					BidType:    schema.CPMBidType,
					Extra: map[string]any{
						"placement_id": "applovin-banner-adaptive",
					},
				},
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
			want: []auction.AdUnit{
				{
					DemandID:   "applovin",
					UID:        "1701972528521547777",
					PriceFloor: ptr(0.2),
					Label:      "applovin-banner-adaptive",
					BidType:    schema.CPMBidType,
					Extra: map[string]any{
						"placement_id": "applovin-banner-adaptive",
					},
				},
				{
					DemandID:   "applovin",
					UID:        "1701972528521547778",
					PriceFloor: ptr(0.3),
					Label:      "applovin-banner-leaderboard",
					BidType:    schema.CPMBidType,
					Extra: map[string]any{
						"placement_id": "applovin-banner-leaderboard",
					},
				},
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
			want: []auction.AdUnit{
				{
					DemandID:   "applovin",
					UID:        "1701972528521547777",
					PriceFloor: ptr(0.2),
					Label:      "applovin-banner-adaptive",
					BidType:    schema.CPMBidType,
					Extra: map[string]any{
						"placement_id": "applovin-banner-adaptive",
					},
				},
			},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.InterstitialType,
				AdFormat:   ad.EmptyFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey},
				PriceFloor: ptr(0.2),
			},
			want: []auction.AdUnit{
				{
					DemandID:   "applovin",
					UID:        "1701972528521547779",
					PriceFloor: ptr(0.3),
					Label:      "applovin-interstitial",
					BidType:    schema.CPMBidType,
					Extra: map[string]any{
						"placement_id": "applovin-interstitial",
					},
				},
				{
					DemandID:   "bidmachine",
					UID:        "1701972528521547780",
					PriceFloor: nil,
					Label:      "bidmachine-interstitial",
					BidType:    schema.RTBBidType,
					Extra:      map[string]any{},
				},
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
			want: []auction.AdUnit{
				{
					DemandID:   "applovin",
					UID:        "1701972528521547781",
					PriceFloor: ptr(0.4),
					Label:      "app2-applovin-banner-mrec",
					BidType:    schema.CPMBidType,
					Extra: map[string]any{
						"placement_id": "app2-applovin-banner-mrec",
					},
				},
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
			want: []auction.AdUnit{
				{
					DemandID:   "applovin",
					UID:        "1701972528521547777",
					PriceFloor: ptr(0.2),
					Label:      "applovin-banner-adaptive",
					BidType:    schema.CPMBidType,
					Extra: map[string]any{
						"placement_id": "applovin-banner-adaptive",
					},
				},
			},
		},
	}

	for _, tC := range testCases {
		got, err := matcher.Match(context.Background(), tC.params)
		if err != nil {
			t.Errorf("Error matching line items: %v", err)
		}

		less := func(a, b auction.AdUnit) bool { return a.Label < b.Label }
		if diff := cmp.Diff(tC.want, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("matcher.Match(ctx, %+v) mismatch (-want, +got)\n%s", tC.params, diff)
		}
	}
}
