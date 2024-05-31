package apihandlers_test

import (
	"context"
	"errors"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v1/apihandlers"
	"net/http"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionmocks "github.com/bidon-io/bidon-backend/internal/auction/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	sdkapimocks "github.com/bidon-io/bidon-backend/internal/sdkapi/v1/apihandlers/mocks"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
)

func testHelperAuctionHandler(t *testing.T) *apihandlers.AuctionHandler {
	app := sdkapi.App{ID: 1}
	geodata := geocoder.GeoData{CountryCode: "US"}
	auctionConfig := &auction.Config{
		ID:  1,
		UID: "1701972528521547776",
		Rounds: []auction.RoundConfig{
			{
				ID:      "ROUND_1",
				Demands: []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey},
				Timeout: 15000,
			},
			{
				ID:      "ROUND_2",
				Demands: []adapter.Key{adapter.UnityAdsKey},
				Bidding: []adapter.Key{adapter.BidmachineKey},
				Timeout: 15000,
			},
			{
				ID:      "ROUND_3",
				Demands: []adapter.Key{adapter.ApplovinKey},
				Timeout: 15000,
			},
			{
				ID:      "ROUND_4",
				Demands: []adapter.Key{adapter.UnityAdsKey, adapter.ApplovinKey},
				Timeout: 15000,
			},
			{
				ID:      "ROUND_5",
				Bidding: []adapter.Key{adapter.BidmachineKey},
				Timeout: 15000,
			},
		},
	}
	lineItems := []auction.LineItem{
		{ID: "test", UID: "1701972528521547776", PriceFloor: 0.1, PlacementID: "", AdUnitID: "test_id"},
	}
	pf := 0.1
	adUnits := []auction.AdUnit{
		{
			DemandID:   "test",
			UID:        "1701972528521547776",
			Label:      "test",
			PriceFloor: &pf,
			BidType:    schema.RTBBidType,
			Extra: map[string]any{
				"placement_id": "test_id",
			},
		},
	}
	segments := []segment.Segment{
		{
			ID:      1,
			UID:     "1701972528521547776",
			Filters: []segment.Filter{{Type: "country", Name: "country", Operator: "IN", Values: []string{"US", "UK"}}},
		},
	}

	appFetcher := &sdkapimocks.AppFetcherMock{
		FetchCachedFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return app, nil
		},
	}
	gcoder := &sdkapimocks.GeocoderMock{
		LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
			return geodata, nil
		},
	}
	configFetcher := &sdkapimocks.ConfigFetcherMock{
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error) {
			return auctionConfig, nil
		},
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, key string, aucUID string) *auction.Config {
			return nil
		},
	}
	lineItemsMatcher := &auctionmocks.LineItemsMatcherMock{
		MatchCachedFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.LineItem, error) {
			return lineItems, nil
		},
	}
	adUnitsMatcher := &auctionmocks.AdUnitsMatcherMock{
		MatchCachedFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
			return adUnits, nil
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return segments, nil
		},
	}
	auctionBuilder := &auction.Builder{
		ConfigFetcher:    configFetcher,
		LineItemsMatcher: lineItemsMatcher,
	}
	auctionBuilderV2 := &auction.BuilderV2{
		ConfigFetcher:  configFetcher,
		AdUnitsMatcher: adUnitsMatcher,
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}

	// Create a new AuctionHandler instance

	handler := &apihandlers.AuctionHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.AuctionRequest, *schema.AuctionRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      gcoder,
		},
		AuctionBuilder:   auctionBuilder,
		AuctionBuilderV2: auctionBuilderV2,
		SegmentMatcher:   segmentMatcher,
		EventLogger:      &event.Logger{Engine: &engine.Log{}},
	}

	return handler
}

func TestAuctionHandler_Handle(t *testing.T) {
	tests := []struct {
		name                 string
		sdkVersion           string
		requestPath          string
		expectedResponsePath string
		expectedStatusCode   int
		wantErr              bool
		err                  error
	}{
		{
			name:                 "OK version 0,4",
			sdkVersion:           "0.4.0",
			requestPath:          "testdata/auction/valid_request.json",
			expectedResponsePath: "testdata/auction/valid_response.json",
			expectedStatusCode:   http.StatusOK,
			wantErr:              false,
		},
		{
			name:               "Err NoAdsFound version 0,4",
			sdkVersion:         "0.4.0",
			requestPath:        "testdata/auction/noads_request.json",
			expectedStatusCode: http.StatusUnprocessableEntity,
			wantErr:            true,
			err:                sdkapi.ErrNoAdsFound,
		},
		{
			name:                 "OK version 0,5",
			sdkVersion:           "0.5",
			requestPath:          "testdata/auction/valid_request.json",
			expectedResponsePath: "testdata/auction/valid_response_v2.json",
			expectedStatusCode:   http.StatusOK,
			wantErr:              false,
		},
		{
			name:               "Err NoAdsFound version 0,5",
			sdkVersion:         "0.5",
			requestPath:        "testdata/auction/noads_request.json",
			expectedStatusCode: http.StatusUnprocessableEntity,
			wantErr:            true,
			err:                sdkapi.ErrNoAdsFound,
		},
		{
			name:               "Err Invalid SDKVesrion",
			sdkVersion:         "",
			requestPath:        "testdata/auction/valid_request.json",
			expectedStatusCode: http.StatusUnprocessableEntity,
			wantErr:            true,
			err:                sdkapi.ErrInvalidSDKVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := os.ReadFile(tt.requestPath)
			if err != nil {
				t.Fatalf("Error reading request file: %v", err)
			}

			handler := testHelperAuctionHandler(t)
			rec, err := ExecuteRequest(t, handler, http.MethodPost, "/auction/interstitial", string(reqBody), &RequestOptions{
				Headers: map[string]string{
					"X-Bidon-Version": tt.sdkVersion,
				},
			})
			CheckResponseCode(t, err, rec.Code, tt.expectedStatusCode)

			if tt.wantErr {
				if !errors.Is(err, tt.err) {
					t.Errorf("Expected error %v, got: %v", tt.err, err)
				}
			} else {
				expectedResponseJson, err := os.ReadFile(tt.expectedResponsePath)
				if err != nil {
					t.Fatalf("Error reading response file: %v", err)
				}
				CheckResponses(t, expectedResponseJson, rec.Body.Bytes())
			}
		})
	}
}
