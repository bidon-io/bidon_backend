package adminstore

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type DemandSourceRepo = resourceRepo[admin.DemandSource, admin.DemandSourceAttrs, db.DemandSource]

func NewDemandSourceRepo(db *db.DB) *DemandSourceRepo {
	return &DemandSourceRepo{
		db:           db,
		mapper:       demandSourceMapper{},
		associations: []string{},
	}
}

type demandSourceMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m demandSourceMapper) dbModel(s *admin.DemandSourceAttrs, id int64) *db.DemandSource {
	return &db.DemandSource{
		Model:     db.Model{ID: id},
		APIKey:    s.ApiKey,
		HumanName: s.HumanName,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m demandSourceMapper) resource(s *db.DemandSource) admin.DemandSource {
	return admin.DemandSource{
		ID:                s.ID,
		DemandSourceAttrs: m.resourceAttrs(s),
	}
}

func (m demandSourceMapper) resourceAttrs(s *db.DemandSource) admin.DemandSourceAttrs {
	return admin.DemandSourceAttrs{
		ApiKey:    s.APIKey,
		HumanName: s.HumanName,
	}
}
