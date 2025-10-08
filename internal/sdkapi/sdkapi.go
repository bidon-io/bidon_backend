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
	Badv       string
	Bcat       string
	Bapp       string
}

func (a *App) GetBadv() string {
	if a == nil {
		return ""
	}
	return a.Badv
}

func (a *App) GetBcat() string {
	if a == nil {
		return ""
	}
	return a.Bcat
}

func (a *App) GetBapp() string {
	if a == nil {
		return ""
	}
	return a.Bapp
}

func GetBlockedAdvertisersList(app *App) []string {
	return parseCommaSeparatedList(app.GetBadv())
}

func GetBlockedCategoriesList(app *App) []string {
	return parseCommaSeparatedList(app.GetBcat())
}

func GetBlockedAppsList(app *App) []string {
	return parseCommaSeparatedList(app.GetBapp())
}

func parseCommaSeparatedList(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for part := range strings.SplitSeq(s, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
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
