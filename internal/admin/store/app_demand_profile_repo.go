package adminstore

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/admin/resource"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
	"gorm.io/gorm"
)

type AppDemandProfileRepo struct {
	*resourceRepo[admin.AppDemandProfile, admin.AppDemandProfileAttrs, db.AppDemandProfile]
}

func NewAppDemandProfileRepo(d *db.DB) *AppDemandProfileRepo {
	return &AppDemandProfileRepo{
		resourceRepo: &resourceRepo[admin.AppDemandProfile, admin.AppDemandProfileAttrs, db.AppDemandProfile]{
			db:           d,
			mapper:       appDemandProfileMapper{db: d},
			associations: []string{"App", "Account", "DemandSource"},
		},
	}
}

func (r *AppDemandProfileRepo) ListOwnedByUser(ctx context.Context, userID int64, _ map[string][]string) (*resource.Collection[admin.AppDemandProfile], error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Table("App").Where(map[string]any{"user_id": userID}))
	}, nil)
}

func (r *AppDemandProfileRepo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.AppDemandProfile, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Table("App").Where(map[string]any{"user_id": userID}))
	})
}

type appDemandProfileMapper struct {
	db *db.DB
}

//lint:ignore U1000 this method is used by generic struct
func (m appDemandProfileMapper) dbModel(attrs *admin.AppDemandProfileAttrs, id int64) *db.AppDemandProfile {
	data, _ := json.Marshal(attrs.Data)

	return &db.AppDemandProfile{
		ID:             id,
		AppID:          attrs.AppID,
		AccountType:    attrs.AccountType,
		AccountID:      attrs.AccountID,
		DemandSourceID: attrs.DemandSourceID,
		Data:           data,
		Enabled:        attrs.Enabled,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m appDemandProfileMapper) resource(p *db.AppDemandProfile) admin.AppDemandProfile {
	return admin.AppDemandProfile{
		ID:                    p.ID,
		PublicUID:             strconv.FormatInt(p.PublicUID.Int64, 10),
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
	var data map[string]any
	_ = json.Unmarshal(p.Data, &data)

	return admin.AppDemandProfileAttrs{
		AppID:          p.AppID,
		DemandSourceID: p.DemandSourceID,
		AccountID:      p.AccountID,
		Data:           data,
		AccountType:    p.AccountType,
		Enabled:        p.Enabled,
	}
}
