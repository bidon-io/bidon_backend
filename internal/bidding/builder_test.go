package bidding_test

import (
	"context"
	"net/http"
	"testing"

	"go.uber.org/goleak"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/bidmachine"
	"github.com/bidon-io/bidon-backend/internal/bidding/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func testApp(id int64) *sdkapi.App {
	return &sdkapi.App{ID: id}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestBuilder_Build(t *testing.T) {
	auctionConfig := auction.Config{
		Bidding: []adapter.Key{adapter.BidmachineKey},
	}

	auctionConfigV2 := auction.Config{
		Demands: []adapter.Key{adapter.UnityAdsKey},
		Bidding: []adapter.Key{adapter.BidmachineKey},
		Timeout: 15000,
	}

	adaptersBuilder := &mocks.AdaptersBuilderMock{
		BuildFunc: func(_ adapter.Key, cfg adapter.ProcessedConfigsMap) (*adapters.Bidder, error) {
			adpt := &bidmachine.BidmachineAdapter{
				Endpoint: cfg[adapter.BidmachineKey]["endpoint"].(string),
				SellerID: cfg[adapter.BidmachineKey]["seller_id"].(string),
			}

			bidder := &adapters.Bidder{
				Adapter: adpt,
				Client:  http.DefaultClient,
			}

			return bidder, nil
		},
	}
	defer http.DefaultClient.CloseIdleConnections()

	notificationHandler := &mocks.NotificationHandlerMock{
		HandleBiddingRoundFunc: func(_ context.Context, _ *schema.AdObject, _ bidding.AuctionResult, _ string, _ string) error {
			return nil
		},
	}

	bidCacher := &mocks.BidCacherMock{ // Pass through bids
		ApplyBidCacheFunc: func(_ context.Context, _ *schema.AuctionRequest, aucRes *bidding.AuctionResult) []adapters.DemandResponse {
			return aucRes.Bids
		},
	}

	tests := []struct {
		name                string
		adaptersBuilder     bidding.AdaptersBuilder
		notificationHandler bidding.NotificationHandler
		buildParams         *bidding.BuildParams
		expectedResult      adapters.DemandResponse
		expectedError       error
	}{
		{
			name:                "successful build",
			adaptersBuilder:     adaptersBuilder,
			notificationHandler: notificationHandler,
			buildParams: &bidding.BuildParams{
				App: testApp(1),
				AdapterConfigs: adapter.ProcessedConfigsMap{
					adapter.BidmachineKey: {
						"endpoint":  "https://example.com",
						"seller_id": "1",
					},
				},
				AuctionRequest: schema.AuctionRequest{
					AdObject: schema.AdObject{
						Demands: map[adapter.Key]map[string]any{
							adapter.BidmachineKey: {
								"token": "token",
							},
						},
					},
					Adapters: schema.Adapters{
						adapter.BidmachineKey: {
							Version:    "1.0.0",
							SDKVersion: "1.0.0",
						},
					},
				},
				BiddingAdapters: auctionConfig.Bidding,
			},
			expectedResult: adapters.DemandResponse{
				Status:   204,
				DemandID: adapter.BidmachineKey,
			},
			expectedError: nil,
		},
		{
			name:                "round-less successful build v2",
			adaptersBuilder:     adaptersBuilder,
			notificationHandler: notificationHandler,
			buildParams: &bidding.BuildParams{
				App: testApp(1),
				AdapterConfigs: adapter.ProcessedConfigsMap{
					adapter.BidmachineKey: {
						"endpoint":  "https://example.com",
						"seller_id": "1",
					},
				},
				AuctionRequest: schema.AuctionRequest{
					AdObject: schema.AdObject{
						Demands: map[adapter.Key]map[string]any{
							adapter.BidmachineKey: {
								"token":           "bid_token",
								"status":          "SUCCESS",
								"token_start_ts":  169157953564,
								"token_finish_ts": 169157953564,
							},
						},
					},
					Adapters: schema.Adapters{
						adapter.BidmachineKey: {
							Version:    "1.0.0",
							SDKVersion: "1.0.0",
						},
					},
				},
				BiddingAdapters: auctionConfigV2.Bidding,
			},
			expectedResult: adapters.DemandResponse{
				Status:   204,
				DemandID: adapter.BidmachineKey,
			},
			expectedError: nil,
		},
		{
			name:                "round-less build v2 with no token",
			adaptersBuilder:     adaptersBuilder,
			notificationHandler: notificationHandler,
			buildParams: &bidding.BuildParams{
				App: testApp(1),
				AuctionRequest: schema.AuctionRequest{
					AdObject: schema.AdObject{
						Demands: map[adapter.Key]map[string]any{
							adapter.BidmachineKey: {"status": "SUCCESS"},
						},
					},
					Adapters: schema.Adapters{
						adapter.BidmachineKey: {
							Version:    "1.0.0",
							SDKVersion: "1.0.0",
						},
					},
				},
				BiddingAdapters: auctionConfigV2.Bidding,
			},
			expectedResult: adapters.DemandResponse{
				Status:   204,
				DemandID: adapter.BidmachineKey,
			},
			expectedError: bidding.ErrNoAdaptersMatched,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &bidding.Builder{
				AdaptersBuilder:     tt.adaptersBuilder,
				NotificationHandler: tt.notificationHandler,
				BidCacher:           bidCacher,
			}

			result, err := builder.HoldAuction(context.Background(), tt.buildParams)

			if err != nil && tt.expectedError == nil {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error: %v, but got nil", tt.expectedError)
			}

			if len(result.Bids) > 0 && result.Bids[0].Status != tt.expectedResult.Status {
				t.Errorf("expected result: %+v, but got %+v", tt.expectedResult, result)
			}
		})
	}
}
