package adminstore

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type AppDemandProfileRepo = resourceRepo[admin.AppDemandProfile, admin.AppDemandProfileAttrs, db.AppDemandProfile]

func NewAppDemandProfileRepo(db *db.DB) *AppDemandProfileRepo {
	return &AppDemandProfileRepo{
		db:           db,
		mapper:       appDemandProfileMapper{},
		associations: []string{"App", "Account", "DemandSource"},
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
		ID:                    p.ID,
		AppDemandProfileAttrs: m.resourceAttrs(p),
		App: admin.App{
			ID:       p.AppID,
			AppAttrs: appMapper{}.resourceAttrs(&p.App),
		},
		Account: admin.DemandSourceAccount{
			ID:                       p.AccountID,
			DemandSourceAccountAttrs: demandSourceAccountMapper{}.resourceAttrs(&p.Account),
		},
		DemandSource: admin.DemandSource{
			ID:                p.DemandSourceID,
			DemandSourceAttrs: demandSourceMapper{}.resourceAttrs(&p.DemandSource),
		},
	}
}

func (m appDemandProfileMapper) resourceAttrs(p *db.AppDemandProfile) admin.AppDemandProfileAttrs {
	return admin.AppDemandProfileAttrs{
		AppID:          p.AppID,
		DemandSourceID: p.DemandSourceID,
		AccountID:      p.AccountID,
		Data:           p.Data,
		AccountType:    p.AccountType,
	}
}
