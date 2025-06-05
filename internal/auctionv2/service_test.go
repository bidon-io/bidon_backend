package auctionv2_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auctionv2"
	"github.com/bidon-io/bidon-backend/internal/auctionv2/mocks"
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

func TestService_Run(t *testing.T) {
	ctx := context.Background()
	auctionConfig := &auction.Config{
		ID:         1,
		UID:        "config_uid",
		PriceFloor: 0.05,
		Timeout:    15000,
	}
	geoData := geocoder.GeoData{}
	request := &schema.AuctionV2Request{
		AdObject: schema.AdObjectV2{
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
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, id, uid string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(ctx context.Context, appID int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(ctx context.Context, params *auctionv2.BuildParams) (*auctionv2.AuctionResult, error) {
			return &auctionv2.AuctionResult{
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

	service := &auctionv2.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	t.Run("Successful Run", func(t *testing.T) {
		params := &auctionv2.ExecutionParams{
			Req:     request,
			AppID:   1,
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(err error) {},
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
		params := &auctionv2.ExecutionParams{
			Req:     &invalidRequest,
			AppID:   1,
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(err error) {},
		}
		_, err := service.Run(ctx, params)
		if !errors.Is(err, sdkapi.ErrInvalidAuctionKey) {
			t.Fatalf("Expected error %v, got %v", sdkapi.ErrInvalidAuctionKey, err)
		}
	})

	t.Run("No Ads Found", func(t *testing.T) {
		auctionBuilder.BuildFunc = func(ctx context.Context, params *auctionv2.BuildParams) (*auctionv2.AuctionResult, error) {
			return nil, auction.ErrNoAdsFound
		}
		params := &auctionv2.ExecutionParams{
			Req:     request,
			AppID:   1,
			Country: "US",
			GeoData: geoData,
			Log:     func(string) {},
			LogErr:  func(err error) {},
		}
		_, err := service.Run(ctx, params)
		if !errors.Is(err, sdkapi.ErrNoAdsFound) {
			t.Fatalf("Expected error %v, got %v", sdkapi.ErrNoAdsFound, err)
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
	request := &schema.AuctionV2Request{
		AdObject: schema.AdObjectV2{
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
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, id, uid string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(ctx context.Context, appID int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(ctx context.Context, params *auctionv2.BuildParams) (*auctionv2.AuctionResult, error) {
			return &auctionv2.AuctionResult{
				AuctionConfiguration: auctionConfig,
				CPMAdUnits: &[]auction.AdUnit{
					{
						DemandID:   string(adapter.BidmachineKey),
						UID:        "123_bidmachine",
						Label:      "bidmachine",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": "123"},
					},
				},
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{},
				},
			}, nil
		},
	}
	eventLogger := &event.Logger{Engine: &engine.Log{}}

	service := &auctionv2.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	params := &auctionv2.ExecutionParams{
		Req:     request,
		AppID:   1,
		Country: "US",
		GeoData: geoData,
		Log:     func(string) {},
		LogErr:  func(err error) {},
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
	request := &schema.AuctionV2Request{
		AdObject: schema.AdObjectV2{
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
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, id, uid string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(ctx context.Context, appID int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(ctx context.Context, params *auctionv2.BuildParams) (*auctionv2.AuctionResult, error) {
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

			return &auctionv2.AuctionResult{
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

	service := &auctionv2.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	params := &auctionv2.ExecutionParams{
		Req:     request,
		AppID:   1,
		Country: "US",
		GeoData: geoData,
		Log:     func(string) {},
		LogErr:  func(err error) {},
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
	request := &schema.AuctionV2Request{
		AdObject: schema.AdObjectV2{
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
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, id, uid string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(ctx context.Context, appID int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(ctx context.Context, params *auctionv2.BuildParams) (*auctionv2.AuctionResult, error) {
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

			return &auctionv2.AuctionResult{
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

	service := &auctionv2.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	params := &auctionv2.ExecutionParams{
		Req:     request,
		AppID:   1,
		Country: "US",
		GeoData: geoData,
		Log:     func(string) {},
		LogErr:  func(err error) {},
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
	request := &schema.AuctionV2Request{
		AdObject: schema.AdObjectV2{
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
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return []segment.Segment{sgmnt}, nil
		},
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	configFetcher := &mocks.ConfigFetcherMock{
		FetchByUIDCachedFunc: func(ctx context.Context, appId int64, id, uid string) *auction.Config {
			return auctionConfig
		},
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error) {
			return auctionConfig, nil
		},
	}
	adapterKeysFetcher := &mocks.AdapterKeysFetcherMock{
		FetchEnabledAdapterKeysFunc: func(ctx context.Context, appID int64, keys []adapter.Key) ([]adapter.Key, error) {
			return keys, nil
		},
	}
	auctionBuilder := &mocks.AuctionBuilderMock{
		BuildFunc: func(ctx context.Context, params *auctionv2.BuildParams) (*auctionv2.AuctionResult, error) {
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

			return &auctionv2.AuctionResult{
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

	service := &auctionv2.Service{
		AdapterKeysFetcher: adapterKeysFetcher,
		ConfigFetcher:      configFetcher,
		AuctionBuilder:     auctionBuilder,
		SegmentMatcher:     segmentMatcher,
		EventLogger:        eventLogger,
	}

	params := &auctionv2.ExecutionParams{
		Req:     request,
		AppID:   1,
		Country: "US",
		GeoData: geoData,
		Log:     func(string) {},
		LogErr:  func(err error) {},
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
