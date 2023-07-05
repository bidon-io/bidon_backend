package store

import (
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type CountryRepo = resourceRepo[admin.Country, admin.CountryAttrs, db.Country]

func NewCountryRepo(db *db.DB) *CountryRepo {
	return &CountryRepo{
		db:           db,
		mapper:       countryMapper{},
		associations: []string{},
	}
}

type countryMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m countryMapper) dbModel(a *admin.CountryAttrs, id int64) *db.Country {
	humanName := sql.NullString{}
	if a.HumanName != "" {
		humanName.String = a.HumanName
		humanName.Valid = true
	}

	return &db.Country{
		Model:      db.Model{ID: id},
		Alpha2Code: a.Alpha2Code,
		Alpha3Code: a.Alpha3Code,
		HumanName:  humanName,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m countryMapper) resource(c *db.Country) admin.Country {
	return admin.Country{
		ID: c.ID,
		CountryAttrs: admin.CountryAttrs{
			HumanName:  c.HumanName.String,
			Alpha2Code: c.Alpha2Code,
			Alpha3Code: c.Alpha3Code,
		},
	}
}
