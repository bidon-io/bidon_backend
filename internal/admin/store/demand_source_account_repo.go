package adminstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
	"gorm.io/gorm"
)

type DemandSourceAccountRepo struct {
	*resourceRepo[admin.DemandSourceAccount, admin.DemandSourceAccountAttrs, db.DemandSourceAccount]
}

func NewDemandSourceAccountRepo(d *db.DB) *DemandSourceAccountRepo {
	return &DemandSourceAccountRepo{
		resourceRepo: &resourceRepo[admin.DemandSourceAccount, admin.DemandSourceAccountAttrs, db.DemandSourceAccount]{
			db:           d,
			mapper:       demandSourceAccountMapper{db: d},
			associations: []string{"User", "DemandSource"},
		},
	}
}

func (r DemandSourceAccountRepo) ListOwnedByUser(ctx context.Context, userID int64) ([]admin.DemandSourceAccount, error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	})
}

func (r DemandSourceAccountRepo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.DemandSourceAccount, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	})
}

func (r DemandSourceAccountRepo) ListOwnedByUserOrShared(ctx context.Context, userID int64) ([]admin.DemandSourceAccount, error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id IN ?", []int64{userID, 0, 1})
	})
}

func (r DemandSourceAccountRepo) FindOwnedByUserOrShared(ctx context.Context, userID int64, id int64) (*admin.DemandSourceAccount, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id IN ?", []int64{userID, 0, 1})
	})
}

type demandSourceAccountMapper struct {
	db *db.DB
}

//lint:ignore U1000 this method is used by generic struct
func (m demandSourceAccountMapper) dbModel(a *admin.DemandSourceAccountAttrs, id int64) *db.DemandSourceAccount {
	extra, _ := json.Marshal(a.Extra)

	var label sql.NullString
	if a.Label != "" {
		label.String = a.Label
		label.Valid = true
	}

	var isBidding sql.NullBool
	if a.IsBidding != nil {
		isBidding.Bool = *a.IsBidding
		isBidding.Valid = true
	}

	return &db.DemandSourceAccount{
		ID:             id,
		DemandSourceID: a.DemandSourceID,
		Label:          label,
		UserID:         a.UserID,
		Type:           a.Type,
		Extra:          extra,
		IsBidding:      isBidding,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m demandSourceAccountMapper) resource(a *db.DemandSourceAccount) admin.DemandSourceAccount {
	return admin.DemandSourceAccount{
		ID:                       a.ID,
		PublicUID:                strconv.FormatInt(a.PublicUID.Int64, 10),
		DemandSourceAccountAttrs: m.resourceAttrs(a),
		User:                     userMapper{}.resource(&a.User),
		DemandSource: admin.DemandSource{
			ID:                a.DemandSource.ID,
			DemandSourceAttrs: demandSourceMapper{}.resourceAttrs(&a.DemandSource),
		},
	}
}

func (m demandSourceAccountMapper) resourceAttrs(a *db.DemandSourceAccount) admin.DemandSourceAccountAttrs {
	var extra map[string]any
	_ = json.Unmarshal(a.Extra, &extra)

	return admin.DemandSourceAccountAttrs{
		UserID:         a.UserID,
		Label:          a.Label.String,
		Type:           a.Type,
		DemandSourceID: a.DemandSourceID,
		IsBidding:      &a.IsBidding.Bool,
		Extra:          extra,
	}
}
