package adminstore

import (
	"context"
	"encoding/json"

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
			mapper:       demandSourceAccountMapper{},
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

type demandSourceAccountMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m demandSourceAccountMapper) dbModel(a *admin.DemandSourceAccountAttrs, id int64) *db.DemandSourceAccount {
	extra, _ := json.Marshal(a.Extra)

	return &db.DemandSourceAccount{
		Model:          db.Model{ID: id},
		DemandSourceID: a.DemandSourceID,
		UserID:         a.UserID,
		Type:           a.Type,
		Extra:          extra,
		IsBidding:      a.IsBidding,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m demandSourceAccountMapper) resource(a *db.DemandSourceAccount) admin.DemandSourceAccount {
	return admin.DemandSourceAccount{
		ID:                       a.ID,
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
		Type:           a.Type,
		DemandSourceID: a.DemandSourceID,
		IsBidding:      a.IsBidding,
		Extra:          extra,
	}
}
