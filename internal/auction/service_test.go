package auction_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auction/mocks"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
)

// MockEventLogger captures logged events for testing
type MockEventLogger struct {
	LoggedEvents []event.Event
}

func (m *MockEventLogger) Produce(message event.LogMessage, _ func(error)) {
	// For testing, we'll unmarshal the message back to an event
	var adEvent event.AdEvent
	if err := json.Unmarshal(message.Value, &adEvent); err == nil {
		m.LoggedEvents = append(m.LoggedEvents, &adEvent)
	}
}

func (m *MockEventLogger) Ping(_ context.Context) error {
	return nil
}

func TestService_Run(t *testing.T) {
	ctx := context.Background()
	auctionConfig := &auction.Config{
		ID:         1,
		UID:        "config_uid",
		PriceFloor: 0.05,
		Timeout:    15000,
	}
	geoData := geocoder.GeoData{}
	request := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionKey: "1ERNSV33K4000",
			PriceFloor: 0.01,
		},
		BaseRequest: schema.BaseRequest{
			Device: schema.Device{
				OS:   "android",
				Type: "phone",
			},
		},
		AdType: ad.BannerType,
		AdCache: []schema.AdCacheObject{
			{Price: 0.02},
		},
	}
	sgmnt := segment.Segment{
		ID:  1,
		UID: "1",
		Filters: []segment.Filter{
			{Type: "country", Operator: "IN", Values: []string{"US"}},
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(_ context.Context, _ int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(_ context.Context, _ int64, _, _ string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(_ context.Context, _ int64, _ ad.Type, _ int64, _ string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(_ context.Context, _ int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(_ context.Context, _ *auction.BuildParams) (*auction.Result, error) {
			return &auction.Result{
				AuctionConfiguration: auctionConfig,
				CPMAdUnits: &[]auction.AdUnit{
					{
						DemandID:   "gam",
						UID:        "123_gam",
						Label:      "gam",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": "123"},
					},
				},
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{
						{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
					},
				},
			}, nil
		},
	}
	eventLogger := &event.Logger{Engine: &engine.Log{}}

	service := &auction.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	t.Run("Successful Run", func(t *testing.T) {
		params := &auction.ExecutionParams{
			Req:     request,
			App:     testApp(1),
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(_ error) {},
		}
		response, err := service.Run(ctx, params)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if response.ConfigID != auctionConfig.ID {
			t.Errorf("Expected ConfigID %d, got %d", auctionConfig.ID, response.ConfigID)
		}
	})

	t.Run("Invalid Auction Key", func(t *testing.T) {
		invalidRequest := *request
		invalidRequest.AdObject.AuctionKey = "invalid_key"
		params := &auction.ExecutionParams{
			Req:     &invalidRequest,
			App:     testApp(1),
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(_ error) {},
		}
		_, err := service.Run(ctx, params)
		if !errors.Is(err, sdkapi.ErrInvalidAuctionKey) {
			t.Fatalf("Expected error %v, got %v", sdkapi.ErrInvalidAuctionKey, err)
		}
	})

	t.Run("No Ads Found", func(t *testing.T) {
		auctionBuilder.BuildFunc = func(_ context.Context, _ *auction.BuildParams) (*auction.Result, error) {
			return nil, auction.ErrNoAdsFound
		}
		params := &auction.ExecutionParams{
			Req:     request,
			App:     testApp(1),
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(_ error) {},
		}
		_, err := service.Run(ctx, params)
		if !errors.Is(err, sdkapi.ErrNoAdsFound) {
			t.Fatalf("Expected error %v, got %v", sdkapi.ErrNoAdsFound, err)
		}
	})

	t.Run("AdapterKeysFetcher Error - Events Still Logged", func(t *testing.T) {
		// Create a mock event logger to capture events
		mockEventLogger := &MockEventLogger{}
		service.EventLogger = &event.Logger{Engine: mockEventLogger}

		// Make AdapterKeysFetcher return an error
		adapterKeysFetcher.FetchEnabledAdapterKeysFunc = func(_ context.Context, _ int64, _ []adapter.Key) ([]adapter.Key, error) {
			return nil, errors.New("adapter keys fetch failed")
		}

		params := &auction.ExecutionParams{
			Req:     request,
			App:     testApp(1),
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(_ error) {},
		}

		_, err := service.Run(ctx, params)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		// Verify that events were still logged despite the error
		if len(mockEventLogger.LoggedEvents) == 0 {
			t.Fatal("Expected events to be logged even on error, but no events were logged")
		}

		// Find the auction_request event
		var auctionEvent *event.AdEvent
		for _, ev := range mockEventLogger.LoggedEvents {
			if adEvent, ok := ev.(*event.AdEvent); ok && adEvent.EventType == "auction_request" {
				auctionEvent = adEvent
				break
			}
		}

		if auctionEvent == nil {
			t.Fatal("Expected auction_request event to be logged")
		}

		// Verify error status and message
		if auctionEvent.Status != event.ErrorAdRequestStatus {
			t.Errorf("Expected Status to be 'ERROR', got '%s'", auctionEvent.Status)
		}
		if auctionEvent.Error == "" {
			t.Error("Expected Error field to contain error message, got empty string")
		}
		if !strings.Contains(auctionEvent.Error, "adapter keys fetch failed") {
			t.Errorf("Expected Error field to contain 'adapter keys fetch failed', got '%s'", auctionEvent.Error)
		}
	})

	t.Run("Invalid Auction Key Error - Events Still Logged", func(t *testing.T) {
		// Create a mock event logger to capture events
		mockEventLogger := &MockEventLogger{}
		service.EventLogger = &event.Logger{Engine: mockEventLogger}

		// Reset AdapterKeysFetcher to success
		adapterKeysFetcher.FetchEnabledAdapterKeysFunc = func(_ context.Context, _ int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		}

		// Create request with invalid auction key
		invalidRequest := *request
		invalidRequest.AdObject.AuctionKey = "invalid_key"

		params := &auction.ExecutionParams{
			Req:     &invalidRequest,
			App:     testApp(1),
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(_ error) {},
		}

		_, err := service.Run(ctx, params)
		if !errors.Is(err, sdkapi.ErrInvalidAuctionKey) {
			t.Fatalf("Expected ErrInvalidAuctionKey, got %v", err)
		}

		// Verify that events were still logged despite the error
		if len(mockEventLogger.LoggedEvents) == 0 {
			t.Fatal("Expected events to be logged even on error, but no events were logged")
		}

		// Find the auction_request event
		var auctionEvent *event.AdEvent
		for _, ev := range mockEventLogger.LoggedEvents {
			if adEvent, ok := ev.(*event.AdEvent); ok && adEvent.EventType == "auction_request" {
				auctionEvent = adEvent
				break
			}
		}

		if auctionEvent == nil {
			t.Fatal("Expected auction_request event to be logged")
		}

		// Verify error status and message
		if auctionEvent.Status != event.ErrorAdRequestStatus {
			t.Errorf("Expected Status to be 'ERROR', got '%s'", auctionEvent.Status)
		}
		if auctionEvent.Error == "" {
			t.Error("Expected Error field to contain error message, got empty string")
		}
	})

	t.Run("Config Match Error - Events Still Logged", func(t *testing.T) {
		// Create a mock event logger to capture events
		mockEventLogger := &MockEventLogger{}
		service.EventLogger = &event.Logger{Engine: mockEventLogger}

		// Make ConfigFetcher.Match return an error
		configFetcher.MatchFunc = func(_ context.Context, _ int64, _ ad.Type, _ int64, _ string) (*auction.Config, error) {
			return nil, errors.New("config match failed")
		}

		// Create request without auction key to trigger Match call
		noKeyRequest := *request
		noKeyRequest.AdObject.AuctionKey = ""

		params := &auction.ExecutionParams{
			Req:     &noKeyRequest,
			App:     testApp(1),
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(_ error) {},
		}

		_, err := service.Run(ctx, params)
		if !errors.Is(err, sdkapi.ErrNoAdsFound) {
			t.Fatalf("Expected ErrNoAdsFound, got %v", err)
		}

		// Verify that events were still logged despite the error
		if len(mockEventLogger.LoggedEvents) == 0 {
			t.Fatal("Expected events to be logged even on error, but no events were logged")
		}

		// Find the auction_request event
		var auctionEvent *event.AdEvent
		for _, ev := range mockEventLogger.LoggedEvents {
			if adEvent, ok := ev.(*event.AdEvent); ok && adEvent.EventType == "auction_request" {
				auctionEvent = adEvent
				break
			}
		}

		if auctionEvent == nil {
			t.Fatal("Expected auction_request event to be logged")
		}

		// Verify error status and message
		if auctionEvent.Status != event.ErrorAdRequestStatus {
			t.Errorf("Expected Status to be 'ERROR', got '%s'", auctionEvent.Status)
		}
		if auctionEvent.Error == "" {
			t.Error("Expected Error field to contain error message, got empty string")
		}
	})

	t.Run("Auction Builder Error - Events Still Logged", func(t *testing.T) {
		// Create a mock event logger to capture events
		mockEventLogger := &MockEventLogger{}
		service.EventLogger = &event.Logger{Engine: mockEventLogger}

		// Reset ConfigFetcher to success
		configFetcher.MatchFunc = func(_ context.Context, _ int64, _ ad.Type, _ int64, _ string) (*auction.Config, error) {
			return auctionConfig, nil
		}

		// Make AuctionBuilder return an error
		auctionBuilder.BuildFunc = func(_ context.Context, _ *auction.BuildParams) (*auction.Result, error) {
			return nil, errors.New("auction build failed")
		}

		// Create request without auction key to trigger Match call
		noKeyRequest := *request
		noKeyRequest.AdObject.AuctionKey = ""

		params := &auction.ExecutionParams{
			Req:     &noKeyRequest,
			App:     testApp(1),
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(_ error) {},
		}

		_, err := service.Run(ctx, params)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		// Verify that events were still logged despite the error
		if len(mockEventLogger.LoggedEvents) == 0 {
			t.Fatal("Expected events to be logged even on error, but no events were logged")
		}

		// Find the auction_request event
		var auctionEvent *event.AdEvent
		for _, ev := range mockEventLogger.LoggedEvents {
			if adEvent, ok := ev.(*event.AdEvent); ok && adEvent.EventType == "auction_request" {
				auctionEvent = adEvent
				break
			}
		}

		if auctionEvent == nil {
			t.Fatal("Expected auction_request event to be logged")
		}

		// Verify error status and message
		if auctionEvent.Status != event.ErrorAdRequestStatus {
			t.Errorf("Expected Status to be 'ERROR', got '%s'", auctionEvent.Status)
		}
		if auctionEvent.Error == "" {
			t.Error("Expected Error field to contain error message, got empty string")
		}
		if !strings.Contains(auctionEvent.Error, "auction build failed") {
			t.Errorf("Expected Error field to contain 'auction build failed', got '%s'", auctionEvent.Error)
		}

		// Verify auction configuration is properly set even in error case
		if auctionEvent.AuctionConfigurationID != auctionConfig.ID {
			t.Errorf("Expected AuctionConfigurationID to be %d, got %d", auctionConfig.ID, auctionEvent.AuctionConfigurationID)
		}
	})

	t.Run("Successful Run - Events Logged with Correct Status", func(t *testing.T) {
		// Create a mock event logger to capture events
		mockEventLogger := &MockEventLogger{}
		service.EventLogger = &event.Logger{Engine: mockEventLogger}

		// Create a successful auction result
		successfulAuctionResult := &auction.Result{
			AuctionConfiguration: auctionConfig,
			CPMAdUnits: &[]auction.AdUnit{
				{
					DemandID:   "gam",
					UID:        "123_gam",
					Label:      "gam",
					PriceFloor: ptr(0.1),
					BidType:    "CPM",
					Extra:      map[string]any{"placement_id": "123"},
				},
			},
			BiddingAuctionResult: &bidding.AuctionResult{
				Bids: []adapters.DemandResponse{
					{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{}},
				},
			},
		}

		// Reset all mocks to success
		configFetcher.MatchFunc = func(_ context.Context, _ int64, _ ad.Type, _ int64, _ string) (*auction.Config, error) {
			return auctionConfig, nil
		}
		auctionBuilder.BuildFunc = func(_ context.Context, _ *auction.BuildParams) (*auction.Result, error) {
			return successfulAuctionResult, nil
		}

		// Create request without auction key to trigger Match call
		noKeyRequest := *request
		noKeyRequest.AdObject.AuctionKey = ""

		params := &auction.ExecutionParams{
			Req:     &noKeyRequest,
			App:     testApp(1),
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(_ error) {},
		}

		_, err := service.Run(ctx, params)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify that events were logged
		if len(mockEventLogger.LoggedEvents) == 0 {
			t.Fatal("Expected events to be logged, but no events were logged")
		}

		// Find the auction_request event
		var auctionEvent *event.AdEvent
		for _, ev := range mockEventLogger.LoggedEvents {
			if adEvent, ok := ev.(*event.AdEvent); ok && adEvent.EventType == "auction_request" {
				auctionEvent = adEvent
				break
			}
		}

		if auctionEvent == nil {
			t.Fatal("Expected auction_request event to be logged")
		}

		// Verify success status - should be "SUCCESS" for success (actual behavior)
		if auctionEvent.Status != "SUCCESS" {
			t.Errorf("Expected Status to be 'SUCCESS' for success, got '%s'", auctionEvent.Status)
		}
		// Verify Error field is empty for success
		if auctionEvent.Error != "" {
			t.Errorf("Expected Error field to be empty for success, got '%s'", auctionEvent.Error)
		}
	})
}

func TestService_Run_BidmachineWithMediator(t *testing.T) {
	ctx := context.Background()
	auctionConfig := &auction.Config{
		ID:         1,
		UID:        "config_uid",
		PriceFloor: 0.05,
		Timeout:    15000,
	}
	geoData := geocoder.GeoData{}
	request := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionKey: "1ERNSV33K4000",
			PriceFloor: 0.01,
		},
		BaseRequest: schema.BaseRequest{
			Device: schema.Device{
				OS:   "android",
				Type: "phone",
			},
			Ext: `{"mediator": "max"}`,
		},
		AdType:  ad.BannerType,
		AdCache: []schema.AdCacheObject{},
	}
	sgmnt := segment.Segment{
		ID:  1,
		UID: "1",
		Filters: []segment.Filter{
			{Type: "country", Operator: "IN", Values: []string{"US"}},
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(_ context.Context, _ int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(_ context.Context, _ int64, _, _ string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(_ context.Context, _ int64, _ ad.Type, _ int64, _ string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(_ context.Context, _ int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(_ context.Context, _ *auction.BuildParams) (*auction.Result, error) {
			return &auction.Result{
				AuctionConfiguration: auctionConfig,
				CPMAdUnits: &[]auction.AdUnit{
					{
						DemandID:   string(adapter.BidmachineKey),
						UID:        "123_bidmachine",
						Label:      "bidmachine",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement": "123"},
					},
				},
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{},
				},
			}, nil
		},
	}
	eventLogger := &event.Logger{Engine: &engine.Log{}}

	service := &auction.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	params := &auction.ExecutionParams{
		Req:     request,
		App:     testApp(1),
		Country: "US",
		GeoData: geoData,
		Log:     func(string) {},
		LogErr:  func(_ error) {},
	}

	request.NormalizeValues()

	response, err := service.Run(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	found := false
	for _, adUnit := range response.AdUnits {
		if adUnit.DemandID == string(adapter.BidmachineKey) {
			if customParams, ok := adUnit.Extra["custom_parameters"].(map[string]any); ok {
				if mediator, ok := customParams["mediator"].(string); ok && mediator == "max" {
					found = true
					break
				}
			}
		}
	}
	if !found {
		t.Error("Expected bidmachine ad unit to have custom_parameters with mediator")
	}
}

func TestService_Run_BiddingWithDemandExt(t *testing.T) {
	ctx := context.Background()
	auctionConfig := &auction.Config{
		ID:         1,
		UID:        "config_uid",
		PriceFloor: 0.05,
		Timeout:    15000,
	}
	geoData := geocoder.GeoData{}
	request := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionKey: "1ERNSV33K4000",
			PriceFloor: 0.01,
		},
		BaseRequest: schema.BaseRequest{
			Device: schema.Device{
				OS:   "android",
				Type: "phone",
			},
			Ext: `{"mediator": "max"}`,
		},
		AdType:  ad.BannerType,
		AdCache: []schema.AdCacheObject{},
	}
	sgmnt := segment.Segment{
		ID:  1,
		UID: "1",
		Filters: []segment.Filter{
			{Type: "country", Operator: "IN", Values: []string{"US"}},
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(_ context.Context, _ int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(_ context.Context, _ int64, _, _ string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(_ context.Context, _ int64, _ ad.Type, _ int64, _ string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(_ context.Context, _ int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(_ context.Context, _ *auction.BuildParams) (*auction.Result, error) {
			adUnits := []auction.AdUnit{
				{
					DemandID: string(adapter.BidmachineKey),
					UID:      "bidmachine_unit_123",
					Label:    "bidmachine_test",
					BidType:  schema.RTBBidType,
					Timeout:  30000,
					Extra:    map[string]any{"test_key": "test_value"},
				},
			}

			return &auction.Result{
				AuctionConfiguration: auctionConfig,
				CPMAdUnits:           &[]auction.AdUnit{},
				AdUnits:              &adUnits,
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{
						{
							DemandID: adapter.BidmachineKey,
							Bid: &adapters.BidDemandResponse{
								ID:      "bid123",
								ImpID:   "imp123",
								Price:   0.15,
								Payload: "test_payload",
							},
						},
					},
				},
			}, nil
		},
	}
	eventLogger := &event.Logger{Engine: &engine.Log{}}

	service := &auction.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	params := &auction.ExecutionParams{
		Req:     request,
		App:     testApp(1),
		Country: "US",
		GeoData: geoData,
		Log:     func(string) {},
		LogErr:  func(_ error) {},
	}

	request.NormalizeValues()

	response, err := service.Run(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.AdUnits) == 0 && len(response.NoBids) == 0 {
		t.Error("Expected at least one ad unit or no bid")
	}
}

func TestService_Run_BidmachineWithMediatorInBidding(t *testing.T) {
	ctx := context.Background()
	auctionConfig := &auction.Config{
		ID:         1,
		UID:        "config_uid",
		PriceFloor: 0.05,
		Timeout:    15000,
	}
	geoData := geocoder.GeoData{}
	request := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionKey: "1ERNSV33K4000",
			PriceFloor: 0.01,
		},
		BaseRequest: schema.BaseRequest{
			Device: schema.Device{
				OS:   "android",
				Type: "phone",
			},
			Ext: `{"mediator": "max"}`,
		},
		AdType:  ad.BannerType,
		AdCache: []schema.AdCacheObject{},
	}
	sgmnt := segment.Segment{
		ID:  1,
		UID: "1",
		Filters: []segment.Filter{
			{Type: "country", Operator: "IN", Values: []string{"US"}},
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(_ context.Context, _ int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(_ context.Context, _ int64, _, _ string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(_ context.Context, _ int64, _ ad.Type, _ int64, _ string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(_ context.Context, _ int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(_ context.Context, _ *auction.BuildParams) (*auction.Result, error) {
			adUnits := []auction.AdUnit{
				{
					DemandID: string(adapter.BidmachineKey),
					UID:      "bidmachine_unit_123",
					Label:    "bidmachine_test",
					BidType:  schema.RTBBidType,
					Timeout:  30000,
					Extra:    map[string]any{"test_key": "test_value"},
				},
			}

			return &auction.Result{
				AuctionConfiguration: auctionConfig,
				CPMAdUnits:           &[]auction.AdUnit{},
				AdUnits:              &adUnits,
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{
						{
							DemandID: adapter.BidmachineKey,
							Bid: &adapters.BidDemandResponse{
								ID:      "bid123",
								ImpID:   "imp123",
								Price:   0.15,
								Payload: "test_payload_bidmachine",
							},
						},
					},
				},
			}, nil
		},
	}
	eventLogger := &event.Logger{Engine: &engine.Log{}}

	service := &auction.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	params := &auction.ExecutionParams{
		Req:     request,
		App:     testApp(1),
		Country: "US",
		GeoData: geoData,
		Log:     func(string) {},
		LogErr:  func(_ error) {},
	}

	request.NormalizeValues()

	response, err := service.Run(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	found := false
	for _, adUnit := range response.AdUnits {
		if adUnit.DemandID == string(adapter.BidmachineKey) {
			if payload, ok := adUnit.Extra["payload"].(string); ok && payload == "test_payload_bidmachine" {
				if customParams, ok := adUnit.Extra["custom_parameters"].(map[string]any); ok {
					if mediator, ok := customParams["mediator"].(string); ok && mediator == "max" {
						found = true
						break
					}
				}
			}
		}
	}
	for _, adUnit := range response.NoBids {
		if adUnit.DemandID == string(adapter.BidmachineKey) {
			if payload, ok := adUnit.Extra["payload"].(string); ok && payload == "test_payload_bidmachine" {
				if customParams, ok := adUnit.Extra["custom_parameters"].(map[string]any); ok {
					if mediator, ok := customParams["mediator"].(string); ok && mediator == "max" {
						found = true
						break
					}
				}
			}
		}
	}
	if !found {
		t.Error("Expected bidmachine ad unit to have payload and custom_parameters with mediator")
	}
}

func TestService_Run_BuildDemandExtVariousAdapters(t *testing.T) {
	ctx := context.Background()
	auctionConfig := &auction.Config{
		ID:         1,
		UID:        "config_uid",
		PriceFloor: 0.05,
		Timeout:    15000,
	}
	geoData := geocoder.GeoData{}
	request := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionKey: "1ERNSV33K4000",
			PriceFloor: 0.01,
		},
		BaseRequest: schema.BaseRequest{
			Device: schema.Device{
				OS:   "android",
				Type: "phone",
			},
			Ext: `{}`,
		},
		AdType:  ad.BannerType,
		AdCache: []schema.AdCacheObject{},
	}
	sgmnt := segment.Segment{
		ID:  1,
		UID: "1",
		Filters: []segment.Filter{
			{Type: "country", Operator: "IN", Values: []string{"US"}},
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(_ context.Context, _ int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(_ context.Context, _ int64, _, _ string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(_ context.Context, _ int64, _ ad.Type, _ int64, _ string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(_ context.Context, _ int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(_ context.Context, _ *auction.BuildParams) (*auction.Result, error) {
			adUnits := []auction.AdUnit{
				{
					DemandID: string(adapter.AmazonKey),
					UID:      "amazon_unit_123",
					Label:    "amazon_test",
					BidType:  schema.RTBBidType,
					Timeout:  30000,
					Extra:    map[string]any{"slot_uuid": "amazon_slot"},
				},
				{
					DemandID: string(adapter.MobileFuseKey),
					UID:      "mobilefuse_unit_123",
					Label:    "mobilefuse_test",
					BidType:  schema.RTBBidType,
					Timeout:  30000,
					Extra:    map[string]any{"test_key": "test_value"},
				},
				{
					DemandID: string(adapter.VKAdsKey),
					UID:      "vkads_unit_123",
					Label:    "vkads_test",
					BidType:  schema.RTBBidType,
					Timeout:  30000,
					Extra:    map[string]any{"test_key": "test_value"},
				},
				{
					DemandID: "unknown_adapter",
					UID:      "unknown_unit_123",
					Label:    "unknown_test",
					BidType:  schema.RTBBidType,
					Timeout:  30000,
					Extra:    map[string]any{"test_key": "test_value"},
				},
			}

			return &auction.Result{
				AuctionConfiguration: auctionConfig,
				CPMAdUnits:           &[]auction.AdUnit{},
				AdUnits:              &adUnits,
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{
						{
							DemandID: adapter.AmazonKey,
							SlotUUID: "amazon_slot",
							Bid: &adapters.BidDemandResponse{
								ID:      "amazon_bid",
								ImpID:   "amazon_imp",
								Price:   0.12,
								Payload: "amazon_payload",
							},
						},
						{
							DemandID: adapter.MobileFuseKey,
							Bid: &adapters.BidDemandResponse{
								ID:         "mobilefuse_bid",
								ImpID:      "mobilefuse_imp",
								Price:      0.13,
								Payload:    "mobilefuse_payload",
								Signaldata: "mobilefuse_signal",
							},
						},
						{
							DemandID: adapter.VKAdsKey,
							Bid: &adapters.BidDemandResponse{
								ID:      "vkads_bid_123",
								ImpID:   "vkads_imp",
								Price:   0.14,
								Payload: "vkads_payload",
							},
						},
						{
							DemandID: "unknown_adapter",
							Bid: &adapters.BidDemandResponse{
								ID:      "unknown_bid",
								ImpID:   "unknown_imp",
								Price:   0.11,
								Payload: "unknown_payload",
							},
						},
					},
				},
			}, nil
		},
	}
	eventLogger := &event.Logger{Engine: &engine.Log{}}

	service := &auction.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	params := &auction.ExecutionParams{
		Req:     request,
		App:     testApp(1),
		Country: "US",
		GeoData: geoData,
		Log:     func(string) {},
		LogErr:  func(_ error) {},
	}

	request.NormalizeValues()

	response, err := service.Run(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	allAdUnits := append(response.AdUnits, response.NoBids...)

	amazonFound := false
	mobilefuseFound := false
	vkadsFound := false
	unknownFound := false

	for _, adUnit := range allAdUnits {
		switch adUnit.DemandID {
		case string(adapter.AmazonKey):
			amazonFound = true
		case string(adapter.MobileFuseKey):
			if signaldata, ok := adUnit.Extra["signaldata"].(string); ok && signaldata == "mobilefuse_signal" {
				mobilefuseFound = true
			}
		case string(adapter.VKAdsKey):
			if bidID, ok := adUnit.Extra["bid_id"].(string); ok && bidID == "vkads_bid_123" {
				vkadsFound = true
			}
		case "unknown_adapter":
			if payload, ok := adUnit.Extra["payload"].(string); ok && payload == "unknown_payload" {
				unknownFound = true
			}
		}
	}

	if !amazonFound {
		t.Error("Expected Amazon ad unit to be found")
	}
	if !mobilefuseFound {
		t.Error("Expected MobileFuse ad unit to have signaldata")
	}
	if !vkadsFound {
		t.Error("Expected VKAds ad unit to have bid_id")
	}
	if !unknownFound {
		t.Error("Expected unknown adapter ad unit to have payload")
	}
}

func TestBidmachineWithPlacementID(t *testing.T) {
	ctx := context.Background()
	auctionConfig := &auction.Config{
		ID:         1,
		UID:        "config_uid",
		PriceFloor: 0.05,
		Timeout:    15000,
	}
	geoData := geocoder.GeoData{}
	request := &schema.AuctionRequest{
		AdObject: schema.AdObject{
			AuctionKey: "1ERNSV33K4000",
			PriceFloor: 0.01,
		},
		BaseRequest: schema.BaseRequest{
			Device: schema.Device{
				OS:   "android",
				Type: "phone",
			},
		},
		AdType: ad.BannerType,
	}
	sgmnt := segment.Segment{
		ID:  1,
		UID: "1",
		Filters: []segment.Filter{
			{Type: "country", Operator: "IN", Values: []string{"US"}},
		},
	}
	segmentFetcher := &segmentmocks.FetcherMock{
		FetchCachedFunc: func(_ context.Context, _ int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(_ context.Context, _ int64, _, _ string) *auction.Config {
			return auctionConfig
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(_ context.Context, _ int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}

	placementID := "test_placement_123"
	priceFloor := 0.1
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(_ context.Context, _ *auction.BuildParams) (*auction.Result, error) {
			return &auction.Result{
				AuctionConfiguration: auctionConfig,
				CPMAdUnits:           &[]auction.AdUnit{},
				AdUnits: &[]auction.AdUnit{
					{
						DemandID:   "bidmachine",
						UID:        "123_bidmachine",
						Label:      "bidmachine",
						PriceFloor: &priceFloor,
						BidType:    schema.RTBBidType,
						Extra: map[string]any{
							"placement": placementID,
						},
					},
				},
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{
						{
							DemandID: adapter.BidmachineKey,
							Bid:      &adapters.BidDemandResponse{Payload: "test_payload", Price: 0.15}, // Higher than price floor
						},
					},
				},
			}, nil
		},
	}
	eventLogger := &event.Logger{Engine: &engine.Log{}}

	service := &auction.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	params := &auction.ExecutionParams{
		Req:     request,
		App:     testApp(1),
		Country: "US",
		GeoData: geoData,
		Log:     func(string) {},
		LogErr:  func(_ error) {},
	}

	response, err := service.Run(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the response contains the BidMachine ad unit with placement_id in ext
	if len(response.AdUnits) == 0 {
		t.Fatal("Expected at least one ad unit in response")
	}

	bidmachineAdUnit := response.AdUnits[0]
	if bidmachineAdUnit.DemandID != "bidmachine" {
		t.Errorf("Expected DemandID to be 'bidmachine', got '%s'", bidmachineAdUnit.DemandID)
	}

	// Check that placement is included in the ext object (from the Extra field)
	if placementValue, exists := bidmachineAdUnit.Extra["placement"]; !exists {
		t.Error("Expected placement to be present in ext object")
	} else if placementValue != placementID {
		t.Errorf("Expected placement to be '%s', got '%v'", placementID, placementValue)
	}

	// Check that payload is also present
	if payload, exists := bidmachineAdUnit.Extra["payload"]; !exists {
		t.Error("Expected payload to be present in ext object")
	} else if payload != "test_payload" {
		t.Errorf("Expected payload to be 'test_payload', got '%v'", payload)
	}
}
