package geocoder

import (
	"context"
	"fmt"
	"net"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/oschwald/maxminddb-golang"
)

// Geocoder represents an geocoder.
type Geocoder struct {
	MaxMindDB *maxminddb.Reader
	DB        *db.DB
	Cache     cache
}

type cache interface {
	Get(context.Context, []byte, func(ctx context.Context) (*db.Country, error)) (*db.Country, error)
}

// GeoData represents the geolocation data.
type GeoData struct {
	CountryCode    string
	CountryID      int64
	CountryCode3   string
	CityName       string
	RegionName     string
	RegionCode     string
	Lat            float64
	Lon            float64
	Accuracy       int
	ZipCode        string
	IPService      int
	UnknownCountry bool
	IPString       string
}

type MmdbGeoData struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Subdivisions []struct {
		Names   map[string]string `maxminddb:"names"`
		ISOCode string            `maxminddb:"iso_code"`
	} `maxminddb:"subdivisions"`
	Location struct {
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		AccuracyRadius int     `maxminddb:"accuracy_radius"`
	} `maxminddb:"location"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	Continent struct {
		Code  string            `maxminddb:"code"`
		Names map[string]string `maxminddb:"names"`
	}
}

func (g *MmdbGeoData) ContinentName() string {
	return g.Continent.Names["en"]
}

func (g *MmdbGeoData) CountryCode() string {
	return g.Country.ISOCode
}

func (g *MmdbGeoData) CityName() string {
	return g.City.Names["en"]
}

func (g *MmdbGeoData) SubdivisionName() string {
	if len(g.Subdivisions) == 0 {
		return ""
	}

	return g.Subdivisions[0].Names["en"]
}

func (g *MmdbGeoData) SubdivisionCode() string {
	if len(g.Subdivisions) == 0 {
		return ""
	}

	return g.Subdivisions[0].ISOCode
}

const (
	MaxMindProviderCode = 3
	UnknownCountryCode  = "ZZ"
)

var DefaultCountryCodesForContinents = map[string]string{
	"Europe": "FR",
	"Asia":   "ID",
}

// Lookup finds the geolocation data for the given IP address.
func (g *Geocoder) Lookup(ctx context.Context, ipString string) (GeoData, error) {
	var geoData GeoData

	if g.MaxMindDB == nil {
		return geoData, fmt.Errorf("maxminddb not set")
	}

	var mmdbGeoData MmdbGeoData
	ip := net.ParseIP(ipString)

	err := g.lookupIP(ip, &mmdbGeoData)
	if err != nil {
		return geoData, err
	}

	countryCode := g.countryCodeFor(mmdbGeoData)
	country, err := g.findCountryCached(ctx, countryCode)
	if err != nil {
		return geoData, err
	}

	geoData.CountryCode = countryCode
	geoData.CountryCode3 = country.Alpha3Code
	geoData.UnknownCountry = countryCode == UnknownCountryCode
	geoData.CountryID = country.ID
	geoData.CityName = mmdbGeoData.CityName()
	geoData.RegionName = mmdbGeoData.SubdivisionName()
	geoData.RegionCode = mmdbGeoData.SubdivisionCode()
	geoData.Lat = mmdbGeoData.Location.Latitude
	geoData.Lon = mmdbGeoData.Location.Longitude
	geoData.Accuracy = mmdbGeoData.Location.AccuracyRadius * 1000 // convert kilometers to meters
	geoData.ZipCode = mmdbGeoData.Postal.Code
	geoData.IPService = MaxMindProviderCode
	geoData.IPString = ipString

	return geoData, nil
}

func (g *Geocoder) lookupIP(ip net.IP, mmdbGeoData *MmdbGeoData) error {
	err := g.MaxMindDB.Lookup(ip, mmdbGeoData)
	if err != nil {
		return err
	}
	return nil
}

func (g *Geocoder) countryCodeFor(mmdbGeoData MmdbGeoData) string {
	if mmdbGeoData.Country.ISOCode != "" {
		return mmdbGeoData.Country.ISOCode
	}

	if code, ok := DefaultCountryCodesForContinents[mmdbGeoData.ContinentName()]; ok {
		return code
	}

	return UnknownCountryCode
}

func (g *Geocoder) findCountryCached(ctx context.Context, countryCode string) (*db.Country, error) {
	return g.Cache.Get(ctx, []byte(countryCode), func(ctx context.Context) (*db.Country, error) {
		return g.findCountry(ctx, countryCode)
	})
}

func (g *Geocoder) findCountry(ctx context.Context, countryCode string) (*db.Country, error) {
	var dbCountry db.Country

	if err := g.DB.WithContext(ctx).Where("alpha_2_code = ?", countryCode).First(&dbCountry).Error; err != nil {
		return nil, err
	}

	return &dbCountry, nil
}
