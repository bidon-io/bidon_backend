package apihandlers_test

import (
	"context"
	"encoding/json"
	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v1/apihandlers/mocks"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"strings"
	"testing"
)

type Handler interface {
	Handle(c echo.Context) error
}

type RequestOptions struct {
	Headers map[string]string
	Params  map[string]string
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
		FetchCachedFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return app, nil
		},
	}
}

func ExecuteRequest(t *testing.T, handler Handler, method, path, body string, options *RequestOptions) (*httptest.ResponseRecorder, error) {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := config.Echo()
	c := e.NewContext(req, rec)

	if options != nil {
		// Set headers in Echo context If any.
		for k, v := range options.Headers {
			c.Request().Header.Set(k, v)
		}
		// Set path parameters in Echo context If any.
		for k, v := range options.Params {
			c.SetParamNames(k)
			c.SetParamValues(v)
		}
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

// CheckResponses is a helper function to assert the response body.
func CheckResponses(t *testing.T, expectedResponseJson, actualResponseJson []byte) {
	t.Helper()

	var actualResponse interface{}
	var expectedResponse interface{}
	err := json.Unmarshal(actualResponseJson, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to parse JSON1: %s", err)
	}
	err = json.Unmarshal(expectedResponseJson, &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to parse JSON2: %s", err)
	}

	if diff := cmp.Diff(actualResponse, expectedResponse); diff != "" {
		t.Errorf("Response mismatch (-want, +got):\n%s", diff)
	}
}
