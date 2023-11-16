package dbtest

import (
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/db"
)

var appCounter atomic.Uint64

func appDefaults(n uint64) func(*db.App) {
	return func(app *db.App) {
		if app.UserID == 0 && app.User.ID == 0 {
			app.User = BuildUser(func(user *db.User) {
				*user = app.User
			})
		}
		if app.PlatformID == 0 {
			app.PlatformID = db.AndroidPlatformID
		}
		if app.HumanName == "" {
			app.HumanName = fmt.Sprintf("Test App %d", n)
		}
		if app.PackageName == (sql.NullString{}) {
			app.PackageName = sql.NullString{
				String: fmt.Sprintf("app_%d.package", n),
				Valid:  true,
			}
		}
		if app.AppKey == (sql.NullString{}) {
			app.AppKey = sql.NullString{
				String: fmt.Sprintf("app_%d_key", n),
				Valid:  true,
			}
		}
		if app.Settings == nil {
			app.Settings = map[string]any{
				"app_num": n,
				"foo":     "bar",
			}
		}
	}
}

func BuildApp(opts ...func(*db.App)) db.App {
	var app db.App

	n := appCounter.Add(1)

	opts = append(opts, appDefaults(n))
	for _, opt := range opts {
		opt(&app)
	}

	return app
}

func CreateApp(t *testing.T, tx *db.DB, opts ...func(*db.App)) db.App {
	t.Helper()

	app := BuildApp(opts...)
	if err := tx.Create(&app).Error; err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	return app
}
