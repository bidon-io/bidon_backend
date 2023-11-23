package store_test

import (
	"context"
	"database/sql"
	"fmt"
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
			Code:      ptr("applovin-banner-banner"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.1")),
			AccountID: applovinAccount.ID,
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547776,
				Valid: true,
			},
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
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547777,
				Valid: true,
			},
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
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547778,
				Valid: true,
			},
		},
		{
			AppID:     apps[0].ID,
			AdType:    db.InterstitialAdType,
			Code:      ptr("applovin-interstitial"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.3")),
			AccountID: applovinAccount.ID,
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547779,
				Valid: true,
			},
		},
		{
			AppID:     apps[0].ID,
			AdType:    db.InterstitialAdType,
			Code:      ptr("bidmachine-interstitial"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.3")),
			AccountID: bidmachineAccount.ID,
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
			Code:      ptr("app2-applovin-banner-mrec"),
			BidFloor:  decimal.NewNullDecimal(decimal.RequireFromString("0.4")),
			AccountID: applovinAccount.ID,
			PublicUID: sql.NullInt64{
				Int64: 1701972528521547781,
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
				{ID: "applovin", UID: "1701972528521547776", PriceFloor: 0.1, AdUnitID: "applovin-banner-banner"},
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2, AdUnitID: "applovin-banner-adaptive"},
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
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2, AdUnitID: "applovin-banner-adaptive"},
				{ID: "applovin", UID: "1701972528521547778", PriceFloor: 0.3, AdUnitID: "applovin-banner-leaderboard"},
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
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2, AdUnitID: "applovin-banner-adaptive"},
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
				{ID: "applovin", UID: "1701972528521547779", PriceFloor: 0.3, AdUnitID: "applovin-interstitial"},
				{ID: "bidmachine", UID: "1701972528521547780", PriceFloor: 0.3, AdUnitID: "bidmachine-interstitial"},
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
				{ID: "applovin", UID: "1701972528521547781", PriceFloor: 0.4, AdUnitID: "app2-applovin-banner-mrec"},
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
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2, AdUnitID: "applovin-banner-adaptive"},
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

func TestLineItemsMatcher_ExtraFields(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	user := dbtest.CreateUser(t, tx)

	keys := adapter.Keys
	demandSources := make([]db.DemandSource, len(keys))
	for i, key := range keys {
		demandSources[i] = dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
			source.APIKey = string(key)
		})
	}

	accounts := make([]db.DemandSourceAccount, len(demandSources))
	for i, demandSource := range demandSources {
		accounts[i] = dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
			account.User = user
			account.DemandSource = demandSource
			account.Extra = dbtest.ValidDemandSourceAccountExtra(t, adapter.Key(demandSource.APIKey))
		})
	}

	app := dbtest.CreateApp(t, tx, func(app *db.App) {
		app.User = user
	})

	for i, account := range accounts {
		dbtest.CreateLineItem(t, tx, func(lineItem *db.LineItem) {
			lineItem.App = app
			lineItem.Account = account
			lineItem.BidFloor = decimal.NewNullDecimal(decimal.RequireFromString("0.15"))
			lineItem.Extra = dbtest.ValidLineItemExtra(t, adapter.Key(demandSources[i].APIKey))
			lineItem.Code = ptr(fmt.Sprintf("code%d", i))
			lineItem.PublicUID = sql.NullInt64{
				Int64: int64(i),
				Valid: true,
			}
		})
	}

	matcher := store.LineItemsMatcher{DB: tx}

	params := func(adapters []adapter.Key) *auction.BuildParams {
		return &auction.BuildParams{
			AppID:      app.ID,
			AdType:     ad.BannerType,
			AdFormat:   ad.BannerFormat,
			DeviceType: device.PhoneType,
			Adapters:   adapters,
		}
	}

	testCases := []struct {
		params *auction.BuildParams
		want   []auction.LineItem
	}{
		{
			params: params([]adapter.Key{adapter.AdmobKey}),
			want: []auction.LineItem{
				{ID: "admob", UID: "0", PriceFloor: 0.15, AdUnitID: "admob_line_item"},
			},
		},
		{
			params: params([]adapter.Key{adapter.ApplovinKey}),
			want: []auction.LineItem{
				{ID: "applovin", UID: "2", PriceFloor: 0.15, AdUnitID: "code2", ZonedID: "applovin_line_item_zone_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.BidmachineKey}),
			want: []auction.LineItem{
				{ID: "bidmachine", UID: "3", PriceFloor: 0.15, AdUnitID: "code3"},
			},
		},
		{
			params: params([]adapter.Key{adapter.DTExchangeKey}),
			want: []auction.LineItem{
				{ID: "dtexchange", UID: "5", PriceFloor: 0.15, AdUnitID: "code5", PlacementID: "dt_exchange_line_item_spot_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.MetaKey}),
			want: []auction.LineItem{
				{ID: "meta", UID: "7", PriceFloor: 0.15, AdUnitID: "code7", PlacementID: "meta_line_item_placement_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.MintegralKey}),
			want: []auction.LineItem{
				{
					ID:          "mintegral",
					UID:         "8",
					PriceFloor:  0.15,
					AdUnitID:    "mintegral_line_item_unit_id",
					PlacementID: "mintegral_line_item_placement_id",
				},
			},
		},
		{
			params: params([]adapter.Key{adapter.MobileFuseKey}),
			want: []auction.LineItem{
				{ID: "mobilefuse", UID: "9", PriceFloor: 0.15, AdUnitID: "code9", PlacementID: "mobile_fuse_line_item_placement_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.UnityAdsKey}),
			want: []auction.LineItem{
				{ID: "unityads", UID: "10", PriceFloor: 0.15, AdUnitID: "code10", PlacementID: "unity_ads_line_item_placement_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.VungleKey}),
			want: []auction.LineItem{
				{ID: "vungle", UID: "11", PriceFloor: 0.15, AdUnitID: "code11", PlacementID: "vungle_line_item_placement_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.BigoAdsKey}),
			want: []auction.LineItem{
				{ID: "bigoads", UID: "4", PriceFloor: 0.15, AdUnitID: "code4", SlotID: "bigo_ads_line_item_slot_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.InmobiKey}),
			want: []auction.LineItem{
				{ID: "inmobi", UID: "6", PriceFloor: 0.15, AdUnitID: "code6", PlacementID: "inmobi_line_item_placement_id"},
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
