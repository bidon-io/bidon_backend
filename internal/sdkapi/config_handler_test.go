package sdkapi_test

import (
	"context"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
	"net/http"
	"os"
	"testing"
)

func SetupConfigHandler() sdkapi.ConfigHandler {
	app := sdkapi.App{ID: 1}
	sgmnt := segment.Segment{
		ID:      1,
		UID:     "1",
		Filters: []segment.Filter{segment.Filter{Type: "country", Operator: "IN", Values: []string{"US"}}},
	}

	segmentFetcher := &segmentmocks.FetcherMock{
		FetchFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	adapterInitConfigsFetcher := &mocks.AdapterInitConfigsFetcherMock{
		FetchAdapterInitConfigsFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]sdkapi.AdapterInitConfig, error) {
			return []sdkapi.AdapterInitConfig{
				&sdkapi.AdmobInitConfig{
					AppID: fmt.Sprintf("admob_app_%d", app.ID),
				},
				&sdkapi.ApplovinInitConfig{
					AppKey: "applovin",
					SDKKey: "applovin",
				},
				&sdkapi.BidmachineInitConfig{
					SellerID:        "1",
					Endpoint:        "x.appbaqend.com",
					MediationConfig: []string{"one", "two"},
				},
				&sdkapi.DTExchangeInitConfig{
					AppID: fmt.Sprintf("dtexchange_app_%d", app.ID),
				},
				&sdkapi.MetaInitConfig{
					AppID:     fmt.Sprintf("meta_app_%d", app.ID),
					AppSecret: fmt.Sprintf("meta_app_%d_secret", app.ID),
				},
				&sdkapi.MintegralInitConfig{
					AppID:  fmt.Sprintf("mintegral_app_%d", app.ID),
					AppKey: "mintegral",
				},
			}, nil
		},
	}

	return sdkapi.ConfigHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.ConfigRequest, *schema.ConfigRequest]{
			AppFetcher: AppFetcherMock(),
			Geocoder:   GeocoderMock(),
		},
		EventLogger:               &event.Logger{Engine: &engine.Log{}},
		SegmentMatcher:            segmentMatcher,
		AdapterInitConfigsFetcher: adapterInitConfigsFetcher,
	}
}

func TestConfigHandler_Handle(t *testing.T) {
	tests := []struct {
		name         string
		requestPath  string
		expectedCode int
		wantErr      bool
	}{
		{
			name:         "valid request",
			requestPath:  "testdata/config/valid_request.json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid request",
			requestPath:  "testdata/config/invalid_request.json",
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
			handler := SetupConfigHandler()
			rec, err := ExecuteRequest(t, &handler, http.MethodPost, "/config", string(reqBody), nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error %v, got: %v", tt.wantErr, err)
			}

			CheckResponseCode(t, err, rec.Code, tt.expectedCode)
		})
	}
}
