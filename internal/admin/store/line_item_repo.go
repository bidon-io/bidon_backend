package store

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/shopspring/decimal"
)

type LineItemRepo = resourceRepo[admin.LineItem, admin.LineItemAttrs, db.LineItem]

func NewLineItemRepo(db *db.DB) *LineItemRepo {
	return &LineItemRepo{
		db:           db,
		mapper:       lineItemMapper{},
		associations: []string{"App", "Account"},
	}
}

type lineItemMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m lineItemMapper) dbModel(i *admin.LineItemAttrs, id int64) *db.LineItem {
	var bidFloor decimal.NullDecimal
	if i.BidFloor != nil {
		bidFloor.Decimal = *i.BidFloor
		bidFloor.Valid = true
	}

	var format sql.NullString
	if i.Format != nil {
		format.String = string(*i.Format)
		format.Valid = true
	}

	return &db.LineItem{
		Model:       db.Model{ID: id},
		AppID:       i.AppID,
		AccountType: i.AccountType,
		AccountID:   i.AccountID,
		HumanName:   i.HumanName,
		Code:        i.Code,
		BidFloor:    bidFloor,
		AdType:      db.AdTypeFromDomain(i.AdType),
		Extra:       i.Extra,
		Format:      format,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m lineItemMapper) resource(i *db.LineItem) admin.LineItem {
	return admin.LineItem{
		ID:            i.ID,
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
		Extra:       i.Extra,
	}
}
