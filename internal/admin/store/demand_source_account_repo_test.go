package adminstore_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/google/go-cmp/cmp"
)

func TestDemandSourceAccountRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewDemandSourceAccountRepo(tx)
	demandSources := make([]*db.DemandSource, 3)
	demandSources[0] = dbtest.CreateDemandSource(t, tx, dbtest.WithDemandSourceOptions(&db.DemandSource{
		APIKey: "applovin",
	}))
	demandSources[1] = dbtest.CreateDemandSource(t, tx, dbtest.WithDemandSourceOptions(&db.DemandSource{
		APIKey: "bidmachine",
	}))
	demandSources[2] = dbtest.CreateDemandSource(t, tx, dbtest.WithDemandSourceOptions(&db.DemandSource{
		APIKey: "unityads",
	}))
	user := dbtest.CreateUser(t, tx, 1)
	accounts := []admin.DemandSourceAccountAttrs{
		{
			UserID:         user.ID,
			Type:           "DemandSourceAccount::Applovin",
			DemandSourceID: demandSources[0].ID,
			IsBidding:      ptr(false),
			Extra:          map[string]any{"key": "value"},
		},
		{
			UserID:         user.ID,
			Type:           "DemandSourceAccount::Bidmachine",
			DemandSourceID: demandSources[1].ID,
			IsBidding:      ptr(true),
			Extra:          map[string]any{"key": "value"},
		},
		{
			UserID:         user.ID,
			Type:           "DemandSourceAccount::UnityAds",
			DemandSourceID: demandSources[2].ID,
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
		want[i].User = *adminstore.UserResource(user)
		want[i].DemandSource = *adminstore.DemandSourceResource(demandSources[i])
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

	repo := adminstore.NewDemandSourceAccountRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	demandSource := dbtest.CreateDemandSource(t, tx, dbtest.WithDemandSourceOptions(&db.DemandSource{
		APIKey: "bidmachine",
	}))
	attrs := &admin.DemandSourceAccountAttrs{
		UserID:         user.ID,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: demandSource.ID,
		IsBidding:      ptr(true),
		Extra:          map[string]any{"key": "value"},
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}
	want.User = *adminstore.UserResource(user)
	want.DemandSource = *adminstore.DemandSourceResource(demandSource)

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

	repo := adminstore.NewDemandSourceAccountRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	demandSource := dbtest.CreateDemandSource(t, tx, dbtest.WithDemandSourceOptions(&db.DemandSource{
		APIKey: "bidmachine",
	}))
	attrs := admin.DemandSourceAccountAttrs{
		UserID:         user.ID,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: demandSource.ID,
		IsBidding:      ptr(true),
		Extra:          map[string]any{"key": "value"},
	}

	account, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, account, nil)
	}

	want := account
	want.Extra = map[string]any{"key": "value2"}
	want.IsBidding = ptr(false)

	updateParams := &admin.DemandSourceAccountAttrs{
		Extra:     want.Extra,
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

	repo := adminstore.NewDemandSourceAccountRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	demandSource := dbtest.CreateDemandSource(t, tx, dbtest.WithDemandSourceOptions(&db.DemandSource{
		APIKey: "bidmachine",
	}))
	attrs := &admin.DemandSourceAccountAttrs{
		UserID:         user.ID,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: demandSource.ID,
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
