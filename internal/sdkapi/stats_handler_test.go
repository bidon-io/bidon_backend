package sdkapi_test

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/mocks"
	"net/http"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func SetupStatsHandler() sdkapi.StatsHandler {
	mockHandler := &mocks.StatsNotificationHandlerMock{}
	mockHandler.HandleStatsFunc = func(contextMoqParam context.Context, stats schema.Stats, config auction.Config) error {
		return nil
	}
	auctionConfig := &auction.Config{
		ID:  1,
		UID: "123",
		Rounds: []auction.RoundConfig{
			{
				ID:      "ROUND_BANNER_1",
				Demands: []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey},
				Timeout: 15000,
			},
		},
	}
	configMatcher := &mocks.ConfigMatcherMock{
		MatchByIdFunc: func(ctx context.Context, appID, id int64) *auction.Config {
			return auctionConfig
		},
	}

	return sdkapi.StatsHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.StatsRequest, *schema.StatsRequest]{
			AppFetcher: AppFetcherMock(),
			Geocoder:   GeocoderMock(),
		},
		ConfigMatcher:       configMatcher,
		EventLogger:         &event.Logger{Engine: &engine.Log{}},
		NotificationHandler: mockHandler,
	}
}

func TestStatsHandler_Handle(t *testing.T) {
	tests := []struct {
		name         string
		requestPath  string
		expectedCode int
		wantErr      bool
	}{
		{
			name:         "valid request",
			requestPath:  "testdata/stats/valid_request.json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid request",
			requestPath:  "testdata/stats/invalid_request.json",
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
			handler := SetupStatsHandler()
			rec, err := ExecuteRequest(t, &handler, http.MethodPost, "/stats/banner", string(reqBody), nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error %v, got: %v", tt.wantErr, err)
			}

			CheckResponseCode(t, err, rec.Code, tt.expectedCode)
		})
	}
}