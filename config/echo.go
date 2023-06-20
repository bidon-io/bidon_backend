package config

import (
	"errors"
	"net/http"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Echo() *echo.Echo {
	e := echo.New()
	e.Debug = Env != ProdEnv
	e.HTTPErrorHandler = HTTPErrorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))

	return e
}

// HTTPErrorHandler is the default error handler for Bidon services.
// Errors that do not wrap echo.HTTPError are considered unexpected and are sent to Sentry.
func HTTPErrorHandler(err error, c echo.Context) {
	var herr *echo.HTTPError
	if !errors.As(err, &herr) {
		hub := sentryecho.GetHubFromContext(c)
		if hub != nil {
			hub.CaptureException(err)
		}

		var message string
		if c.Echo().Debug {
			message = err.Error()
		} else {
			message = http.StatusText(http.StatusInternalServerError)
		}

		herr = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: message,
		}
	}

	response := map[string]any{
		"error": map[string]any{
			"code":    herr.Code,
			"message": herr.Message,
		},
	}

	if err := c.JSON(herr.Code, response); err != nil {
		c.Logger().Error(err)
	}
}
