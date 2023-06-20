package sdkapi

import (
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
