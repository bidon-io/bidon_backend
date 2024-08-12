package apihandlers_test

import (
	"context"
	"errors"
	"github.com/bidon-io/bidon-backend/internal/auctionv2"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers"
	"net/http"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"

	auctionv2mocks "github.com/bidon-io/bidon-backend/internal/auctionv2/mocks"
	handlersmocks "github.com/bidon-io/bidon-backend/internal/sdkapi/v2/apihandlers/mocks"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
)

func testHelperAuctionV2Handler(t *testing.T) *apihandlers.AuctionHandler {
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
	gamPf := 0.8
	adUnits := []auction.AdUnit{
		{
			DemandID:   "amazon",
			Label:      "amazon",
			PriceFloor: &pf,
			UID:        "123_amazon",
			BidType:    schema.RTBBidType,
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
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
		{
			DemandID: "vungle",
			Label:    "vungle",
			UID:      "123_vungle",
			BidType:  schema.RTBBidType,
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
		{
			DemandID:   "gam",
			Label:      "gam",
			PriceFloor: &gamPf,
			UID:        "123_gam",
			BidType:    schema.CPMBidType,
			Extra: map[string]any{
				"placement_id": "123",
			},
		},
	}

	adUnitsMatcher := &auctionv2mocks.AdUnitsMatcherMock{
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
	biddingAdaptersConfigBuilder := &handlersmocks.AdaptersConfigBuilderMock{
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
	biddingBuilder := &handlersmocks.BiddingBuilderMock{
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
	auctionBuilderV2 := &auctionv2.Builder{
		ConfigFetcher:                configFetcher,
		AdUnitsMatcher:               adUnitsMatcher,
		BiddingBuilder:               biddingBuilder,
		BiddingAdaptersConfigBuilder: biddingAdaptersConfigBuilder,
	}

	handler := &apihandlers.AuctionHandler{
		BaseHandler: &apihandlers.BaseHandler[schema.AuctionV2Request, *schema.AuctionV2Request]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      gcoder,
		},
		AuctionBuilder: auctionBuilderV2,
		SegmentMatcher: segmentMatcher,
		EventLogger:    &event.Logger{Engine: &engine.Log{}},
	}

	return handler
}

func TestAuctionV2Handler_Handle(t *testing.T) {
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

			handler := testHelperAuctionV2Handler(t)
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
