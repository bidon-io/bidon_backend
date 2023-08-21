package dbtest

import (
	"database/sql"
	"fmt"

	"github.com/bidon-io/bidon-backend/internal/db"
)

type AppFactory struct {
	User        func(int) db.User
	PlatformID  func(int) db.PlatformID
	HumanName   func(int) string
	PackageName func(int) string
	AppKey      func(int) string
	Settings    func(int) map[string]any
}

func (f AppFactory) Build(i int) db.App {
	app := db.App{}

	var user db.User
	if f.User == nil {
		user = UserFactory{}.Build(i)
	} else {
		user = f.User(i)
	}
	app.UserID = user.ID
	app.User = user

	if f.PlatformID == nil {
		app.PlatformID = db.AndroidPlatformID
	} else {
		app.PlatformID = f.PlatformID(i)
	}

	if f.HumanName == nil {
		app.HumanName = fmt.Sprintf("Test App %d", i)
	} else {
		app.HumanName = f.HumanName(i)
	}

	if f.PackageName == nil {
		app.PackageName = sql.NullString{
			String: fmt.Sprintf("app.package%d", i),
			Valid:  true,
		}
	} else {
		app.PackageName = sql.NullString{
			String: f.PackageName(i),
			Valid:  true,
		}
	}

	if f.AppKey == nil {
		app.AppKey = sql.NullString{
			String: fmt.Sprintf("appkey%d", i),
			Valid:  true,
		}
	} else {
		app.AppKey = sql.NullString{
			String: f.AppKey(i),
			Valid:  true,
		}
	}

	if f.Settings == nil {
		app.Settings = map[string]any{"foo": "bar"}
	} else {
		app.Settings = f.Settings(i)
	}

	return app
}
