package store

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/shopspring/decimal"
)

type LineItemRepo = resourceRepo[admin.LineItem, admin.LineItemAttrs, lineItem, *lineItem]

type lineItem struct {
	Model
	AppID       int64               `gorm:"column:app_id;type:bigint;not null"`
	AccountType string              `gorm:"column:account_type;type:varchar;not null"`
	AccountID   int64               `gorm:"column:account_id;type:bigint;not null"`
	HumanName   string              `gorm:"column:human_name;type:varchar;not null"`
	Code        *string             `gorm:"column:code;type:varchar;not null"`
	BidFloor    decimal.NullDecimal `gorm:"column:bid_floor;type:numeric"`
	AdType      int32               `gorm:"column:ad_type;type:integer;not null"`
	Extra       map[string]any      `gorm:"column:extra;type:jsonb;default:'{}';serializer:json"`
	Width       int32               `gorm:"column:width;type:integer;default:0;not null"`
	Height      int32               `gorm:"column:height;type:integer;default:0;not null"`
	Format      sql.NullString      `gorm:"column:format;type:varchar"`
}

//lint:ignore U1000 this method is used by generic struct
func (i *lineItem) initFromResourceAttrs(attrs *admin.LineItemAttrs) {
	i.AppID = attrs.AppID
	i.AccountType = attrs.AccountType
	i.AccountID = attrs.AccountID
	i.HumanName = attrs.HumanName
	i.Code = attrs.Code
	if attrs.BidFloor != nil {
		i.BidFloor.Decimal = *attrs.BidFloor
		i.BidFloor.Valid = true
	}
	i.AdType = dbAdType(attrs.AdType)
	i.Extra = attrs.Extra
	if attrs.Format != nil {
		i.Format.String = string(*attrs.Format)
		i.Format.Valid = true
	}
}

//lint:ignore U1000 this method is used by generic struct
func (i *lineItem) toResource() admin.LineItem {
	format := admin.LineItemFormat(i.Format.String)
	return admin.LineItem{
		ID: i.ID,
		LineItemAttrs: admin.LineItemAttrs{
			HumanName:   i.HumanName,
			AppID:       i.AppID,
			BidFloor:    &i.BidFloor.Decimal,
			AdType:      adType(i.AdType),
			Format:      &format,
			AccountID:   i.AccountID,
			AccountType: i.AccountType,
			Code:        i.Code,
			Extra:       i.Extra,
		},
	}
}
