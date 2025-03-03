package apihandlers_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers/mocks"
)

func SetupStatsV2Handler() apihandlers.StatsHandler {
	mockHandler := &mocks.StatsNotificationHandlerMock{}
	mockHandler.HandleStatsFunc = func(contextMoqParam context.Context, stats schema.StatsV2, config *auction.Config, _ string, _ string) {
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
	configFetcher := mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, key string, aucUID string) *auction.Config {
			return auctionConfig
		},
	}

	return apihandlers.StatsHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.StatsV2Request, *schema.StatsV2Request]{
			AppFetcher:    AppFetcherMock(),
			ConfigFetcher: &configFetcher,
			Geocoder:      GeocoderMock(),
		},
		EventLogger:         &event.Logger{Engine: &engine.Log{}},
		NotificationHandler: mockHandler,
	}
}

func TestStatsV2Handler_Handle(t *testing.T) {
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
			handler := SetupStatsV2Handler()
			rec, err := ExecuteRequest(t, &handler, http.MethodPost, "/v2/stats/banner", string(reqBody), nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error %v, got: %v", tt.wantErr, err)
			}

			CheckResponseCode(t, err, rec.Code, tt.expectedCode)
		})
	}
}
