package sdkapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var ErrAppNotValid = echo.NewHTTPError(http.StatusUnprocessableEntity, "App is not valid")
var ErrNoAdsFound = echo.NewHTTPError(http.StatusUnprocessableEntity, "No ads found")
var ErrInvalidSDKVersion = echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid SDK version")
var ErrNoAdaptersFound = echo.NewHTTPError(http.StatusUnprocessableEntity, "No adapters found")
