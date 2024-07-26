package adminstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

func TestLineItemRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewLineItemRepo(tx)

	users := make([]db.User, 2)
	for i := range users {
		users[i] = dbtest.CreateUser(t, tx)
	}
	apps := make([]db.App, 2)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx, func(app *db.App) {
			app.User = users[i]
		})
	}

	applovinDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.ApplovinKey)
		source.HumanName = source.APIKey
	})
	applovinAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.User = users[0]
		account.DemandSource = applovinDemandSource
	})

	bidmachineDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.BidmachineKey)
		source.HumanName = source.APIKey
	})
	bidmachineAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.User = users[0]
		account.DemandSource = bidmachineDemandSource
	})

	unityAdsDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.UnityAdsKey)
		source.HumanName = source.APIKey
	})
	unityAdsAccount1 := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.User = users[0]
		account.DemandSource = unityAdsDemandSource
	})
	unityAdsAccount2 := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.User = users[1]
		account.DemandSource = unityAdsDemandSource
	})

	items := []struct {
		*admin.LineItemAttrs
		App     db.App
		Account db.DemandSourceAccount
	}{
		{
			LineItemAttrs: &admin.LineItemAttrs{
				HumanName:   "banner",
				AppID:       apps[0].ID,
				BidFloor:    ptr(decimal.NewFromInt(1)),
				AdType:      ad.BannerType,
				Format:      ptr(ad.BannerFormat),
				AccountID:   applovinAccount.ID,
				AccountType: applovinAccount.Type,
				Extra:       map[string]any{"key": "value"},
			},
			App:     apps[0],
			Account: applovinAccount,
		},
		{
			LineItemAttrs: &admin.LineItemAttrs{
				HumanName:   "interstitial",
				AppID:       apps[0].ID,
				BidFloor:    ptr(decimal.Decimal{}),
				AdType:      ad.InterstitialType,
				Format:      ptr(ad.EmptyFormat),
				AccountID:   bidmachineAccount.ID,
				AccountType: bidmachineAccount.Type,
				Extra:       map[string]any{"key": "value"},
			},
			App:     apps[0],
			Account: bidmachineAccount,
		},
		{
			LineItemAttrs: &admin.LineItemAttrs{
				HumanName:   "rewarded",
				AppID:       apps[0].ID,
				BidFloor:    ptr(decimal.NewFromInt(3)),
				AdType:      ad.RewardedType,
				Format:      ptr(ad.EmptyFormat),
				AccountID:   unityAdsAccount1.ID,
				AccountType: unityAdsAccount1.Type,
				Extra:       map[string]any{"key": "value"},
			},
			App:     apps[0],
			Account: unityAdsAccount1,
		},
		{
			LineItemAttrs: &admin.LineItemAttrs{
				HumanName:   "rewarded App 2",
				AppID:       apps[1].ID,
				BidFloor:    ptr(decimal.NewFromInt(3)),
				AdType:      ad.RewardedType,
				Format:      ptr(ad.EmptyFormat),
				AccountID:   unityAdsAccount2.ID,
				AccountType: unityAdsAccount2.Type,
				IsBidding:   ptr(true),
				Extra:       map[string]any{"key": "value"},
			},
			App:     apps[1],
			Account: unityAdsAccount2,
		},
	}

	allItems := make([]admin.LineItem, len(items))
	for i, attrs := range items {
		item, err := repo.Create(context.Background(), attrs.LineItemAttrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; allItems %T, %v", &attrs, nil, err, item, nil)
		}

		allItems[i] = *item
		allItems[i].Account = adminstore.DemandSourceAccountAttrsWithId(&attrs.Account)
		allItems[i].App = adminstore.AppAttrsWithId(&attrs.App)
	}

	testcases := []struct {
		name    string
		qParams map[string][]string
		want    []admin.LineItem
		wantErr bool
	}{
		{
			name:    "no filters",
			qParams: nil,
			want:    allItems,
		},
		{
			name: "filter by user_id",
			qParams: map[string][]string{
				"user_id": {fmt.Sprint(users[0].ID)},
			},
			want: allItems[:3],
		},
		{
			name: "filter by app_id",
			qParams: map[string][]string{
				"app_id": {fmt.Sprint(apps[0].ID)},
			},
			want: allItems[:3],
		},
		{
			name: "filter by ad_type",
			qParams: map[string][]string{
				"ad_type": {string(ad.RewardedType)},
			},
			want: allItems[2:],
		},
		{
			name: "filter by account_id",
			qParams: map[string][]string{
				"account_id": {fmt.Sprint(unityAdsAccount1.ID)},
			},
			want: allItems[2:3],
		},
		{
			name: "filter by account_type",
			qParams: map[string][]string{
				"account_type": {unityAdsAccount1.Type},
			},
			want: allItems[2:],
		},
		{
			name: "filter by is_bidding true",
			qParams: map[string][]string{
				"is_bidding": {"true"},
			},
			want: allItems[3:],
		},
		{
			name: "filter by is_bidding false",
			qParams: map[string][]string{
				"is_bidding": {"false"},
			},
			want: allItems[:3],
		},
		{
			name: "filter by AppID and AccountID",
			qParams: map[string][]string{
				"app_id":     {fmt.Sprint(apps[0].ID)},
				"account_id": {fmt.Sprint(applovinAccount.ID)},
			},
			want: allItems[:1],
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := repo.List(context.Background(), tc.qParams)
			if err != nil {
				t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, tc.want, nil)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestLineItemRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewLineItemRepo(tx)

	app := dbtest.CreateApp(t, tx)
	applovinDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.ApplovinKey)
		source.HumanName = source.APIKey
	})
	applovinAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = applovinDemandSource
	})

	attrs := &admin.LineItemAttrs{
		HumanName:   "banner",
		AppID:       app.ID,
		BidFloor:    ptr(decimal.NewFromInt(1)),
		AdType:      ad.BannerType,
		Format:      ptr(ad.BannerFormat),
		AccountID:   applovinAccount.ID,
		AccountType: applovinAccount.Type,
		Extra:       map[string]any{"key": "value"},
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}
	want.App = adminstore.AppAttrsWithId(&app)
	want.Account = adminstore.DemandSourceAccountAttrsWithId(&applovinAccount)

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestLineItemRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewLineItemRepo(tx)

	app := dbtest.CreateApp(t, tx)
	applovinDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.ApplovinKey)
		source.HumanName = source.APIKey
	})
	applovinAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = applovinDemandSource
	})

	attrs := admin.LineItemAttrs{
		HumanName:   "banner",
		AppID:       app.ID,
		BidFloor:    ptr(decimal.NewFromInt(1)),
		AdType:      ad.BannerType,
		Format:      ptr(ad.BannerFormat),
		AccountID:   applovinAccount.ID,
		AccountType: applovinAccount.Type,
		Extra:       map[string]any{"key": "value"},
	}

	item, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, item, nil)
	}

	want := item
	want.BidFloor = ptr(decimal.Decimal{})
	want.Format = ptr(ad.EmptyFormat)

	updateParams := &admin.LineItemAttrs{
		BidFloor: want.BidFloor,
		Format:   want.Format,
	}
	got, err := repo.Update(context.Background(), item.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestLineItemRepo_Delete(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewLineItemRepo(tx)

	app := dbtest.CreateApp(t, tx)
	applovinDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.ApplovinKey)
		source.HumanName = source.APIKey
	})
	applovinAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = applovinDemandSource
	})
	attrs := &admin.LineItemAttrs{
		HumanName:   "banner",
		AppID:       app.ID,
		BidFloor:    ptr(decimal.NewFromInt(1)),
		AdType:      ad.BannerType,
		Format:      ptr(ad.BannerFormat),
		AccountID:   applovinAccount.ID,
		AccountType: applovinAccount.Type,
		Extra:       map[string]any{"key": "value"},
	}
	item, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, item, nil)
	}

	err = repo.Delete(context.Background(), item.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", item.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), item.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", item.ID, got, err, nil, "record not found")
	}
}
