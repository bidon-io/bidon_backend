package store

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
)

type AuctionConfigurationRepo = resourceRepo[admin.AuctionConfiguration, admin.AuctionConfigurationAttrs, auctionConfiguration, *auctionConfiguration]

type auctionConfiguration struct {
	Model
	Name       sql.NullString                    `gorm:"column:name;type:varchar"`
	AppID      int64                             `gorm:"column:app_id;type:bigint;not null"`
	AdType     int32                             `gorm:"column:ad_type;type:integer;not null"`
	Rounds     []admin.AuctionRoundConfiguration `gorm:"column:rounds;type:jsonb;default:'[]';serializer:json"`
	Pricefloor float64                           `gorm:"column:pricefloor;type:double precision;not null"`
}

//lint:ignore U1000 this method is used by generic struct
func (a *auctionConfiguration) initFromResourceAttrs(attrs *admin.AuctionConfigurationAttrs) {
	name := sql.NullString{}
	if attrs.Name != "" {
		name.String = attrs.Name
		name.Valid = true
	}

	a.Name = name
	a.AppID = attrs.AppID
	a.AdType = dbAdType(attrs.AdType)
	a.Rounds = attrs.Rounds
	a.Pricefloor = attrs.Pricefloor
}

//lint:ignore U1000 this method is used by generic struct
func (a *auctionConfiguration) toResource() admin.AuctionConfiguration {
	return admin.AuctionConfiguration{
		ID: a.ID,
		AuctionConfigurationAttrs: admin.AuctionConfigurationAttrs{
			Name:       a.Name.String,
			AppID:      a.AppID,
			AdType:     adType(a.AdType),
			Rounds:     a.Rounds,
			Pricefloor: a.Pricefloor,
		},
	}
}
