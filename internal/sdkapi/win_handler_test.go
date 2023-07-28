package sdkapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

func SetupWinHandler() sdkapi.WinHandler {
	app := sdkapi.App{ID: 1}
	geodata := geocoder.GeoData{CountryCode: "US"}

	// Create a mock WinNotificationHandler
	mockHandler := &mocks.WinNotificationHandlerMock{}
	mockHandler.HandleWinFunc = func(ctx context.Context, imp *schema.Imp, responses []*adapters.DemandResponse) error {
		return nil
	}
	appFetcher := &mocks.AppFetcherMock{
		FetchFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return app, nil
		},
	}
	geocoder := &mocks.GeocoderMock{
		LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
			return geodata, nil
		},
	}

	// Create a new WinHandler instance
	return sdkapi.WinHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.WinRequest, *schema.WinRequest]{
			AppFetcher: appFetcher,
			Geocoder:   geocoder,
		},
		EventLogger:         &event.Logger{Engine: &engine.Log{}},
		NotificationHandler: mockHandler,
	}
}

func TestWinHandler_Handle(t *testing.T) {
	rec := httptest.NewRecorder()

	reqBody, err := os.ReadFile("testdata/win/valid_request.json")
	if err != nil {
		t.Fatalf("Error reading request file: %v", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/win", strings.NewReader(string(reqBody[:])))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	handler := SetupWinHandler()

	// Call the Handle method with the mock request and response
	e := config.Echo("sdkapi-test", nil)
	c := e.NewContext(req, rec)

	err = handler.Handle(c)
	// Verify that the Handle method returns no error and the response status code is 200
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestWinHandler_Handle_InvalidRequest(t *testing.T) {
	// Create a mock request and response with an invalid request body
	rec := httptest.NewRecorder()

	reqBody := `{"imp": {"id": "imp-1"}, "responses": [{"bid": {"id": "bid-1", "impid": "imp-1", "price": "invalid"}}]}`

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/win", strings.NewReader(string(reqBody[:])))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	handler := SetupWinHandler()

	// Call the Handle method with the mock request and response
	e := config.Echo("sdkapi-test", nil)
	c := e.NewContext(req, rec)
	err := handler.Handle(c)

	// Verify that the Handle method returns an error
	if err == nil {
		t.Fatal("Expected error, got no error")
	}
}
