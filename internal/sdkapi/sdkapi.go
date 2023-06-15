package sdkapi

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckBidonHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("X-Bidon-Version") == "" {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Request should contain X-Bidon-Version header")
		}

		return next(c)
	}
}

func ErrorHandler(err error, c echo.Context) {
	var herr *echo.HTTPError
	if !errors.As(err, &herr) {
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
