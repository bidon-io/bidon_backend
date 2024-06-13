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
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestBuilder_Build(t *testing.T) {
	auctionConfig := auction.Config{
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
		},
	}

	auctionConfigV2 := auction.Config{
		Demands: []adapter.Key{adapter.UnityAdsKey},
		Bidding: []adapter.Key{adapter.BidmachineKey},
		Timeout: 15000,
	}

	adaptersBuilder := &mocks.AdaptersBuilderMock{
		BuildFunc: func(adapterKey adapter.Key, cfg adapter.ProcessedConfigsMap) (*adapters.Bidder, error) {
			adpt := &bidmachine.BidmachineAdapter{
				Endpoint: "https://example.com",
				SellerID: "1",
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
		HandleBiddingRoundFunc: func(_ context.Context, _ *schema.Imp, _ bidding.AuctionResult, _ string, _ string) error {
			return nil
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
				AppID: 1,
				BiddingRequest: schema.BiddingRequest{
					Imp: schema.Imp{
						RoundID: "ROUND_2",
						Demands: map[adapter.Key]map[string]any{
							adapter.BidmachineKey: {
								"bid_token": "token",
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
				AuctionConfig: auctionConfig,
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
				AppID: 1,
				BiddingRequest: schema.BiddingRequest{
					Imp: schema.Imp{
						Demands: map[adapter.Key]map[string]any{
							adapter.BidmachineKey: {
								"bid_token": "token",
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
				AuctionConfig: auctionConfigV2,
			},
			expectedResult: adapters.DemandResponse{
				Status:   204,
				DemandID: adapter.BidmachineKey,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &bidding.Builder{
				AdaptersBuilder:     tt.adaptersBuilder,
				NotificationHandler: tt.notificationHandler,
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
