package adminstore

import (
	"context"
	"encoding/json"
	"strconv"

	"gorm.io/gorm"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/resource"
	"github.com/bidon-io/bidon-backend/internal/db"
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

func (r *AppDemandProfileRepo) List(ctx context.Context, qParams map[string][]string) (*resource.Collection[admin.AppDemandProfile], error) {
	filters := queryToAppDemandProfilesFilters(qParams)
	pgn := PaginationFromQueryParams[db.AppDemandProfile](qParams)
	return r.list(ctx, filters.apply, pgn)
}

func (r *AppDemandProfileRepo) ListOwnedByUser(ctx context.Context, userID int64, qParams map[string][]string) (*resource.Collection[admin.AppDemandProfile], error) {
	filters := queryToAppDemandProfilesFilters(qParams)
	filters.UserID = userID
	pgn := PaginationFromQueryParams[db.AppDemandProfile](qParams)
	return r.list(ctx, filters.apply, pgn)
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

type appDemandProfilesFilters struct {
	UserID         int64
	AppID          int64
	AccountID      int64
	DemandSourceID int64
}

func (f *appDemandProfilesFilters) apply(db *gorm.DB) *gorm.DB {
	if f.UserID != 0 {
		db = db.Joins("INNER JOIN apps ON apps.id = app_demand_profiles.app_id").Where("apps.user_id = ?", f.UserID)
	}
	if f.AppID != 0 {
		db = db.Where("app_id = ?", f.AppID)
	}
	if f.AccountID != 0 {
		db = db.Where("account_id = ?", f.AccountID)
	}
	if f.DemandSourceID != 0 {
		db = db.Where("demand_source_id = ?", f.DemandSourceID)
	}
	return db
}

func queryToAppDemandProfilesFilters(qParams map[string][]string) appDemandProfilesFilters {
	filters := appDemandProfilesFilters{}
	if v, ok := qParams["user_id"]; ok {
		filters.UserID, _ = strconv.ParseInt(v[0], 10, 64)
	}
	if v, ok := qParams["app_id"]; ok {
		filters.AppID, _ = strconv.ParseInt(v[0], 10, 64)
	}
	if v, ok := qParams["account_id"]; ok {
		filters.AccountID, _ = strconv.ParseInt(v[0], 10, 64)
	}
	if v, ok := qParams["demand_source_id"]; ok {
		filters.DemandSourceID, _ = strconv.ParseInt(v[0], 10, 64)
	}
	return filters
}
