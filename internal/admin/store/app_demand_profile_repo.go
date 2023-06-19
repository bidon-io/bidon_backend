package store

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type AppDemandProfileRepo = resourceRepo[admin.AppDemandProfile, admin.AppDemandProfileAttrs, db.AppDemandProfile]

func NewAppDemandProfileRepo(db *db.DB) *AppDemandProfileRepo {
	return &AppDemandProfileRepo{
		db:     db,
		mapper: appDemandProfileMapper{},
	}
}

type appDemandProfileMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m appDemandProfileMapper) dbModel(attrs *admin.AppDemandProfileAttrs, id int64) *db.AppDemandProfile {
	return &db.AppDemandProfile{
		Model:          db.Model{ID: id},
		AppID:          attrs.AppID,
		AccountType:    attrs.AccountType,
		AccountID:      attrs.AccountID,
		DemandSourceID: attrs.DemandSourceID,
		Data:           attrs.Data,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m appDemandProfileMapper) resource(p *db.AppDemandProfile) admin.AppDemandProfile {
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
