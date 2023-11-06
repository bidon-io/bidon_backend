package sdkapi

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

// BaseHandler provides common functionality between sdkapi handlers
type BaseHandler[T any, PT rawRequest[T]] struct {
	AppFetcher    AppFetcher
	ConfigFetcher ConfigFetcher
	Geocoder      Geocoder
}

// App represents an app for the purposes of the SDK API
type App struct {
	ID int64
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . AppFetcher ConfigFetcher Geocoder

type AppFetcher interface {
	Fetch(ctx context.Context, appKey, appBundle string) (App, error)
}

type ConfigFetcher interface {
	FetchByUIDCached(ctx context.Context, appId int64, id, uid string) *auction.Config
	Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error)
}

type Geocoder interface {
	Lookup(ctx context.Context, ipString string) (geocoder.GeoData, error)
}

func (b *BaseHandler[T, PT]) resolveRequest(c echo.Context) (*request[T, PT], error) {
	var raw T

	if err := c.Bind(&raw); err != nil {
		return nil, err
	}

	req := PT(&raw)
	req.NormalizeValues()
	sdkVersion := c.Request().Header.Get("X-Bidon-Version")
	req.SetSDKVersion(sdkVersion)

	if err := c.Validate(&raw); err != nil {
		return nil, err
	}

	rawApp := req.GetApp()

	app, err := b.AppFetcher.Fetch(c.Request().Context(), rawApp.Key, rawApp.Bundle)
	if err != nil {
		return nil, err
	}

	var auctionConfig *auction.Config
	if b.ConfigFetcher != nil {
		id, uid := req.GetAuctionConfigurationParams()
		auctionConfig = b.ConfigFetcher.FetchByUIDCached(c.Request().Context(), app.ID, id, uid)
	}
	if auctionConfig != nil {
		req.SetAuctionConfigurationParams(auctionConfig.ID, auctionConfig.UID)
	}

	geoData, err := b.Geocoder.Lookup(c.Request().Context(), c.RealIP())
	if err != nil {
		c.Logger().Infof("Failed to lookup ip: %v", err)
	}

	return &request[T, PT]{
		raw:           raw,
		app:           app,
		auctionConfig: auctionConfig,
		geoData:       geoData,
	}, nil
}

type rawRequest[T any] interface {
	*T
	GetApp() schema.App
	GetGeo() schema.Geo
	SetSDKVersion(string)
	NormalizeValues()
	GetAuctionConfigurationParams() (string, string)
	SetAuctionConfigurationParams(int64, string)
}

// request wraps raw request and includes additional data that is needed for all sdkapi handlers
type request[T any, PT rawRequest[T]] struct {
	raw           T
	app           App
	auctionConfig *auction.Config
	geoData       geocoder.GeoData
}

func (r *request[T, PT]) countryCode() string {
	if r.geoData.CountryCode != "" {
		return r.geoData.CountryCode
	}

	geo := PT(&r.raw).GetGeo()
	if geo.Country != "" {
		return geo.Country
	}

	return geocoder.UnknownCountryCode
}
