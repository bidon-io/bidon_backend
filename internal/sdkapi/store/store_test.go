package store

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/google/go-cmp/cmp"
)

var testDB *db.DB

func TestMain(m *testing.M) {
	testDB = dbtest.Prepare()

	os.Exit(m.Run())
}

func TestAppFetcher_Fetch(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	user := dbtest.CreateUser(t, tx, 1)
	app := &db.App{
		UserID:      user.ID,
		AppKey:      sql.NullString{String: "asdf", Valid: true},
		PackageName: sql.NullString{String: "com.example.app", Valid: true},
	}
	if err := tx.Create(app).Error; err != nil {
		t.Fatalf("Error creating app: %v", err)
	}

	fetcher := &AppFetcher{DB: tx}

	testCases := []struct {
		name      string
		appKey    string
		appBundle string
		want      any
	}{
		{
			name:      "App matches",
			appKey:    app.AppKey.String,
			appBundle: app.PackageName.String,
			want:      sdkapi.App{ID: app.ID},
		},
		{
			name:      "App key does not match",
			appKey:    "fdsa",
			appBundle: app.PackageName.String,
			want:      sdkapi.ErrAppNotValid,
		},
		{
			name:      "App bundle does not match",
			appKey:    app.AppKey.String,
			appBundle: "not.found",
			want:      sdkapi.ErrAppNotValid,
		},
		{
			name:      "Nothing matches",
			appKey:    "fdsa",
			appBundle: "not.found",
			want:      sdkapi.ErrAppNotValid,
		},
	}

	for _, tC := range testCases {
		app, err := fetcher.Fetch(context.Background(), tC.appKey, tC.appBundle)

		var got any
		switch tC.want.(type) {
		case sdkapi.App:
			got = app
		case error:
			got = err
		}

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("fetcher.Fetch -> %v mismatch (-want +got):\n%s", tC.name, diff)
		}
	}
}
