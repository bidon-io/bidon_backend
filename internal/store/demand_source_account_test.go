package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/store"
	"github.com/google/go-cmp/cmp"
)

func TestDemandSourceAccountRepo_List(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.DemandSourceAccountRepo{DB: tx}

	accounts := []admin.DemandSourceAccountAttrs{
		{
			UserID:         1,
			Type:           "DemandSourceAccount::Applovin",
			DemandSourceID: 1,
			IsBidding:      ptr(false),
			Extra:          map[string]any{"key": "value"},
		},
		{
			UserID:         1,
			Type:           "DemandSourceAccount::Bidmachine",
			DemandSourceID: 2,
			IsBidding:      ptr(true),
			Extra:          map[string]any{"key": "value"},
		},
		{
			UserID:         1,
			Type:           "DemandSourceAccount::UnityAds",
			DemandSourceID: 3,
			IsBidding:      nil,
			Extra:          map[string]any{"key": "value"},
		},
	}

	want := make([]admin.DemandSourceAccount, len(accounts))
	for i, attrs := range accounts {
		account, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, account, nil)
		}

		want[i] = *account
	}

	got, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestDemandSourceAccountRepo_Find(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.DemandSourceAccountRepo{DB: tx}

	attrs := &admin.DemandSourceAccountAttrs{
		UserID:         1,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: 2,
		IsBidding:      ptr(true),
		Extra:          map[string]any{"key": "value"},
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

func TestDemandSourceAccountRepo_Update(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.DemandSourceAccountRepo{DB: tx}

	attrs := admin.DemandSourceAccountAttrs{
		UserID:         1,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: 2,
		IsBidding:      ptr(true),
		Extra:          map[string]any{"key": "value"},
	}

	account, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, account, nil)
	}

	want := account
	want.UserID = 2
	want.IsBidding = ptr(false)

	updateParams := &admin.DemandSourceAccountAttrs{
		UserID:    want.UserID,
		IsBidding: want.IsBidding,
	}
	got, err := repo.Update(context.Background(), account.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestDemandSourceAccountRepo_Delete(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.DemandSourceAccountRepo{DB: tx}

	attrs := &admin.DemandSourceAccountAttrs{
		UserID:         1,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: 2,
		IsBidding:      ptr(true),
		Extra:          map[string]any{"key": "value"},
	}
	account, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, account, nil)
	}

	err = repo.Delete(context.Background(), account.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", account.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), account.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", account.ID, got, err, nil, "record not found")
	}
}
