package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/google/go-cmp/cmp"
)

func createDemandSource(t *testing.T, tx *db.DB, APIKey string) *db.DemandSource {
	t.Helper()

	demandSource := &db.DemandSource{
		APIKey:    APIKey,
		HumanName: APIKey,
	}
	err := tx.Create(demandSource).Error
	if err != nil {
		t.Fatalf("Error creating demand source: %v", err)
	}

	return demandSource
}

func TestDemandSourceAccountRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewDemandSourceAccountRepo(tx)

	accounts := []admin.DemandSourceAccountAttrs{
		{
			UserID:         1,
			Type:           "DemandSourceAccount::Applovin",
			DemandSourceID: createDemandSource(t, tx, "applovin").ID,
			IsBidding:      ptr(false),
			Extra:          map[string]any{"key": "value"},
		},
		{
			UserID:         1,
			Type:           "DemandSourceAccount::Bidmachine",
			DemandSourceID: createDemandSource(t, tx, "bidmachine").ID,
			IsBidding:      ptr(true),
			Extra:          map[string]any{"key": "value"},
		},
		{
			UserID:         1,
			Type:           "DemandSourceAccount::UnityAds",
			DemandSourceID: createDemandSource(t, tx, "unityads").ID,
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
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewDemandSourceAccountRepo(tx)

	attrs := &admin.DemandSourceAccountAttrs{
		UserID:         1,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: createDemandSource(t, tx, "bidmachine").ID,
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
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewDemandSourceAccountRepo(tx)

	attrs := admin.DemandSourceAccountAttrs{
		UserID:         1,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: createDemandSource(t, tx, "bidmachine").ID,
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
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewDemandSourceAccountRepo(tx)

	attrs := &admin.DemandSourceAccountAttrs{
		UserID:         1,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: createDemandSource(t, tx, "bidmachine").ID,
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
