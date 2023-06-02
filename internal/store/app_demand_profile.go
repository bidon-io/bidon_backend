package store

import "github.com/bidon-io/bidon-backend/internal/admin"

type AppDemandProfileRepo = resourceRepo[admin.AppDemandProfile, admin.AppDemandProfileAttrs, appDemandProfile, *appDemandProfile]

type appDemandProfile struct {
	Model
	AppID          int64          `gorm:"column:app_id;type:bigint;not null"`
	AccountType    string         `gorm:"column:account_type;type:varchar;not null"`
	AccountID      int64          `gorm:"column:account_id;type:bigint;not null"`
	DemandSourceID int64          `gorm:"column:demand_source_id;type:bigint;not null"`
	Data           map[string]any `gorm:"column:data;type:jsonb;default:'{}';serializer:json"`
}

//lint:ignore U1000 this method is used by generic struct
func (p *appDemandProfile) initFromResourceAttrs(attrs *admin.AppDemandProfileAttrs) {
	p.AppID = attrs.AppID
	p.AccountType = attrs.AccountType
	p.AccountID = attrs.AccountID
	p.DemandSourceID = attrs.DemandSourceID
	p.Data = attrs.Data
}

//lint:ignore U1000 this method is used by generic struct
func (p *appDemandProfile) toResource() admin.AppDemandProfile {
	return admin.AppDemandProfile{
		ID: p.ID,
		AppDemandProfileAttrs: admin.AppDemandProfileAttrs{
			AppID:          p.AppID,
			DemandSourceID: p.DemandSourceID,
			AccountID:      p.AccountID,
			Data:           p.Data,
			AccountType:    p.AccountType,
		},
	}
}
