package auctionv2_test

import (
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/internal/ad"
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
	"testing"
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
			Regulations: &schema.Regulations{
				COPPA: true,
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
					{
						DemandID:   "applovin",
						UID:        "123_applovin",
						Label:      "applovin",
						PriceFloor: ptr(0.1),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": "123"},
					},
					{
						DemandID:   "unity",
						UID:        "123_unity",
						Label:      "unity",
						PriceFloor: ptr(0.3),
						BidType:    "CPM",
						Extra:      map[string]any{"placement_id": "123"},
					},
				},
				AdUnits: &[]auction.AdUnit{
					{
						DemandID: "bidmachine",
						UID:      "123_bidmachine",
						Label:    "bidmachine",
						BidType:  "RTB",
					},
					{
						DemandID: "mobilefuse",
						UID:      "123_mobilefuse",
						Label:    "mobilefuse",
						BidType:  "RTB",
					},
				},
				BiddingAuctionResult: &bidding.AuctionResult{
					Bids: []adapters.DemandResponse{
						{DemandID: "bidmachine", Bid: &adapters.BidDemandResponse{Price: 0.2, Payload: "token"}},
						{DemandID: "mobilefuse", Bid: &adapters.BidDemandResponse{Price: 0.5, Signaldata: "token"}},
						{DemandID: "meta", Bid: &adapters.BidDemandResponse{}},
					},
				},
			}, nil
		},
	}
	eventLogger := &event.Logger{Engine: &engine.Log{}}

	service := &auctionv2.Service{
		ConfigFetcher:  configFetcher,
		AuctionBuilder: auctionBuilder,
		SegmentMatcher: segmentMatcher,
		EventLogger:    eventLogger,
	}

	t.Run("Successful Run", func(t *testing.T) {
		responseUnits := []auction.AdUnit{
			{
				DemandID:   "mobilefuse",
				UID:        "123_mobilefuse",
				Label:      "mobilefuse",
				PriceFloor: ptr(0.5),
				BidType:    "RTB",
				Extra:      map[string]any{"signaldata": "token"},
			},
			{
				DemandID:   "bidmachine",
				UID:        "123_bidmachine",
				Label:      "bidmachine",
				PriceFloor: ptr(0.2),
				BidType:    "RTB",
				Extra:      map[string]any{"payload": "token"},
			},
			{
				DemandID:   "unity",
				UID:        "123_unity",
				Label:      "unity",
				PriceFloor: ptr(0.3),
				BidType:    "CPM",
				Extra:      map[string]any{"placement_id": "123"},
			},
			{
				DemandID:   "gam",
				UID:        "123_gam",
				Label:      "gam",
				PriceFloor: ptr(0.1),
				BidType:    "CPM",
				Extra:      map[string]any{"placement_id": "123"},
			},
		}
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
		if diff := cmp.Diff(response.AdUnits, responseUnits); diff != "" {
			t.Errorf("Expected \n(-want, +got)\n%s", diff)
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
