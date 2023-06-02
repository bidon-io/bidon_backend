package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/store"
	"github.com/google/go-cmp/cmp"
)

func TestAppDemandProfileRepo_List(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AppDemandProfileRepo{DB: tx}

	profiles := []admin.AppDemandProfileAttrs{
		{
			AppID:          1,
			DemandSourceID: 1,
			AccountID:      1,
			Data:           map[string]any{"api_key": "asdf"},
			AccountType:    "DemandSourceAccount::Applovin",
		},
		{
			AppID:          2,
			DemandSourceID: 2,
			AccountID:      2,
			Data:           map[string]any{"api_key": "asdf"},
			AccountType:    "DemandSourceAccount::Bidmachine",
		},
	}

	want := make([]admin.AppDemandProfile, len(profiles))
	for i, attrs := range profiles {
		profile, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, profile, nil)
		}

		want[i] = *profile
	}

	got, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAppDemandProfileRepo_Find(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AppDemandProfileRepo{DB: tx}

	attrs := &admin.AppDemandProfileAttrs{
		AppID:          1,
		DemandSourceID: 1,
		AccountID:      1,
		Data:           map[string]any{"api_key": "asdf"},
		AccountType:    "DemandSourceAccount::Applovin",
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

func TestAppDemandProfileRepo_Update(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AppDemandProfileRepo{DB: tx}

	attrs := admin.AppDemandProfileAttrs{
		AppID:          1,
		DemandSourceID: 1,
		AccountID:      1,
		Data:           map[string]any{"api_key": "asdf"},
		AccountType:    "DemandSourceAccount::Applovin",
	}

	profile, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, profile, nil)
	}

	want := profile
	want.AppID = 2

	updateParams := &admin.AppDemandProfileAttrs{
		AppID: want.AppID,
	}
	got, err := repo.Update(context.Background(), profile.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAppDemandProfileRepo_Delete(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AppDemandProfileRepo{DB: tx}

	attrs := &admin.AppDemandProfileAttrs{
		AppID:          1,
		DemandSourceID: 1,
		AccountID:      1,
		Data:           map[string]any{"api_key": "asdf"},
		AccountType:    "DemandSourceAccount::Applovin",
	}
	profile, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, profile, nil)
	}

	err = repo.Delete(context.Background(), profile.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", profile.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), profile.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", profile.ID, got, err, nil, "record not found")
	}
}
