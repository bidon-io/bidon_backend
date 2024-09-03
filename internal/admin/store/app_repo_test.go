package adminstore_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	adminstore "github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/google/go-cmp/cmp"
)

func TestAppRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAppRepo(tx)

	users := make([]db.User, 2)
	for i := range users {
		users[i] = dbtest.CreateUser(t, tx)
	}

	apps := []admin.AppAttrs{
		{
			PlatformID:  admin.IOSPlatformID,
			HumanName:   "App 1",
			PackageName: "com.example.app1",
			UserID:      users[0].ID,
		},
		{
			PlatformID:  admin.AndroidPlatformID,
			HumanName:   "App 2",
			PackageName: "com.example.app2",
			UserID:      users[1].ID,
		},
	}

	want := make([]admin.App, len(apps))
	for i, attrs := range apps {
		app, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, app, nil)
		}

		want[i] = *app
		want[i].User = *adminstore.UserResource(&users[i])
	}

	got, err := repo.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAppRepo_ListOwnedByUser(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	users := make([]db.User, 2)
	for i := range users {
		users[i] = dbtest.CreateUser(t, tx)
	}

	dbFirstUserApps := make([]db.App, 2)
	for i := range dbFirstUserApps {
		dbFirstUserApps[i] = dbtest.CreateApp(t, tx, func(app *db.App) {
			app.User = users[0]
		})
	}
	dbSecondUserApps := make([]db.App, 2)
	for i := range dbSecondUserApps {
		dbSecondUserApps[i] = dbtest.CreateApp(t, tx, func(app *db.App) {
			app.User = users[1]
		})
	}

	firstUserApps := make([]admin.App, 2)
	secondUserApps := make([]admin.App, 2)
	for i := 0; i < 2; i++ {
		firstUserApps[i] = adminstore.AppResource(&dbFirstUserApps[i])
		secondUserApps[i] = adminstore.AppResource(&dbSecondUserApps[i])
	}

	repo := adminstore.NewAppRepo(tx)

	tests := []struct {
		name   string
		userID int64
		want   []admin.App
	}{
		{
			"first user",
			users[0].ID,
			firstUserApps,
		},
		{
			"second user",
			users[1].ID,
			secondUserApps,
		},
		{
			"non-existent user",
			999,
			[]admin.App{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.ListOwnedByUser(context.Background(), tt.userID, nil)
			if err != nil {
				t.Fatalf("ListOwnedByUser() got %v; want %+v", err, tt.want)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("ListOwnedByUser() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestAppRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAppRepo(tx)

	user := dbtest.CreateUser(t, tx)
	attrs := &admin.AppAttrs{
		PlatformID:  admin.IOSPlatformID,
		HumanName:   "App 1",
		PackageName: "com.example.app1",
		UserID:      user.ID,
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}
	want.User = *adminstore.UserResource(&user)

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAppRepo_FindOwnedByUser(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	users := make([]db.User, 2)
	for i := range users {
		users[i] = dbtest.CreateUser(t, tx)
	}

	dbFirstUserApps := make([]db.App, 2)
	for i := range dbFirstUserApps {
		dbFirstUserApps[i] = dbtest.CreateApp(t, tx, func(app *db.App) {
			app.User = users[0]
		})
	}
	dbSecondUserApps := make([]db.App, 2)
	for i := range dbSecondUserApps {
		dbSecondUserApps[i] = dbtest.CreateApp(t, tx, func(app *db.App) {
			app.User = users[1]
		})
	}

	firstUserApps := make([]admin.App, 2)
	secondUserApps := make([]admin.App, 2)
	for i := 0; i < 2; i++ {
		firstUserApps[i] = adminstore.AppResource(&dbFirstUserApps[i])
		secondUserApps[i] = adminstore.AppResource(&dbSecondUserApps[i])
	}

	repo := adminstore.NewAppRepo(tx)

	tests := []struct {
		name    string
		userID  int64
		id      int64
		want    *admin.App
		wantErr bool
	}{
		{
			"first user, first user's app",
			users[0].ID,
			firstUserApps[0].ID,
			&firstUserApps[0],
			false,
		},
		{
			"first user, second user's app",
			users[0].ID,
			secondUserApps[0].ID,
			nil,
			true,
		},
		{
			"second user, second user's app",
			users[1].ID,
			secondUserApps[0].ID,
			&secondUserApps[0],
			false,
		},
		{
			"second user, first user's app",
			users[1].ID,
			firstUserApps[0].ID,
			nil,
			true,
		},
		{
			"non-existent user",
			999,
			999,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindOwnedByUser(context.Background(), tt.userID, tt.id)
			if tt.wantErr {
				if err == nil {
					t.Errorf("FindOwnedByUser() = %+v; want error", got)
				}
			} else if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FindOwnedByUser() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestAppRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAppRepo(tx)

	user := dbtest.CreateUser(t, tx)
	attrs := admin.AppAttrs{
		PlatformID:  admin.IOSPlatformID,
		HumanName:   "App 1",
		PackageName: "com.example.app1",
		UserID:      user.ID,
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

	repo := adminstore.NewAppRepo(tx)

	user := dbtest.CreateUser(t, tx)
	attrs := &admin.AppAttrs{
		PlatformID:  admin.IOSPlatformID,
		HumanName:   "App 1",
		PackageName: "com.example.app1",
		UserID:      user.ID,
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
