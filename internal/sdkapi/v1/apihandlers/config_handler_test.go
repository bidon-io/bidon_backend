package apihandlers_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/v1/apihandlers"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v1/apihandlers/mocks"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
)

func SetupConfigHandler() apihandlers.ConfigHandler {
	app := sdkapi.App{ID: 1}
	sgmnt := segment.Segment{
		ID:  1,
		UID: "1",
		Filters: []segment.Filter{
			{Type: "country", Operator: "IN", Values: []string{"US"}},
		},
	}

	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	adapterInitConfigsFetcher := &mocks.AdapterInitConfigsFetcherMock{
		FetchAdapterInitConfigsFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key, setAmazonSlots bool, setOrder bool) ([]sdkapi.AdapterInitConfig, error) {
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
				&sdkapi.BigoAdsInitConfig{
					AppID: fmt.Sprintf("bigo_app_%d", app.ID),
				},
				&sdkapi.DTExchangeInitConfig{
					AppID: fmt.Sprintf("dtexchange_app_%d", app.ID),
				},
				&sdkapi.GAMInitConfig{
					AppID:       fmt.Sprintf("dtexchange_app_%d", app.ID),
					NetworkCode: "network_code",
				},
				&sdkapi.MetaInitConfig{
					AppID:     fmt.Sprintf("meta_app_%d", app.ID),
					AppSecret: fmt.Sprintf("meta_app_%d_secret", app.ID),
				},
				&sdkapi.MintegralInitConfig{
					AppID:  fmt.Sprintf("mintegral_app_%d", app.ID),
					AppKey: "mintegral",
				},
				&sdkapi.UnityAdsInitConfig{
					GameID: fmt.Sprintf("unity_game_%d", app.ID),
				},
				&sdkapi.VungleInitConfig{
					AppID: fmt.Sprintf("vungle_app_%d", app.ID),
				},
				&sdkapi.MobileFuseInitConfig{
					PublisherID: fmt.Sprintf("mobilefuse_publisher_%d", app.ID),
					AppKey:      fmt.Sprintf("mobilefuse_app_%d", app.ID),
				},
				&sdkapi.InmobiInitConfig{
					AccountID: fmt.Sprintf("inmobi_account_%d", app.ID),
					AppKey:    fmt.Sprintf("inmobi_app_%d", app.ID),
				},
				&sdkapi.AmazonInitConfig{
					AppKey: fmt.Sprintf("amazon_app_%d", app.ID),
				},
			}, nil
		},
	}

	return apihandlers.ConfigHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.ConfigRequest, *schema.ConfigRequest]{
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
		sdkVersion   string
		requestPath  string
		expectedCode int
		wantErr      bool
	}{
		{
			name:         "valid request",
			sdkVersion:   "0.4.0",
			requestPath:  "testdata/config/valid_request.json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid request",
			sdkVersion:   "0.4.0",
			requestPath:  "testdata/config/invalid_request.json",
			expectedCode: http.StatusUnprocessableEntity,
			wantErr:      true,
		},
		{
			name:         "valid request",
			sdkVersion:   "0.5.0",
			requestPath:  "testdata/config/valid_request.json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid request android",
			sdkVersion:   "0.5.0",
			requestPath:  "testdata/config/valid_request_android.json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid sdk version",
			sdkVersion:   "",
			requestPath:  "testdata/config/valid_request.json",
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
			rec, err := ExecuteRequest(t, &handler, http.MethodPost, "/config", string(reqBody), &RequestOptions{
				Headers: map[string]string{
					"X-Bidon-Version": tt.sdkVersion,
				},
			})

			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error %v, got: %v", tt.wantErr, err)
			}

			CheckResponseCode(t, err, rec.Code, tt.expectedCode)
		})
	}
}
