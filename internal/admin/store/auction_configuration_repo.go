package store

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type AuctionConfigurationRepo = resourceRepo[admin.AuctionConfiguration, admin.AuctionConfigurationAttrs, db.AuctionConfiguration]

func NewAuctionConfigurationRepo(db *db.DB) *AuctionConfigurationRepo {
	return &AuctionConfigurationRepo{
		db:           db,
		mapper:       auctionConfigurationMapper{},
		associations: []string{"App", "Segment"},
	}
}

type auctionConfigurationMapper struct{}

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
		Model:                    db.Model{ID: id},
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
		SegmentID:                segmentID,
		ExternalWinNotifications: c.ExternalWinNotifications,
	}
}
