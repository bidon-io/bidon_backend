package store

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type AuctionConfigurationRepo = resourceRepo[admin.AuctionConfiguration, admin.AuctionConfigurationAttrs, db.AuctionConfiguration]

func NewAuctionConfigurationRepo(db *db.DB) *AuctionConfigurationRepo {
	return &AuctionConfigurationRepo{
		db:     db,
		mapper: auctionConfigurationMapper{},
	}
}

type auctionConfigurationMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m auctionConfigurationMapper) dbModel(c *admin.AuctionConfigurationAttrs) *db.AuctionConfiguration {
	name := sql.NullString{}
	if c.Name != "" {
		name.String = c.Name
		name.Valid = true
	}

	return &db.AuctionConfiguration{
		Name:       name,
		AppID:      c.AppID,
		AdType:     db.AdTypeFromDomain(c.AdType),
		Rounds:     c.Rounds,
		Pricefloor: c.Pricefloor,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m auctionConfigurationMapper) resource(c *db.AuctionConfiguration) admin.AuctionConfiguration {
	return admin.AuctionConfiguration{
		ID: c.ID,
		AuctionConfigurationAttrs: admin.AuctionConfigurationAttrs{
			Name:       c.Name.String,
			AppID:      c.AppID,
			AdType:     c.AdType.Domain(),
			Rounds:     c.Rounds,
			Pricefloor: c.Pricefloor,
		},
	}
}
