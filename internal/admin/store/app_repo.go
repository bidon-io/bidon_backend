package adminstore

import (
	"context"
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
	"gorm.io/gorm"
)

type AppRepo struct {
	*resourceRepo[admin.App, admin.AppAttrs, db.App]
}

func NewAppRepo(d *db.DB) *AppRepo {
	return &AppRepo{
		resourceRepo: &resourceRepo[admin.App, admin.AppAttrs, db.App]{
			db:           d,
			mapper:       appMapper{},
			associations: []string{"User"},
		},
	}
}

func (r *AppRepo) ListOwnedByUser(ctx context.Context, userID int64) ([]admin.App, error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	})
}

func (r *AppRepo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.App, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	})
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
		User:     userMapper{}.resource(&a.User),
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
