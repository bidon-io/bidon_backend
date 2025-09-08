package sdkapi

import (
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App represents an app for the purposes of the SDK API
type App struct {
	ID         int64
	StoreID    string
	StoreURL   string
	Categories []string
}

func CheckBidonHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		skipper := skipIfUtilRoutes()
		if skipper(c) {
			return next(c)
		}

		if c.Request().Header.Get("X-Bidon-Version") == "" {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Request should contain X-Bidon-Version header")
		}

		return next(c)
	}
}

func LogError(c echo.Context, err error) {
	c.Logger().Error(err)

	hub := sentryecho.GetHubFromContext(c)
	if hub != nil {
		client, scope := hub.Client(), hub.Scope()
		client.CaptureException(
			err,
			&sentry.EventHint{Context: c.Request().Context()},
			scope,
		)
	}
}

func skipIfUtilRoutes() middleware.Skipper {
	return func(c echo.Context) bool {
		return strings.EqualFold(c.Path(), "/openapi.json")
	}
}
