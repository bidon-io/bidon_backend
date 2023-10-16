package sdkapi_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func SetupClickHandler() sdkapi.ClickHandler {
	return sdkapi.ClickHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.ClickRequest, *schema.ClickRequest]{
			AppFetcher: AppFetcherMock(),
			Geocoder:   GeocoderMock(),
		},
		EventLogger: &event.Logger{Engine: &engine.Log{}},
	}
}

func TestClickHandler_Handle(t *testing.T) {
	tests := []struct {
		name         string
		requestPath  string
		expectedCode int
		wantErr      bool
	}{
		{
			name:         "valid request",
			requestPath:  "testdata/click/valid_request.json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid request",
			requestPath:  "testdata/click/invalid_request.json",
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
			handler := SetupClickHandler()
			rec, err := ExecuteRequest(t, &handler, http.MethodPost, "/click", string(reqBody), nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error %v, got: %v", tt.wantErr, err)
			}

			CheckResponseCode(t, err, rec.Code, tt.expectedCode)
		})
	}
}
