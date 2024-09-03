package apihandlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers/mocks"
	"github.com/labstack/echo/v4"
)

func SetupLossHandler() apihandlers.LossHandler {
	app := sdkapi.App{ID: 1}
	geodata := geocoder.GeoData{CountryCode: "US"}

	// Create a mock LossNotificationHandler
	mockHandler := &mocks.LossNotificationHandlerMock{}
	mockHandler.HandleLossFunc = func(ctx context.Context, imp *schema.Bid) error {
		return nil
	}
	appFetcher := &mocks.AppFetcherMock{
		FetchCachedFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return app, nil
		},
	}
	geocoder := &mocks.GeocoderMock{
		LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
			return geodata, nil
		},
	}

	// Create a new LossHandler instance
	return apihandlers.LossHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.LossRequest, *schema.LossRequest]{
			AppFetcher: appFetcher,
			Geocoder:   geocoder,
		},
		EventLogger:         &event.Logger{Engine: &engine.Log{}},
		NotificationHandler: mockHandler,
	}
}

func TestLossHandler_Handle(t *testing.T) {
	rec := httptest.NewRecorder()

	reqBody, err := os.ReadFile("testdata/loss/valid_request.json")
	if err != nil {
		t.Fatalf("Error reading request file: %v", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/loss", strings.NewReader(string(reqBody[:])))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	handler := SetupLossHandler()

	// Call the Handle method with the mock request and response
	e := config.Echo()
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

func TestLossHandler_Handle_InvalidRequest(t *testing.T) {
	// Create a mock request and response with an invalid request body
	rec := httptest.NewRecorder()

	reqBody := `{"imp": {"id": "imp-1"}, "responses": [{"bid": {"id": "bid-1", "impid": "imp-1", "price": "invalid"}}]}`

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/loss", strings.NewReader(string(reqBody[:])))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	handler := SetupLossHandler()

	// Call the Handle method with the mock request and response
	e := config.Echo()
	c := e.NewContext(req, rec)
	err := handler.Handle(c)

	// Verify that the Handle method returns an error
	if err == nil {
		t.Fatal("Expected error, got no error")
	}
}
