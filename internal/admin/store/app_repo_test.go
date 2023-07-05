package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/google/go-cmp/cmp"
)

func TestAppRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewAppRepo(tx)

	users := dbtest.CreateUsersList(t, tx, 2)
	apps := []admin.AppAttrs{
		{
			PlatformID:  admin.IOSPlatformID,
			HumanName:   "App 1",
			PackageName: "com.example.app1",
			UserID:      users[0].ID,
			AppKey:      "qwerty",
			Settings:    map[string]any{"setting1": 1, "setting2": 2},
		},
		{
			PlatformID:  admin.AndroidPlatformID,
			HumanName:   "App 2",
			PackageName: "com.example.app2",
			UserID:      users[1].ID,
			AppKey:      "asdfg",
			Settings:    map[string]any{"setting1": 1, "setting2": 2},
		},
	}

	want := make([]admin.App, len(apps))
	for i, attrs := range apps {
		app, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, app, nil)
		}

		want[i] = *app
		want[i].User = *store.UserResource(users[i])
	}

	got, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAppRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewAppRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	attrs := &admin.AppAttrs{
		PlatformID:  admin.IOSPlatformID,
		HumanName:   "App 1",
		PackageName: "com.example.app1",
		UserID:      user.ID,
		AppKey:      "qwerty",
		Settings:    map[string]any{"setting1": 1, "setting2": 2},
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}
	want.User = *store.UserResource(user)

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAppRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewAppRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	attrs := admin.AppAttrs{
		PlatformID:  admin.IOSPlatformID,
		HumanName:   "App 1",
		PackageName: "com.example.app1",
		UserID:      user.ID,
		AppKey:      "qwerty",
		Settings:    map[string]any{"setting1": 1, "setting2": 2},
	}

	app, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, app, nil)
	}

	want := app
	want.PlatformID = admin.AndroidPlatformID

	updateParams := &admin.AppAttrs{
		PlatformID: want.PlatformID,
	}
	got, err := repo.Update(context.Background(), app.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAppRepo_Delete(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewAppRepo(tx)

	user := dbtest.CreateUser(t, tx, 1)
	attrs := &admin.AppAttrs{
		PlatformID:  admin.IOSPlatformID,
		HumanName:   "App 1",
		PackageName: "com.example.app1",
		UserID:      user.ID,
		AppKey:      "qwerty",
		Settings:    map[string]any{"setting1": 1, "setting2": 2},
	}
	app, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, app, nil)
	}

	err = repo.Delete(context.Background(), app.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", app.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), app.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", app.ID, got, err, nil, "record not found")
	}
}
