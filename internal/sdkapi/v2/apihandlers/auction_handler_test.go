package apihandlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionmocks "github.com/bidon-io/bidon-backend/internal/auction/mocks"
	"github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers"
	handlersmocks "github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers/mocks"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
	"github.com/labstack/echo/v4"
)

func testHelperAuctionHandler() *apihandlers.AuctionHandler {
	app := sdkapi.App{ID: 1}
	geodata := geocoder.GeoData{CountryCode: "US"}
	segments := []segment.Segment{
		{
			ID:      1,
			UID:     "1701972528521547776",
			Filters: []segment.Filter{{Type: "country", Name: "country", Operator: "IN", Values: []string{"US", "UK"}}},
		},
	}
	auctionConfig := &auction.Config{
		ID:        1,
		UID:       "1701972528521547776",
		Demands:   []adapter.Key{adapter.ApplovinKey, adapter.BidmachineKey},
		AdUnitIDs: []int64{1, 2, 3},
		Timeout:   30000,
	}
	pf := 0.1
	applovinPf := 0.8
	adUnits := []auction.AdUnit{
		{
			DemandID:   "amazon",
			Label:      "amazon",
			PriceFloor: &pf,
			UID:        "123_amazon",
			BidType:    schema.RTBBidType,
			Timeout:    store.AdUnitTimeout,
			Extra: map[string]any{
				"slot_uuid": "uuid1",
			},
		},
		{
			DemandID:   "meta",
			Label:      "meta",
			PriceFloor: &pf,
			UID:        "123_meta",
			BidType:    schema.RTBBidType,
			Timeout:    store.AdUnitTimeout,
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
		{
			DemandID:   "mobilefuse",
			Label:      "mobilefuse",
			PriceFloor: &pf,
			UID:        "123_mobilefuse",
			BidType:    schema.RTBBidType,
			Timeout:    store.AdUnitTimeout,
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
		{
			DemandID: "vungle",
			Label:    "vungle",
			UID:      "123_vungle",
			BidType:  schema.RTBBidType,
			Timeout:  store.AdUnitTimeout,
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
		{
			DemandID:   "applovin",
			Label:      "Applovin",
			PriceFloor: &applovinPf,
			UID:        "123_applovin",
			BidType:    schema.CPMBidType,
			Timeout:    store.AdUnitTimeout,
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
	}

	adUnitsMatcher := &auctionmocks.AdUnitsMatcherMock{
		MatchCachedFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
			return adUnits, nil
		},
	}
	appFetcher := &handlersmocks.AppFetcherMock{
		FetchCachedFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return app, nil
		},
	}
	gcoder := &handlersmocks.GeocoderMock{
		LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
			return geodata, nil
		},
	}
	configFetcher := &handlersmocks.ConfigFetcherMock{
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error) {
			return auctionConfig, nil
		},
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, key string, aucUID string) *auction.Config {
			return nil
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return segments, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	biddingAdaptersConfigBuilder := &auctionmocks.BiddingAdaptersConfigBuilderMock{
		BuildFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key, adUnitsMap *auction.AdUnitsMap) (adapter.ProcessedConfigsMap, error) {
			return adapter.ProcessedConfigsMap{
				adapter.ApplovinKey: map[string]any{
					"app_key": "123",
				},
				adapter.BidmachineKey: map[string]any{},
				adapter.MetaKey: map[string]any{
					"app_id":     "123",
					"app_secret": "123",
					"seller_id":  "123",
					"tag_id":     "123",
				},
				adapter.MobileFuseKey: map[string]any{
					"tag_id": "123",
				},
				adapter.AmazonKey: map[string]any{
					"price_points_map": map[string]any{
						"price_point_1": map[string]any{"price": 0.1, "price_point": "price_point_1"},
					},
				},
			}, nil
		},
	}
	biddingBuilder := &auctionmocks.BiddingBuilderMock{
		HoldAuctionFunc: func(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error) {
			return bidding.AuctionResult{
				RoundNumber: 0,
				Bids: []adapters.DemandResponse{
					{
						DemandID: "amazon",
						SlotUUID: "uuid1",
						Bid: &adapters.BidDemandResponse{
							Price:    0.5,
							ID:       "111",
							ImpID:    "222",
							DemandID: adapter.MetaKey,
						},
					},
					{
						DemandID: "meta",
						Bid: &adapters.BidDemandResponse{
							Payload:  "payload",
							Price:    0.6,
							ID:       "123",
							ImpID:    "456",
							DemandID: adapter.MetaKey,
						},
					},
					{
						DemandID: "mobilefuse",
						Bid: &adapters.BidDemandResponse{
							Price:      0.7,
							ID:         "333",
							ImpID:      "444",
							DemandID:   adapter.MobileFuseKey,
							Signaldata: "signal_data",
						},
					},
					{
						DemandID: "vungle",
					},
				},
			}, nil
		},
	}
	auctionBuilderV2 := &auction.Builder{
		AdUnitsMatcher:               adUnitsMatcher,
		BiddingBuilder:               biddingBuilder,
		BiddingAdaptersConfigBuilder: biddingAdaptersConfigBuilder,
	}
	adapterKeyFetcher := &auctionmocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(ctx context.Context, appID int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionService := &auction.Service{
		AdapterKeysFetcher: adapterKeyFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilderV2,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        &event.Logger{Engine: &engine.Log{}},
	}

	handler := &apihandlers.AuctionHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.AuctionRequest, *schema.AuctionRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      gcoder,
		},
		AuctionService: auctionService,
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
			name:                 "OK",
			sdkVersion:           "0.5",
			requestPath:          "testdata/auction/valid_request.json",
			expectedResponsePath: "testdata/auction/valid_response.json",
			expectedStatusCode:   http.StatusOK,
			wantErr:              false,
		},
		{
			name:                 "OK",
			sdkVersion:           "0.5",
			requestPath:          "testdata/auction/valid_request_coppa.json",
			expectedResponsePath: "testdata/auction/valid_response_coppa.json",
			expectedStatusCode:   http.StatusOK,
			wantErr:              false,
		},
		{
			name:               "NoAdsFound",
			sdkVersion:         "0.5",
			requestPath:        "testdata/auction/noads_request.json",
			expectedStatusCode: http.StatusUnprocessableEntity,
			wantErr:            true,
			err:                sdkapi.ErrNoAdsFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := os.ReadFile(tt.requestPath)
			if err != nil {
				t.Fatalf("Error reading request file: %v", err)
			}

			handler := testHelperAuctionHandler()
			rec, err := ExecuteRequest(t, handler, http.MethodPost, "/v2/auction/interstitial", string(reqBody), &RequestOptions{
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

func TestAuctionHandler_EmptyResponseForNonAndroidMaxSDKVersions(t *testing.T) {
	tests := []struct {
		name          string
		requestFile   string
		endpoint      string
		sdkVersion    string
		shouldBeEmpty bool
		description   string
	}{
		{
			name:          "iOS_MAX_rewarded_0.7.0_should_return_empty",
			requestFile:   "testdata/auction/ios_max_mediator_rewarded_request.json",
			endpoint:      "/v2/auction/rewarded",
			sdkVersion:    "0.7.0",
			shouldBeEmpty: true,
			description:   "iOS with MAX mediator, rewarded ad type and SDK version 0.7.0 should return empty ads",
		},
		{
			name:          "iOS_MAX_rewarded_0.7.5_should_return_empty",
			requestFile:   "testdata/auction/ios_max_mediator_rewarded_request.json",
			endpoint:      "/v2/auction/rewarded",
			sdkVersion:    "0.7.5",
			shouldBeEmpty: true,
			description:   "iOS with MAX mediator, rewarded ad type and SDK version 0.7.5 should return empty ads",
		},
		{
			name:          "iOS_MAX_rewarded_0.8.1_should_return_empty",
			requestFile:   "testdata/auction/ios_max_mediator_rewarded_request.json",
			endpoint:      "/v2/auction/rewarded",
			sdkVersion:    "0.8.1",
			shouldBeEmpty: true,
			description:   "iOS with MAX mediator, rewarded ad type and SDK version 0.8.1 should return empty ads",
		},
		{
			name:          "iOS_MAX_rewarded_0.6.0_should_not_be_empty",
			requestFile:   "testdata/auction/ios_max_mediator_rewarded_request.json",
			endpoint:      "/v2/auction/rewarded",
			sdkVersion:    "0.6.0",
			shouldBeEmpty: false,
			description:   "iOS with MAX mediator, rewarded ad type and SDK version 0.6.0 should NOT return empty ads",
		},
		{
			name:          "iOS_MAX_rewarded_0.8.0_should_not_be_empty",
			requestFile:   "testdata/auction/ios_max_mediator_rewarded_request.json",
			endpoint:      "/v2/auction/rewarded",
			sdkVersion:    "0.8.0",
			shouldBeEmpty: false,
			description:   "iOS with MAX mediator, rewarded ad type and SDK version 0.8.0 should NOT return empty ads",
		},
		{
			name:          "iOS_MAX_rewarded_0.8.2_should_not_be_empty",
			requestFile:   "testdata/auction/ios_max_mediator_rewarded_request.json",
			endpoint:      "/v2/auction/rewarded",
			sdkVersion:    "0.8.2",
			shouldBeEmpty: false,
			description:   "iOS with MAX mediator, rewarded ad type and SDK version 0.8.2 should NOT return empty ads",
		},
		{
			name:          "iOS_MAX_interstitial_0.7.0_should_not_be_empty",
			requestFile:   "testdata/auction/ios_max_mediator_interstitial_request.json",
			endpoint:      "/v2/auction/interstitial",
			sdkVersion:    "0.7.0",
			shouldBeEmpty: false,
			description:   "iOS with MAX mediator, interstitial ad type and SDK version 0.7.0 should NOT return empty ads",
		},
		{
			name:          "Android_MAX_rewarded_0.7.0_should_not_be_empty",
			requestFile:   "testdata/auction/android_max_mediator_rewarded_request.json",
			endpoint:      "/v2/auction/rewarded",
			sdkVersion:    "0.7.0",
			shouldBeEmpty: false,
			description:   "Android with MAX mediator, rewarded ad type and SDK version 0.7.0 should NOT return empty ads",
		},
		{
			name:          "iOS_other_mediator_rewarded_0.7.0_should_not_be_empty",
			requestFile:   "testdata/auction/ios_other_mediator_rewarded_request.json",
			endpoint:      "/v2/auction/rewarded",
			sdkVersion:    "0.7.0",
			shouldBeEmpty: false,
			description:   "iOS with other mediator, rewarded ad type and SDK version 0.7.0 should NOT return empty ads",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := os.ReadFile(tt.requestFile)
			if err != nil {
				t.Fatalf("Error reading request file: %v", err)
			}

			handler := testHelperAuctionHandler()

			// Extract ad_type from endpoint path
			adType := ""
			if strings.Contains(tt.endpoint, "/rewarded") {
				adType = "rewarded"
			} else if strings.Contains(tt.endpoint, "/interstitial") {
				adType = "interstitial"
			} else if strings.Contains(tt.endpoint, "/banner") {
				adType = "banner"
			}

			rec, err := ExecuteRequest(t, handler, http.MethodPost, tt.endpoint, string(reqBody), &RequestOptions{
				Headers: map[string]string{
					"X-Bidon-Version": tt.sdkVersion,
				},
				Params: map[string]string{
					"ad_type": adType,
				},
			})
			if tt.shouldBeEmpty {
				if err == nil {
					t.Fatalf("%s: Expected error to be returned, got nil", tt.description)
				}
				if echoErr, ok := err.(*echo.HTTPError); ok {
					if echoErr.Code != http.StatusUnprocessableEntity {
						t.Errorf("%s: Expected error code 422, got %d", tt.description, echoErr.Code)
					}
					if echoErr.Message != "No ads found" {
						t.Errorf("%s: Expected error message 'No ads found', got '%v'", tt.description, echoErr.Message)
					}
				} else {
					t.Errorf("%s: Expected echo.HTTPError, got %T: %v", tt.description, err, err)
				}
			} else {
				if err != nil {
					t.Fatalf("ExecuteRequest failed: %v", err)
				}
			}

			if !tt.shouldBeEmpty {
				if rec.Code != http.StatusOK {
					t.Errorf("%s: Expected status %d, got %d", tt.description, http.StatusOK, rec.Code)
				}

				var response auction.Response
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if len(response.AdUnits) == 0 {
					t.Errorf("%s: Expected non-empty AdUnits, got empty", tt.description)
				}
			}
		})
	}
}
