package adminstore

import (
	"context"
	"database/sql"
	"math/big"
	"strconv"
	"strings"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
	"gorm.io/gorm"
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

func (r *AuctionConfigurationV2Repo) List(ctx context.Context) ([]admin.AuctionConfigurationV2, error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("settings->>'v2' = ?", "true")
	})
}

func (r *AuctionConfigurationV2Repo) ListOwnedByUser(ctx context.Context, userID int64) ([]admin.AuctionConfigurationV2, error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Select("user_id").Where(map[string]any{"user_id": userID}).Where("settings->>'v2' = ?", "true"))
	})
}

func (r *AuctionConfigurationV2Repo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.AuctionConfigurationV2, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Select("user_id").Where(map[string]any{"user_id": userID}).Where("settings->>'v2' = ?", "true"))
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
		ExternalWinNotifications: c.ExternalWinNotifications,
		Demands:                  db.StringArrayToAdapterKeys(&c.Demands),
		Bidding:                  db.StringArrayToAdapterKeys(&c.Bidding),
		AdUnitIDs:                c.AdUnitIds,
		Timeout:                  c.Timeout,
		Settings:                 c.Settings,
	}
}
