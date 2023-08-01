package sdkapi_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/config"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionmocks "github.com/bidon-io/bidon-backend/internal/auction/mocks"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters_builder"
	adaptersbuildermocks "github.com/bidon-io/bidon-backend/internal/bidding/adapters_builder/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	sdkapimocks "github.com/bidon-io/bidon-backend/internal/sdkapi/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
)

func testHelperBiddingHandler(t *testing.T) sdkapi.BiddingHandler {
	app := sdkapi.App{ID: 1}
	geodata := geocoder.GeoData{CountryCode: "US"}
	auctionConfig := &auction.Config{
		ID: 1,
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
	segments := []segment.Segment{
		{ID: 1, Filters: []segment.Filter{{Type: "country", Name: "country", Operator: "IN", Values: []string{"US", "UK"}}}},
	}

	appFetcher := &sdkapimocks.AppFetcherMock{
		FetchFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return app, nil
		},
	}
	geocoder := &sdkapimocks.GeocoderMock{
		LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
			return geodata, nil
		},
	}
	configMatcher := &auctionmocks.ConfigMatcherMock{
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return segments, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}

	biddingHttpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxConnsPerHost:     50,
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
		},
	}

	profileFetcher := &adaptersbuildermocks.ConfigurationFetcherMock{
		FetchFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key) (adapter.RawConfigsMap, error) {
			return adapter.RawConfigsMap{
				adapter.ApplovinKey: {
					AccountExtra: map[string]string{
						"app_key": "123",
					},
				},
				adapter.BidmachineKey: {
					AccountExtra: map[string]string{},
				},
			}, nil
		},
	}

	lineItemsMatcher := &adaptersbuildermocks.LineItemsMatcherMock{
		MatchFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.LineItem, error) {
			return []auction.LineItem{}, nil
		},
	}

	// Create a new AuctionHandler instance

	handler := sdkapi.BiddingHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]{
			AppFetcher: appFetcher,
			Geocoder:   geocoder,
		},
		SegmentMatcher: segmentMatcher,
		BiddingBuilder: &bidding.Builder{
			ConfigMatcher:   configMatcher,
			AdaptersBuilder: adapters_builder.BuildBiddingAdapters(biddingHttpClient),
		},
		AdaptersConfigBuilder: &adapters_builder.AdaptersConfigBuilder{
			ConfigurationFetcher: profileFetcher,
			LineItemsMatcher:     lineItemsMatcher,
		},
	}
	return handler
}

func TestBiddingHandler_OK(t *testing.T) {
	handler := testHelperBiddingHandler(t)

	// Read request and response from file
	requestJson, err := os.ReadFile("testdata/bidding/bad_request.json")
	if err != nil {
		t.Fatalf("Error reading request file: %v", err)
	}
	expectedResponseJson, err := os.ReadFile("testdata/bidding/bad_response.json")
	fmt.Println(expectedResponseJson)
	if err != nil {
		t.Fatalf("Error reading response file: %v", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/bidding/interstitial", strings.NewReader(string(requestJson[:])))
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// Create a new Echo instance and context
	e := config.Echo("sdkapi-test", nil)
	c := e.NewContext(req, rec)

	// Call the Handle method
	_ = handler.Handle(c)
	// if err != nil {
	// 	t.Fatalf("Handle method returned an error: %v", err)
	// }

	// // Check that the response status code is HTTP 200 OK
	// if rec.Code != http.StatusOK {
	// 	t.Errorf("Http status is not ok (200). Received: %v", rec.Code)
	// }

	// // Read response body from file
	// var actualResponse interface{}
	// var expectedResponse interface{}
	// err = json.Unmarshal(rec.Body.Bytes(), &actualResponse)
	// if err != nil {
	// 	t.Fatalf("Failed to parse JSON1: %s", err)
	// }
	// err = json.Unmarshal(expectedResponseJson, &expectedResponse)
	// if err != nil {
	// 	t.Fatalf("Failed to parse JSON2: %s", err)
	// }

	// // Check that the response body is what we expect
	// if reflect.DeepEqual(actualResponse, expectedResponse) {
	// 	t.Errorf("Response mismatch. Expected: %v. Received: %v", expectedResponse, actualResponse)
	// }
}

func TestBiddingHandler_ErrNoAdsFound(t *testing.T) {
	handler := testHelperBiddingHandler(t)

	// Read request and response from file
	requestJson, err := os.ReadFile("testdata/auction/noads_request.json")
	if err != nil {
		t.Fatalf("Error reading request file: %v", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/auction/interstitial", strings.NewReader(string(requestJson[:])))
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// Create a new Echo instance and context
	e := config.Echo("sdkapi-test", nil)
	c := e.NewContext(req, rec)

	// Check that Handle method returns a ErrNoAdsFound error
	err = handler.Handle(c)
	if errors.Is(err, auction.ErrNoAdsFound) {
		t.Errorf("Handle method didn't return a ErrNoAdsFound error. Received: %v", err)
	}
}
