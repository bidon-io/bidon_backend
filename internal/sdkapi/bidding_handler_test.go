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
	"time"

	"github.com/bidon-io/bidon-backend/internal/adapter/store"

	"github.com/bidon-io/bidon-backend/config"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters_builder"
	biddingmocks "github.com/bidon-io/bidon-backend/internal/bidding/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	sdkapimocks "github.com/bidon-io/bidon-backend/internal/sdkapi/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
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
				Bidding: []adapter.Key{adapter.MetaKey, adapter.AmazonKey},
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

	appFetcher := &sdkapimocks.AppFetcherMock{
		FetchCachedFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
			return app, nil
		},
	}
	geocoder := &sdkapimocks.GeocoderMock{
		LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
			return geodata, nil
		},
	}
	configFetcher := &sdkapimocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, key string, aucUID string) *auction.Config {
			return auctionConfig
		},
	}

	biddingHttpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxConnsPerHost:     50,
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
		},
	}
	pf := 0.1
	adUnitsMap := &store.AdUnitsMap{
		adapter.MetaKey: []auction.AdUnit{
			{
				DemandID:   "meta",
				Label:      "meta",
				PriceFloor: &pf,
				UID:        "123_meta",
				Extra: map[string]any{
					"placement_id": "123",
				},
			},
		},
	}

	adUnitsMapBuilder := &sdkapimocks.AdUnitsMapBuilderMock{
		BuildFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp) (store.AdUnitsMap, error) {
			return *adUnitsMap, nil
		},
	}

	adapterConfigBuilder := &sdkapimocks.AdaptersConfigBuilderMock{
		BuildFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp, adUnitsMap *store.AdUnitsMap) (adapter.ProcessedConfigsMap, error) {
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
				adapter.AmazonKey: map[string]any{
					"price_points_map": map[string]any{
						"price_point_1": map[string]any{"price": 0.1, "price_point": "price_point_1"},
					},
				},
			}, nil
		},
	}

	notificationMock := &biddingmocks.NotificationHandlerMock{
		HandleRoundFunc: func(context.Context, *schema.Imp, bidding.AuctionResult) error {
			return nil
		},
	}

	// Create a new AuctionHandler instance

	handler := sdkapi.BiddingHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geocoder,
		},
		BiddingBuilder: &bidding.Builder{
			AdaptersBuilder:     adapters_builder.BuildBiddingAdapters(biddingHttpClient),
			NotificationHandler: notificationMock,
		},
		AdaptersConfigBuilder: adapterConfigBuilder,
		AdUnitsMapBuilder:     adUnitsMapBuilder,
		EventLogger:           &event.Logger{Engine: &engine.Log{}},
	}
	return handler
}

func TestBiddingHandler_OK(t *testing.T) {
	handler := testHelperBiddingHandler(t)

	// Read request and response from file
	requestJson, err := os.ReadFile("testdata/bidding/ok_request.json")
	if err != nil {
		t.Fatalf("Error reading request file: %v", err)
	}
	expectedResponseJson, err := os.ReadFile("testdata/bidding/ok_response.json")
	if err != nil {
		t.Fatalf("Error reading response file: %v", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/bidding/interstitial", strings.NewReader(string(requestJson[:])))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Bidon-Version", "0.4.0")

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// Create a new Echo instance and context
	e := config.Echo()
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
	e := config.Echo()
	c := e.NewContext(req, rec)

	// Check that Handle method returns a ErrNoAdsFound error
	err = handler.Handle(c)
	if errors.Is(err, auction.ErrNoAdsFound) {
		t.Errorf("Handle method didn't return a ErrNoAdsFound error. Received: %v", err)
	}
}
