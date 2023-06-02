package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/store"
	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

func TestLineItemRepo_List(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.LineItemRepo{DB: tx}

	items := []admin.LineItemAttrs{
		{
			HumanName:   "banner",
			AppID:       1,
			BidFloor:    ptr(decimal.NewFromInt(1)),
			AdType:      admin.BannerAdType,
			Format:      ptr(admin.BannerLineItemFormat),
			AccountID:   1,
			AccountType: "DemandSourceAccount::Applovin",
			Code:        ptr("12345"),
			Extra:       map[string]any{"key": "value"},
		},
		{
			HumanName:   "interstitial",
			AppID:       2,
			BidFloor:    ptr(decimal.Decimal{}),
			AdType:      admin.InterstitialAdType,
			Format:      ptr(admin.EmptyLineItemFormat),
			AccountID:   2,
			AccountType: "DemandSourceAccount::Bidmachine",
			Code:        ptr(""),
			Extra:       map[string]any{"key": "value"},
		},
		{
			HumanName:   "rewarded",
			AppID:       3,
			BidFloor:    ptr(decimal.NewFromInt(3)),
			AdType:      admin.RewardedAdType,
			Format:      ptr(admin.EmptyLineItemFormat),
			AccountID:   3,
			AccountType: "DemandSourceAccount::UnityAds",
			Code:        ptr("54321"),
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
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.LineItemRepo{DB: tx}

	attrs := &admin.LineItemAttrs{
		HumanName:   "banner",
		AppID:       1,
		BidFloor:    ptr(decimal.NewFromInt(1)),
		AdType:      admin.BannerAdType,
		Format:      ptr(admin.BannerLineItemFormat),
		AccountID:   1,
		AccountType: "DemandSourceAccount::Applovin",
		Code:        ptr("12345"),
		Extra:       map[string]any{"key": "value"},
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestLineItemRepo_Update(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.LineItemRepo{DB: tx}

	attrs := admin.LineItemAttrs{
		HumanName:   "banner",
		AppID:       1,
		BidFloor:    ptr(decimal.NewFromInt(1)),
		AdType:      admin.BannerAdType,
		Format:      ptr(admin.BannerLineItemFormat),
		AccountID:   1,
		AccountType: "DemandSourceAccount::Applovin",
		Code:        ptr("12345"),
		Extra:       map[string]any{"key": "value"},
	}

	item, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, item, nil)
	}

	want := item
	want.AppID = 2
	want.BidFloor = ptr(decimal.Decimal{})
	want.Format = ptr(admin.EmptyLineItemFormat)
	want.Code = ptr("")

	updateParams := &admin.LineItemAttrs{
		AppID:    want.AppID,
		BidFloor: want.BidFloor,
		Format:   want.Format,
		Code:     want.Code,
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
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.LineItemRepo{DB: tx}

	attrs := &admin.LineItemAttrs{
		HumanName:   "banner",
		AppID:       1,
		BidFloor:    ptr(decimal.NewFromInt(1)),
		AdType:      admin.BannerAdType,
		Format:      ptr(admin.BannerLineItemFormat),
		AccountID:   1,
		AccountType: "DemandSourceAccount::Applovin",
		Code:        ptr("12345"),
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
