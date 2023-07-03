package sdkapi

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"golang.org/x/exp/maps"
)

// request wraps raw schema.Request and includes additional data that is needed for all sdkapi handlers
type request struct {
	raw     schema.Request
	app     *App
	geoData *geocoder.GeoData
}

// App represents an app for the purposes of the SDK API
type App struct {
	ID int64
}

func (r *request) adFormat() ad.Format {
	obj := r.raw.AdObject
	if obj.Banner != nil {
		return obj.Banner.Format
	}

	return ad.EmptyFormat
}

func (r *request) adapterKeys() []adapter.Key {
	return maps.Keys(r.raw.Adapters)
}

func (r *request) countryCode() string {
	if r.geoData == nil {
		var countryCode string

		if r.raw.Device.Geo != nil {
			countryCode = r.raw.Device.Geo.Country
		} else if r.raw.Geo != nil {
			countryCode = r.raw.Geo.Country
		} else {
			countryCode = geocoder.UNKNOWN_COUNTRY_CODE
		}

		return countryCode
	}

	return r.geoData.CountryCode
}
