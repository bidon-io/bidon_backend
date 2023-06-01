package store

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
)

type AppRepo = resourceRepo[admin.App, admin.AppAttrs, app, *app]

type app struct {
	Model
	UserID      int64          `gorm:"column:user_id;type:bigint;not null"`
	PlatformID  int32          `gorm:"column:platform_id;type:integer;not null"`
	HumanName   string         `gorm:"column:human_name;type:varchar;not null"`
	PackageName sql.NullString `gorm:"column:package_name;type:varchar"`
	AppKey      sql.NullString `gorm:"column:app_key;type:varchar"`
	Settings    map[string]any `gorm:"column:settings;type:jsonb;default:'{}';serializer:json"`
}

//lint:ignore U1000 this method is used by generic struct
func (a *app) initFromResourceAttrs(attrs *admin.AppAttrs) {
	packageName := sql.NullString{}
	if attrs.PackageName != "" {
		packageName.String = attrs.PackageName
		packageName.Valid = true
	}

	appKey := sql.NullString{}
	if attrs.AppKey != "" {
		appKey.String = attrs.AppKey
		appKey.Valid = true
	}

	a.UserID = attrs.UserID
	a.PlatformID = dbPlatformID(attrs.PlatformID)
	a.HumanName = attrs.HumanName
	a.PackageName = packageName
	a.AppKey = appKey
	a.Settings = attrs.Settings
}

//lint:ignore U1000 this method is used by generic struct
func (a *app) toResource() admin.App {
	return admin.App{
		ID: a.ID,
		AppAttrs: admin.AppAttrs{
			UserID:      a.UserID,
			PlatformID:  platformID(a.PlatformID),
			HumanName:   a.HumanName,
			PackageName: a.PackageName.String,
			AppKey:      a.AppKey.String,
			Settings:    a.Settings,
		},
	}
}
