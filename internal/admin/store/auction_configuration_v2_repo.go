package adminstore

import (
	"context"
	"database/sql"
	"math/big"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/resource"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type AuctionConfigurationV2Repo struct {
	*resourceRepo[admin.AuctionConfigurationV2, admin.AuctionConfigurationV2Attrs, db.AuctionConfiguration]
}

func NewAuctionConfigurationV2Repo(d *db.DB) *AuctionConfigurationV2Repo {
	return &AuctionConfigurationV2Repo{
		resourceRepo: &resourceRepo[admin.AuctionConfigurationV2, admin.AuctionConfigurationV2Attrs, db.AuctionConfiguration]{
			db:           d,
			mapper:       auctionConfigurationV2Mapper{db: d},
			associations: []string{"App", "Segment"},
		},
	}
}

func (r *AuctionConfigurationV2Repo) List(ctx context.Context, qParams map[string][]string) (*resource.Collection[admin.AuctionConfigurationV2], error) {
	filters := queryToAuctionConfigurationFilters(qParams)
	pgn := PaginationFromQueryParams[db.AuctionConfiguration](qParams)

	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(filters.apply, r.v2Scope)
	}, pgn)
}

func (r *AuctionConfigurationV2Repo) ListOwnedByUser(ctx context.Context, userID int64, qParams map[string][]string) (*resource.Collection[admin.AuctionConfigurationV2], error) {
	filters := queryToAuctionConfigurationFilters(qParams)
	filters.UserID = userID
	pgn := PaginationFromQueryParams[db.AuctionConfiguration](qParams)

	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(filters.apply, r.v2Scope)
	}, pgn)
}

func (r *AuctionConfigurationV2Repo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.AuctionConfigurationV2, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Table("App").Where(map[string]any{"user_id": userID}).Where("settings->>'v2' = ?", "true"))
	})
}

type auctionConfigurationV2Mapper struct {
	db *db.DB
}

//lint:ignore U1000 this method is used by generic struct
func (m auctionConfigurationV2Mapper) dbModel(c *admin.AuctionConfigurationV2Attrs, id int64) *db.AuctionConfiguration {
	name := sql.NullString{}
	if c.Name != "" {
		name.String = c.Name
		name.Valid = true
	}
	segmentID := sql.NullInt64{}
	if c.SegmentID != nil {
		segmentID.Int64 = *c.SegmentID
		segmentID.Valid = true
	}

	model := &db.AuctionConfiguration{
		ID:                       id,
		Name:                     name,
		AppID:                    c.AppID,
		AdType:                   db.AdTypeFromDomain(c.AdType),
		Pricefloor:               c.Pricefloor,
		SegmentID:                &segmentID,
		IsDefault:                c.IsDefault,
		ExternalWinNotifications: c.ExternalWinNotifications,
		Demands:                  db.AdapterKeysToStringArray(c.Demands),
		Bidding:                  db.AdapterKeysToStringArray(c.Bidding),
		AdUnitIds:                c.AdUnitIDs,
		Timeout:                  c.Timeout,
		Settings:                 c.Settings,
	}

	if id == 0 {
		if c.Settings != nil {
			model.Settings["v2"] = true
		} else {
			model.Settings = map[string]any{"v2": true}
		}
	}

	return model
}

//lint:ignore U1000 this method is used by generic struct
func (m auctionConfigurationV2Mapper) resource(c *db.AuctionConfiguration) admin.AuctionConfigurationV2 {
	var segment *admin.Segment
	if c.Segment != nil {
		segment = &admin.Segment{
			ID:           c.Segment.ID,
			SegmentAttrs: segmentMapper{}.resourceAttrs(c.Segment),
		}
	}

	return admin.AuctionConfigurationV2{
		ID:                          c.ID,
		PublicUID:                   strconv.FormatInt(c.PublicUID.Int64, 10),
		AuctionKey:                  strings.ToUpper(big.NewInt(c.PublicUID.Int64).Text(32)),
		AuctionConfigurationV2Attrs: m.resourceAttrs(c),
		App: admin.App{
			ID:       c.App.ID,
			AppAttrs: appMapper{}.resourceAttrs(&c.App),
		},
		Segment: segment,
	}
}

func (m auctionConfigurationV2Mapper) resourceAttrs(c *db.AuctionConfiguration) admin.AuctionConfigurationV2Attrs {
	var segmentID *int64
	if c.SegmentID != nil && c.SegmentID.Valid {
		segmentID = &c.SegmentID.Int64
	} else {
		segmentID = nil
	}

	return admin.AuctionConfigurationV2Attrs{
		Name:                     c.Name.String,
		AppID:                    c.AppID,
		AdType:                   c.AdType.Domain(),
		Pricefloor:               c.Pricefloor,
		SegmentID:                segmentID,
		IsDefault:                c.IsDefault,
		ExternalWinNotifications: c.ExternalWinNotifications,
		Demands:                  db.StringArrayToAdapterKeys(&c.Demands),
		Bidding:                  db.StringArrayToAdapterKeys(&c.Bidding),
		AdUnitIDs:                c.AdUnitIds,
		Timeout:                  c.Timeout,
		Settings:                 c.Settings,
	}
}

type AuctionConfigurationFilters struct {
	UserID    int64
	AppID     int64
	AdType    db.AdType
	SegmentID int64
	IsDefault *bool
	Name      string
}

func (f *AuctionConfigurationFilters) apply(db *gorm.DB) *gorm.DB {
	if f.UserID != 0 {
		db = db.Joins("INNER JOIN apps ON apps.id = auction_configurations.app_id").Where("apps.user_id = ?", f.UserID)
	}
	if f.AppID != 0 {
		db = db.Where("app_id = ?", f.AppID)
	}
	if f.AdType != 0 {
		db = db.Where("ad_type = ?", f.AdType)
	}
	if f.SegmentID != 0 {
		db = db.Where("segment_id = ?", f.SegmentID)
	}
	if f.IsDefault != nil {
		if *f.IsDefault {
			db = db.Where("is_default = true")
		} else {
			db = db.Where("is_default = false OR is_default IS NULL")
		}
	}
	if f.Name != "" {
		db = db.Where("name ILIKE ?", "%"+f.Name+"%")
	}
	return db
}

func (r *AuctionConfigurationV2Repo) v2Scope(db *gorm.DB) *gorm.DB {
	return db.Where("settings->>'v2' = ?", "true")
}

func queryToAuctionConfigurationFilters(qParams map[string][]string) AuctionConfigurationFilters {
	filters := AuctionConfigurationFilters{}
	if v, ok := qParams["user_id"]; ok {
		filters.UserID, _ = strconv.ParseInt(v[0], 10, 64)
	}
	if v, ok := qParams["app_id"]; ok {
		filters.AppID, _ = strconv.ParseInt(v[0], 10, 64)
	}
	if v, ok := qParams["ad_type"]; ok {
		dbAdType := db.AdTypeFromDomain(ad.Type(v[0]))
		filters.AdType = dbAdType
	}
	if v, ok := qParams["segment_id"]; ok {
		filters.SegmentID, _ = strconv.ParseInt(v[0], 10, 64)
	}
	if v, ok := qParams["is_default"]; ok {
		b := v[0] == "true"
		filters.IsDefault = &b
	}
	if v, ok := qParams["name"]; ok {
		filters.Name = v[0]
	}
	return filters
}
