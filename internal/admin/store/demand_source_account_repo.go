package adminstore

import (
	"encoding/json"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type DemandSourceAccountRepo = resourceRepo[admin.DemandSourceAccount, admin.DemandSourceAccountAttrs, db.DemandSourceAccount]

func NewDemandSourceAccountRepo(db *db.DB) *DemandSourceAccountRepo {
	return &DemandSourceAccountRepo{
		db:           db,
		mapper:       demandSourceAccountMapper{},
		associations: []string{"User", "DemandSource"},
	}
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
		User: admin.User{
			ID:        a.User.ID,
			UserAttrs: userMapper{}.resourceAttrs(&a.User),
		},
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
