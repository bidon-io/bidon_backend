package store

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
)

type CountryRepo = resourceRepo[admin.Country, admin.CountryAttrs, country, *country]

type country struct {
	Model
	Alpha2Code string         `gorm:"column:alpha_2_code;type:varchar;not null"`
	Alpha3Code string         `gorm:"column:alpha_3_code;type:varchar;not null"`
	HumanName  sql.NullString `gorm:"column:human_name;type:varchar"`
}

//lint:ignore U1000 this method is used by generic struct
func (c *country) initFromResourceAttrs(attrs *admin.CountryAttrs) {
	c.Alpha2Code = attrs.Alpha2Code
	c.Alpha3Code = attrs.Alpha3Code
	if attrs.HumanName != "" {
		c.HumanName.String = attrs.HumanName
		c.HumanName.Valid = true
	}
}

//lint:ignore U1000 this method is used by generic struct
func (c *country) toResource() admin.Country {
	return admin.Country{
		ID: c.ID,
		CountryAttrs: admin.CountryAttrs{
			HumanName:  c.HumanName.String,
			Alpha2Code: c.Alpha2Code,
			Alpha3Code: c.Alpha3Code,
		},
	}
}
