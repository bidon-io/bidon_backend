package sdkapi

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

// BaseHandler provides common functionality between sdkapi handlers
type BaseHandler struct {
	AppFetcher AppFetcher
	Geocoder   Geocoder
}

type AppFetcher interface {
	Fetch(ctx context.Context, appKey, appBundle string) (*App, error)
}

type Geocoder interface {
	FindGeoData(ctx context.Context, ipString string) (*geocoder.GeoData, error)
}

func (b *BaseHandler) resolveRequest(c echo.Context) (*request, error) {
	var raw schema.Request
	if err := c.Bind(&raw); err != nil {
		return nil, err
	}

	app, err := b.AppFetcher.Fetch(c.Request().Context(), raw.App.Key, raw.App.Bundle)
	if err != nil {
		return nil, err
	}

	geoData, err := b.Geocoder.FindGeoData(c.Request().Context(), c.RealIP())
	if err != nil {
		c.Logger().Infof("Failed to lookup ip: %v", err)
	}

	return &request{
		raw:     raw,
		app:     app,
		geoData: geoData,
	}, nil
}
