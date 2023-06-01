package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/store"
	"github.com/google/go-cmp/cmp"
)

func TestAppRepo_List(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AppRepo{DB: tx}

	apps := []admin.AppAttrs{
		{
			PlatformID:  admin.IOSPlatformID,
			HumanName:   "App 1",
			PackageName: "com.example.app1",
			UserID:      1,
			AppKey:      "qwerty",
			Settings:    map[string]any{"setting1": 1, "setting2": 2},
		},
		{
			PlatformID:  admin.AndroidPlatformID,
			HumanName:   "App 2",
			PackageName: "com.example.app2",
			UserID:      2,
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
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AppRepo{DB: tx}

	attrs := &admin.AppAttrs{
		PlatformID:  admin.IOSPlatformID,
		HumanName:   "App 1",
		PackageName: "com.example.app1",
		UserID:      1,
		AppKey:      "qwerty",
		Settings:    map[string]any{"setting1": 1, "setting2": 2},
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

func TestAppRepo_Update(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AppRepo{DB: tx}

	attrs := admin.AppAttrs{
		PlatformID:  admin.IOSPlatformID,
		HumanName:   "App 1",
		PackageName: "com.example.app1",
		UserID:      1,
		AppKey:      "qwerty",
		Settings:    map[string]any{"setting1": 1, "setting2": 2},
	}

	app, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, app, nil)
	}

	want := app
	want.UserID = 2

	updateParams := &admin.AppAttrs{
		UserID: want.UserID,
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
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AppRepo{DB: tx}

	attrs := &admin.AppAttrs{
		PlatformID:  admin.IOSPlatformID,
		HumanName:   "App 1",
		PackageName: "com.example.app1",
		UserID:      1,
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
