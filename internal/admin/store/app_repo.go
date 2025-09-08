package adminstore

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"strconv"

	"gorm.io/gorm"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/resource"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type AppRepo struct {
	*resourceRepo[admin.App, admin.AppAttrs, db.App]
}

func NewAppRepo(d *db.DB) *AppRepo {
	return &AppRepo{
		resourceRepo: &resourceRepo[admin.App, admin.AppAttrs, db.App]{
			db:           d,
			mapper:       appMapper{db: d},
			associations: []string{"User"},
		},
	}
}

func (r *AppRepo) ListOwnedByUser(ctx context.Context, userID int64, _ map[string][]string) (*resource.Collection[admin.App], error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}, nil)
}

func (r *AppRepo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.App, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	})
}

type appMapper struct {
	db *db.DB
}

//lint:ignore U1000 this method is used by generic struct
func (m appMapper) dbModel(a *admin.AppAttrs, id int64) *db.App {
	packageName := sql.NullString{}
	if a.PackageName != "" {
		packageName.String = a.PackageName
		packageName.Valid = true
	}

	appKey := sql.NullString{}
	if id == 0 {
		appKey.String, _ = m.generateAppKey()
		appKey.Valid = true
	}

	storeID := sql.NullString{}
	if a.StoreID != "" {
		storeID.String = a.StoreID
		storeID.Valid = true
	}

	storeURL := sql.NullString{}
	if a.StoreURL != "" {
		storeURL.String = a.StoreURL
		storeURL.Valid = true
	}

	return &db.App{
		ID:          id,
		UserID:      a.UserID,
		PlatformID:  dbPlatformID(a.PlatformID),
		HumanName:   a.HumanName,
		PackageName: packageName,
		AppKey:      appKey,
		StoreID:     storeID,
		StoreURL:    storeURL,
		Categories:  a.Categories,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m appMapper) resource(a *db.App) admin.App {
	return admin.App{
		ID:        a.ID,
		PublicUID: strconv.FormatInt(a.PublicUID.Int64, 10),
		AppAttrs:  m.resourceAttrs(a),
		User:      userMapper{}.resource(&a.User),
	}
}

func (m appMapper) resourceAttrs(a *db.App) admin.AppAttrs {
	return admin.AppAttrs{
		UserID:      a.UserID,
		PlatformID:  platformID(a.PlatformID),
		HumanName:   a.HumanName,
		PackageName: a.PackageName.String,
		AppKey:      a.AppKey.String,
		StoreID:     a.StoreID.String,
		StoreURL:    a.StoreURL.String,
		Categories:  a.Categories,
	}
}

func (m appMapper) generateAppKey() (string, error) {
	keyBytes := make([]byte, 24)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(keyBytes), nil
}
