package adminstore_test

import (
	"context"
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

	user := dbtest.CreateUser(t, tx)
	app := dbtest.CreateApp(t, tx, func(app *db.App) {
		app.User = user
	})

	applovinDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.ApplovinKey)
		source.HumanName = source.APIKey
	})
	applovinAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.User = user
		account.DemandSource = applovinDemandSource
	})

	bidmachineDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.BidmachineKey)
		source.HumanName = source.APIKey
	})
	bidmachineAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.User = user
		account.DemandSource = bidmachineDemandSource
	})

	unityAdsDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.UnityAdsKey)
		source.HumanName = source.APIKey
	})
	unityAdsAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.User = user
		account.DemandSource = unityAdsDemandSource
	})

	accounts := []db.DemandSourceAccount{applovinAccount, bidmachineAccount, unityAdsAccount}
	items := []admin.LineItemAttrs{
		{
			HumanName:   "banner",
			AppID:       app.ID,
			BidFloor:    ptr(decimal.NewFromInt(1)),
			AdType:      ad.BannerType,
			Format:      ptr(ad.BannerFormat),
			AccountID:   applovinAccount.ID,
			AccountType: applovinAccount.Type,
			Extra:       map[string]any{"key": "value"},
		},
		{
			HumanName:   "interstitial",
			AppID:       app.ID,
			BidFloor:    ptr(decimal.Decimal{}),
			AdType:      ad.InterstitialType,
			Format:      ptr(ad.EmptyFormat),
			AccountID:   bidmachineAccount.ID,
			AccountType: bidmachineAccount.Type,
			Extra:       map[string]any{"key": "value"},
		},
		{
			HumanName:   "rewarded",
			AppID:       app.ID,
			BidFloor:    ptr(decimal.NewFromInt(3)),
			AdType:      ad.RewardedType,
			Format:      ptr(ad.EmptyFormat),
			AccountID:   unityAdsAccount.ID,
			AccountType: unityAdsAccount.Type,
			Extra:       map[string]any{"key": "value"},
		},
	}

	want := make([]admin.LineItem, len(items))
	for i, attrs := range items {
		item, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, item, nil)
		}

		want[i] = *item
		want[i].Account = adminstore.DemandSourceAccountAttrsWithId(&accounts[i])
		want[i].App = adminstore.AppAttrsWithId(&app)
	}

	got, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
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
