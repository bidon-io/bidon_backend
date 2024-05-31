package apihandlers_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers/mocks"
)

func SetupShowHandler() apihandlers.ShowHandler {
	mockHandler := &mocks.ShowNotificationHandlerMock{}
	mockHandler.HandleShowFunc = func(ctx context.Context, imp *schema.Bid, _ string, _ string) {}
	return apihandlers.ShowHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.ShowRequest, *schema.ShowRequest]{
			AppFetcher: AppFetcherMock(),
			Geocoder:   GeocoderMock(),
		},
		EventLogger:         &event.Logger{Engine: &engine.Log{}},
		NotificationHandler: mockHandler,
	}
}

func TestShowHandler_Handle(t *testing.T) {
	tests := []struct {
		name         string
		requestPath  string
		expectedCode int
		wantErr      bool
	}{
		{
			name:         "valid request",
			requestPath:  "testdata/show/valid_request.json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid request",
			requestPath:  "testdata/show/invalid_request.json",
			expectedCode: http.StatusUnprocessableEntity,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := os.ReadFile(tt.requestPath)
			if err != nil {
				t.Fatalf("Error reading request file: %v", err)
			}
			handler := SetupShowHandler()
			rec, err := ExecuteRequest(t, &handler, http.MethodPost, "/v2/show", string(reqBody), nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error %v, got: %v", tt.wantErr, err)
			}

			CheckResponseCode(t, err, rec.Code, tt.expectedCode)
		})
	}
}
