package sdkapi_test

import (
	"context"
	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/mocks"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"strings"
	"testing"
)

type Handler interface {
	Handle(c echo.Context) error
}

func GeocoderMock() *mocks.GeocoderMock {
	geodata := geocoder.GeoData{CountryCode: "US"}
	return &mocks.GeocoderMock{
		LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
			return geodata, nil
		},
	}
}

func AppFetcherMock() *mocks.AppFetcherMock {
	app := sdkapi.App{ID: 1}
	return &mocks.AppFetcherMock{
		FetchFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return app, nil
		},
	}
}

func ExecuteRequest(t *testing.T, handler Handler, method, path, body string, params map[string]string) (*httptest.ResponseRecorder, error) {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := config.Echo()
	c := e.NewContext(req, rec)

	// Set path parameters in Echo context If any.
	for k, v := range params {
		c.SetParamNames(k)
		c.SetParamValues(v)
	}

	err := handler.Handle(c)

	return rec, err
}

// checkResponseCode is a helper function to assert the status code of the response.
func CheckResponseCode(t *testing.T, err error, actualCode int, expectedCode int) {
	t.Helper()

	// Check if the error is of type *echo.HTTPError and compare the code.
	if echoError, ok := err.(*echo.HTTPError); ok {
		if echoError.Code != expectedCode {
			t.Fatalf("Expected status code %d, got %d", expectedCode, echoError.Code)
		}
	} else if err == nil {
		// If there's no error, compare the recorded code with the expected code.
		if actualCode != expectedCode {
			t.Fatalf("Expected status code %d, got %d", expectedCode, actualCode)
		}
	} else {
		t.Fatalf("An unexpected error occurred: %v", err)
	}
}
