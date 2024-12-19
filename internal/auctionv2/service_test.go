package auctionv2_test

import (
	"context"
	"errors"

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
		ConfigFetcher:  configFetcher,
		AuctionBuilder: auctionBuilder,
		SegmentMatcher: segmentMatcher,
		EventLogger:    eventLogger,
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
