package adminstore

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type LineItemRepo struct {
	*resourceRepo[admin.LineItem, admin.LineItemAttrs, db.LineItem]
}

func NewLineItemRepo(d *db.DB) *LineItemRepo {
	return &LineItemRepo{
		resourceRepo: &resourceRepo[admin.LineItem, admin.LineItemAttrs, db.LineItem]{
			db:           d,
			mapper:       lineItemMapper{db: d},
			associations: []string{"App", "Account"},
		},
	}
}

func (r *LineItemRepo) ListOwnedByUser(ctx context.Context, userID int64) ([]admin.LineItem, error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Select("user_id").Where(map[string]any{"user_id": userID}))
	})
}

func (r *LineItemRepo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.LineItem, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Select("user_id").Where(map[string]any{"user_id": userID}))
	})
}

func (r *LineItemRepo) CreateMany(ctx context.Context, items []admin.LineItemAttrs) error {
	dbItems := make([]*db.LineItem, len(items))
	for i := range items {
		dbItems[i] = r.mapper.dbModel(&items[i], 0)
	}
	return r.db.WithContext(ctx).Create(&dbItems).Error
}

type lineItemMapper struct {
	db *db.DB
}

//lint:ignore U1000 this method is used by generic struct
func (m lineItemMapper) dbModel(i *admin.LineItemAttrs, id int64) *db.LineItem {
	bidFloor := decimal.NullDecimal{}
	if i.BidFloor != nil {
		bidFloor.Decimal = *i.BidFloor
		bidFloor.Valid = true
	}

	format := sql.NullString{}
	if i.Format != nil {
		format.String = string(*i.Format)
		format.Valid = true
	}

	var isBidding sql.NullBool
	if i.IsBidding != nil {
		isBidding.Bool = *i.IsBidding
		isBidding.Valid = true
	}

	// TODO: remove this hack
	// sets code to one of the following values: ad_unit_id, zone_id, placement_id, slot_id, slot_uuid
	if i.Code == nil {
		keysToCheck := []string{"zone_id", "placement_id", "slot_id", "slot_uuid", "ad_unit_id", "unit_id", "spot_id"}

		for _, key := range keysToCheck {
			if value, ok := i.Extra[key].(string); ok {
				i.Code = &value
				break
			}
		}
	}

	return &db.LineItem{
		ID:          id,
		AppID:       i.AppID,
		AccountType: i.AccountType,
		AccountID:   i.AccountID,
		HumanName:   i.HumanName,
		Code:        i.Code,
		BidFloor:    bidFloor,
		AdType:      db.AdTypeFromDomain(i.AdType),
		Extra:       i.Extra,
		Format:      format,
		IsBidding:   isBidding,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m lineItemMapper) resource(i *db.LineItem) admin.LineItem {
	return admin.LineItem{
		ID:            i.ID,
		PublicUID:     strconv.FormatInt(i.PublicUID.Int64, 10),
		LineItemAttrs: m.resourceAttrs(i),
		App: admin.App{
			ID:       i.AppID,
			AppAttrs: appMapper{}.resourceAttrs(&i.App),
		},
		Account: admin.DemandSourceAccount{
			ID:                       i.AccountID,
			DemandSourceAccountAttrs: demandSourceAccountMapper{}.resourceAttrs(&i.Account),
		},
	}
}

func (m lineItemMapper) resourceAttrs(i *db.LineItem) admin.LineItemAttrs {
	format := ad.Format(i.Format.String)
	return admin.LineItemAttrs{
		HumanName:   i.HumanName,
		AppID:       i.AppID,
		BidFloor:    &i.BidFloor.Decimal,
		AdType:      i.AdType.Domain(),
		Format:      &format,
		AccountID:   i.AccountID,
		AccountType: i.AccountType,
		Code:        i.Code,
		IsBidding:   &i.IsBidding.Bool,
		Extra:       i.Extra,
	}
}
