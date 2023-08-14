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

func TestAppDemandProfileRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAppDemandProfileRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	apps := make([]*db.App, 2)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx, i, user)
	}
	demandSources := dbtest.CreateDemandSourcesList(t, tx, 2)
	accounts := make([]*db.DemandSourceAccount, 2)
	for i := range accounts {
		accounts[i] = dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(&db.DemandSourceAccount{
			UserID:         user.ID,
			DemandSourceID: demandSources[i].ID,
		}))
	}
	profiles := []admin.AppDemandProfileAttrs{
		{
			AppID:          apps[0].ID,
			DemandSourceID: demandSources[0].ID,
			AccountID:      accounts[0].ID,
			Data:           map[string]any{"api_key": "asdf"},
			AccountType:    "DemandSourceAccount::Applovin",
		},
		{
			AppID:          apps[1].ID,
			DemandSourceID: demandSources[1].ID,
			AccountID:      accounts[1].ID,
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
		want[i].App = *adminstore.AppResource(apps[i])
		want[i].DemandSource = *adminstore.DemandSourceResource(demandSources[i])
		want[i].Account = *adminstore.DemandSourceAccountResource(accounts[i])
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
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAppDemandProfileRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	app := dbtest.CreateApp(t, tx, 1, user)
	demandSource := dbtest.CreateDemandSource(t, tx)
	account := dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(&db.DemandSourceAccount{
		UserID:         user.ID,
		DemandSourceID: demandSource.ID,
	}))
	attrs := &admin.AppDemandProfileAttrs{
		AppID:          app.ID,
		DemandSourceID: demandSource.ID,
		AccountID:      account.ID,
		Data:           map[string]any{"api_key": "asdf"},
		AccountType:    "DemandSourceAccount::Applovin",
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}
	want.App = *adminstore.AppResource(app)
	want.Account = *adminstore.DemandSourceAccountResource(account)
	want.DemandSource = *adminstore.DemandSourceResource(demandSource)

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAppDemandProfileRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAppDemandProfileRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	app := dbtest.CreateApp(t, tx, 1, user)
	demandSource := dbtest.CreateDemandSource(t, tx)
	account := dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(&db.DemandSourceAccount{
		UserID:         user.ID,
		DemandSourceID: demandSource.ID,
	}))
	attrs := admin.AppDemandProfileAttrs{
		AppID:          app.ID,
		DemandSourceID: demandSource.ID,
		AccountID:      account.ID,
		Data:           map[string]any{"api_key": "asdf"},
		AccountType:    "DemandSourceAccount::Applovin",
	}

	profile, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, profile, nil)
	}

	want := profile
	want.Data = map[string]any{"api_key": "new_api_key"}

	updateParams := &admin.AppDemandProfileAttrs{
		Data: want.Data,
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
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAppDemandProfileRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	app := dbtest.CreateApp(t, tx, 1, user)
	demandSource := dbtest.CreateDemandSource(t, tx)
	account := dbtest.CreateDemandSourceAccount(t, tx, dbtest.WithDemandSourceAccountOptions(&db.DemandSourceAccount{
		UserID:         user.ID,
		DemandSourceID: demandSource.ID,
	}))
	attrs := &admin.AppDemandProfileAttrs{
		AppID:          app.ID,
		DemandSourceID: demandSource.ID,
		AccountID:      account.ID,
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
