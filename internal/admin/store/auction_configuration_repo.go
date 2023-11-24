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

type AuctionConfigurationRepo struct {
	*resourceRepo[admin.AuctionConfiguration, admin.AuctionConfigurationAttrs, db.AuctionConfiguration]
}

func NewAuctionConfigurationRepo(d *db.DB) *AuctionConfigurationRepo {
	return &AuctionConfigurationRepo{
		resourceRepo: &resourceRepo[admin.AuctionConfiguration, admin.AuctionConfigurationAttrs, db.AuctionConfiguration]{
			db:           d,
			mapper:       auctionConfigurationMapper{db: d},
			associations: []string{"App", "Segment"},
		},
	}
}

func (r *AuctionConfigurationRepo) ListOwnedByUser(ctx context.Context, userID int64) ([]admin.AuctionConfiguration, error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Select("user_id").Where(map[string]any{"user_id": userID}))
	})
}

func (r *AuctionConfigurationRepo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.AuctionConfiguration, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Select("user_id").Where(map[string]any{"user_id": userID}))
	})
}

type auctionConfigurationMapper struct {
	db *db.DB
}

//lint:ignore U1000 this method is used by generic struct
func (m auctionConfigurationMapper) dbModel(c *admin.AuctionConfigurationAttrs, id int64) *db.AuctionConfiguration {
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

	return &db.AuctionConfiguration{
		ID:                       id,
		Name:                     name,
		AppID:                    c.AppID,
		AdType:                   db.AdTypeFromDomain(c.AdType),
		Rounds:                   c.Rounds,
		Pricefloor:               c.Pricefloor,
		SegmentID:                &segmentID,
		ExternalWinNotifications: c.ExternalWinNotifications,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m auctionConfigurationMapper) resource(c *db.AuctionConfiguration) admin.AuctionConfiguration {
	var segment *admin.Segment
	if c.Segment != nil {
		segment = &admin.Segment{
			ID:           c.Segment.ID,
			SegmentAttrs: segmentMapper{}.resourceAttrs(c.Segment),
		}
	}

	return admin.AuctionConfiguration{
		ID:                        c.ID,
		PublicUID:                 strconv.FormatInt(c.PublicUID.Int64, 10),
		AuctionKey:                strings.ToUpper(big.NewInt(c.PublicUID.Int64).Text(32)),
		AuctionConfigurationAttrs: m.resourceAttrs(c),
		App: admin.App{
			ID:       c.App.ID,
			AppAttrs: appMapper{}.resourceAttrs(&c.App),
		},
		Segment: segment,
	}
}

func (m auctionConfigurationMapper) resourceAttrs(c *db.AuctionConfiguration) admin.AuctionConfigurationAttrs {
	var segmentID *int64
	if c.SegmentID != nil && c.SegmentID.Valid {
		segmentID = &c.SegmentID.Int64
	} else {
		segmentID = nil
	}

	return admin.AuctionConfigurationAttrs{
		Name:                     c.Name.String,
		AppID:                    c.AppID,
		AdType:                   c.AdType.Domain(),
		Rounds:                   c.Rounds,
		Pricefloor:               c.Pricefloor,
		SegmentID:                segmentID,
		ExternalWinNotifications: c.ExternalWinNotifications,
	}
}
