package store_test

import (
	"context"
	"database/sql"
	"github.com/bidon-io/bidon-backend/config"
	"testing"
	"time"

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

func TestLineItemsMatcher_MatchCached(t *testing.T) {
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

	matcher := store.LineItemsMatcher{
		DB:    tx,
		Cache: config.NewMemoryCacheOf[[]auction.LineItem](time.Minute),
	}
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
				{ID: "applovin", UID: "1701972528521547776", PriceFloor: 0.1},
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2},
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
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2},
				{ID: "applovin", UID: "1701972528521547778", PriceFloor: 0.3},
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
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2},
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
				{ID: "applovin", UID: "1701972528521547779", PriceFloor: 0.3},
				{ID: "bidmachine", UID: "1701972528521547780", PriceFloor: 0.3},
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
				{ID: "applovin", UID: "1701972528521547781", PriceFloor: 0.4},
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
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2},
			},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.BannerFormat,
				DeviceType: device.PhoneType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
			},
			want: []auction.LineItem{
				{ID: "applovin", UID: "1701972528521547776", PriceFloor: 0.1},
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2},
			},
		},
		{
			params: &auction.BuildParams{
				AppID:      apps[0].ID,
				AdType:     ad.BannerType,
				AdFormat:   ad.LeaderboardFormat,
				DeviceType: device.TabletType,
				Adapters:   []adapter.Key{adapter.ApplovinKey},
			},
			want: []auction.LineItem{
				{ID: "applovin", UID: "1701972528521547777", PriceFloor: 0.2},
				{ID: "applovin", UID: "1701972528521547778", PriceFloor: 0.3},
			},
		},
	}

	for _, tC := range testCases {
		got, err := matcher.MatchCached(context.Background(), tC.params)
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
			lineItem.PublicUID = sql.NullInt64{
				Int64: int64(i),
				Valid: true,
			}
		})
	}

	matcher := store.LineItemsMatcher{
		DB:    tx,
		Cache: config.NewMemoryCacheOf[[]auction.LineItem](time.Minute),
	}

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
				{ID: "admob", PriceFloor: 0.15, AdUnitID: "admob_line_item"},
			},
		},
		{
			params: params([]adapter.Key{adapter.ApplovinKey}),
			want: []auction.LineItem{
				{ID: "applovin", PriceFloor: 0.15, ZonedID: "applovin_line_item_zone_id", AdUnitID: "applovin_line_item_zone_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.BidmachineKey}),
			want: []auction.LineItem{
				{ID: "bidmachine", PriceFloor: 0.15},
			},
		},
		{
			params: params([]adapter.Key{adapter.DTExchangeKey}),
			want: []auction.LineItem{
				{ID: "dtexchange", PriceFloor: 0.15, PlacementID: "dt_exchange_line_item_spot_id", AdUnitID: "dt_exchange_line_item_spot_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.GAMKey}),
			want: []auction.LineItem{
				{ID: "gam", PriceFloor: 0.15, AdUnitID: "gam_line_item"},
			},
		},
		{
			params: params([]adapter.Key{adapter.MetaKey}),
			want: []auction.LineItem{
				{ID: "meta", PriceFloor: 0.15, PlacementID: "meta_line_item_placement_id", AdUnitID: "meta_line_item_placement_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.MintegralKey}),
			want: []auction.LineItem{
				{
					ID:          "mintegral",
					PriceFloor:  0.15,
					AdUnitID:    "mintegral_line_item_unit_id",
					PlacementID: "mintegral_line_item_placement_id",
				},
			},
		},
		{
			params: params([]adapter.Key{adapter.MobileFuseKey}),
			want: []auction.LineItem{
				{ID: "mobilefuse", PriceFloor: 0.15, PlacementID: "mobile_fuse_line_item_placement_id", AdUnitID: "mobile_fuse_line_item_placement_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.UnityAdsKey}),
			want: []auction.LineItem{
				{ID: "unityads", PriceFloor: 0.15, PlacementID: "unity_ads_line_item_placement_id", AdUnitID: "unity_ads_line_item_placement_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.VKAdsKey}),
			want: []auction.LineItem{
				{ID: "vkads", PriceFloor: 0.15, SlotID: "vk_ads_line_item_slot_id", Mediation: "vk_ads_line_item_mediation", AdUnitID: "vk_ads_line_item_slot_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.VungleKey}),
			want: []auction.LineItem{
				{ID: "vungle", PriceFloor: 0.15, PlacementID: "vungle_line_item_placement_id", AdUnitID: "vungle_line_item_placement_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.BigoAdsKey}),
			want: []auction.LineItem{
				{ID: "bigoads", PriceFloor: 0.15, SlotID: "bigo_ads_line_item_slot_id", AdUnitID: "bigo_ads_line_item_slot_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.InmobiKey}),
			want: []auction.LineItem{
				{ID: "inmobi", PriceFloor: 0.15, PlacementID: "inmobi_line_item_placement_id", AdUnitID: "inmobi_line_item_placement_id"},
			},
		},
		{
			params: params([]adapter.Key{adapter.YandexKey}),
			want: []auction.LineItem{
				{ID: "yandex", PriceFloor: 0.15, AdUnitID: "yandex_line_item_ad_unit_id"},
			},
		},
	}

	for _, tC := range testCases {
		got, err := matcher.Match(context.Background(), tC.params)
		if err != nil {
			t.Errorf("Error matching line items: %v", err)
		}
		for i := range got {
			got[i].UID = "" // UID is not deterministic
		}

		less := func(a, b auction.LineItem) bool { return a.AdUnitID < b.AdUnitID }
		if diff := cmp.Diff(tC.want, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("matcher.Match(ctx, %+v) mismatch (-want, +got)\n%s", tC.params, diff)
		}
	}
}
