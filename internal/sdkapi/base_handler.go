package sdkapi

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

// BaseHandler provides common functionality between sdkapi handlers
type BaseHandler struct {
	AppFetcher AppFetcher
}

type AppFetcher interface {
	Fetch(ctx context.Context, appKey, appBundle string) (*App, error)
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

	return &request{
		raw: raw,
		app: app,
	}, nil
}
