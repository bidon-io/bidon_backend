package sdkapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bidon-io/bidon-backend/config"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionmocks "github.com/bidon-io/bidon-backend/internal/auction/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	sdkapimocks "github.com/bidon-io/bidon-backend/internal/sdkapi/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
)

func testHelperAuctionHandler(t *testing.T) *sdkapi.AuctionHandler {
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
	lineItems := []auction.LineItem{
		{ID: "test", PriceFloor: 0.1, PlacementID: "", AdUnitID: "test_id"},
	}
	segments := []segment.Segment{
		{ID: 1, Filters: []segment.Filter{{Type: "country", Name: "country", Operator: "IN", Values: []string{"US", "UK"}}}},
	}

	appFetcher := &sdkapimocks.AppFetcherMock{
		FetchFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return app, nil
		},
	}
	gcoder := &sdkapimocks.GeocoderMock{
		LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
			return geodata, nil
		},
	}
	configMatcher := &auctionmocks.ConfigMatcherMock{
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	lineItemsMatcher := &auctionmocks.LineItemsMatcherMock{
		MatchFunc: func(ctx context.Context, params *auction.BuildParams) ([]auction.LineItem, error) {
			return lineItems, nil
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return segments, nil
		},
	}
	auctionBuilder := &auction.Builder{
		ConfigMatcher:    configMatcher,
		LineItemsMatcher: lineItemsMatcher,
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}

	// Create a new AuctionHandler instance

	handler := &sdkapi.AuctionHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.AuctionRequest, *schema.AuctionRequest]{
			AppFetcher: appFetcher,
			Geocoder:   gcoder,
		},
		AuctionBuilder: auctionBuilder,
		SegmentMatcher: segmentMatcher,
		EventLogger:    &event.Logger{Engine: &engine.Log{}},
	}

	return handler
}

func TestAuctionHandler_OK(t *testing.T) {
	handler := testHelperAuctionHandler(t)

	// Read request and response from file
	requestJson, err := os.ReadFile("testdata/auction/valid_request.json")
	if err != nil {
		t.Fatalf("Error reading request file: %v", err)
	}
	expectedResponseJson, err := os.ReadFile("testdata/auction/valid_response.json")
	if err != nil {
		t.Fatalf("Error reading response file: %v", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/auction/interstitial", strings.NewReader(string(requestJson[:])))
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// Create a new Echo instance and context
	e := config.Echo("sdkapi-test", nil)
	c := e.NewContext(req, rec)

	// Call the Handle method
	err = handler.Handle(c)
	if err != nil {
		t.Fatalf("Handle method returned an error: %v", err)
	}

	// Check that the response status code is HTTP 200 OK
	if rec.Code != http.StatusOK {
		t.Errorf("Http status is not ok (200). Received: %v", rec.Code)
	}

	// Read response body from file
	var actualResponse interface{}
	var expectedResponse interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Failed to parse JSON1: %s", err)
	}
	err = json.Unmarshal(expectedResponseJson, &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to parse JSON2: %s", err)
	}

	// Check that the response body is what we expect
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Response mismatch. Expected: %v. Received: %v", expectedResponse, actualResponse)
	}
}

func TestAuctionHandler_ErrNoAdsFound(t *testing.T) {
	handler := testHelperAuctionHandler(t)

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
