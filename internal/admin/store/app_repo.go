package adminstore

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type AppRepo = resourceRepo[admin.App, admin.AppAttrs, db.App]

func NewAppRepo(db *db.DB) *AppRepo {
	return &AppRepo{
		db:           db,
		mapper:       appMapper{},
		associations: []string{"User"},
	}
}

type appMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m appMapper) dbModel(a *admin.AppAttrs, id int64) *db.App {
	packageName := sql.NullString{}
	if a.PackageName != "" {
		packageName.String = a.PackageName
		packageName.Valid = true
	}

	appKey := sql.NullString{}
	if a.AppKey != "" {
		appKey.String = a.AppKey
		appKey.Valid = true
	}

	return &db.App{
		Model:       db.Model{ID: id},
		UserID:      a.UserID,
		PlatformID:  dbPlatformID(a.PlatformID),
		HumanName:   a.HumanName,
		PackageName: packageName,
		AppKey:      appKey,
		Settings:    a.Settings,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m appMapper) resource(a *db.App) admin.App {
	return admin.App{
		ID:       a.ID,
		AppAttrs: m.resourceAttrs(a),
		User: admin.User{
			ID:        a.User.ID,
			UserAttrs: userMapper{}.resourceAttrs(&a.User),
		},
	}
}

func (m appMapper) resourceAttrs(a *db.App) admin.AppAttrs {
	return admin.AppAttrs{
		UserID:      a.UserID,
		PlatformID:  platformID(a.PlatformID),
		HumanName:   a.HumanName,
		PackageName: a.PackageName.String,
		AppKey:      a.AppKey.String,
		Settings:    a.Settings,
	}
}
