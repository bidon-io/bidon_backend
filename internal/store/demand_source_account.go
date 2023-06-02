package store

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
)

type DemandSourceAccountRepo = resourceRepo[admin.DemandSourceAccount, admin.DemandSourceAccountAttrs, demandSourceAccount, *demandSourceAccount]

type demandSourceAccount struct {
	Model
	DemandSourceID int64          `gorm:"column:demand_source_id;type:bigint;not null"`
	UserID         int64          `gorm:"column:user_id;type:bigint;not null"`
	Type           string         `gorm:"column:type;type:varchar;not null"`
	Extra          map[string]any `gorm:"column:extra;type:jsonb;default:'{}';serializer:json"`
	IsBidding      *bool          `gorm:"column:bidding;type:boolean;default:false"`
	IsDefault      sql.NullBool   `gorm:"column:is_default;type:boolean"`
}

//lint:ignore U1000 this method is used by generic struct
func (a *demandSourceAccount) initFromResourceAttrs(attrs *admin.DemandSourceAccountAttrs) {
	a.DemandSourceID = attrs.DemandSourceID
	a.UserID = attrs.UserID
	a.Type = attrs.Type
	a.Extra = attrs.Extra
	a.IsBidding = attrs.IsBidding
}

//lint:ignore U1000 this method is used by generic struct
func (a *demandSourceAccount) toResource() admin.DemandSourceAccount {
	return admin.DemandSourceAccount{
		ID: a.ID,
		DemandSourceAccountAttrs: admin.DemandSourceAccountAttrs{
			UserID:         a.UserID,
			Type:           a.Type,
			DemandSourceID: a.DemandSourceID,
			IsBidding:      a.IsBidding,
			Extra:          a.Extra,
		},
	}
}
