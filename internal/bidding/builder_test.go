package bidding_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters/bidmachine"
	"github.com/bidon-io/bidon-backend/internal/bidding/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func TestBuilder_Build(t *testing.T) {
	configMatcher := &mocks.ConfigMatcherMock{
		MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error) {
			cfg := &auction.Config{
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
			return cfg, nil
		},
	}

	adaptersBuilder := &mocks.AdaptersBuilderMock{
		BuildFunc: func(adapterKey adapter.Key, cfg adapter.ProcessedConfigsMap) (adapters.Bidder, error) {
			adpt := &bidmachine.BidmachineAdapter{
				Endpoint: "https://example.com",
				SellerID: "1",
			}

			bidder := adapters.Bidder{
				Adapter: adpt,
				Client:  http.DefaultClient,
			}

			return bidder, nil
		},
	}

	notificationHanler := &mocks.NotificationHandlerMock{
		HandleRoundFunc: func(ctx context.Context, imp *schema.Imp, responses []adapters.DemandResponse) error {
			return nil
		},
	}

	tests := []struct {
		name                string
		configMatcher       bidding.ConfigMatcher
		adaptersBuilder     bidding.AdaptersBuilder
		notificationHandler bidding.NotificationHandler
		buildParams         *bidding.BuildParams
		expectedResult      adapters.DemandResponse
		expectedError       error
	}{
		{
			name:                "successful build",
			configMatcher:       configMatcher,
			adaptersBuilder:     adaptersBuilder,
			notificationHandler: notificationHanler,
			buildParams: &bidding.BuildParams{
				AppID:     1,
				SegmentID: 1,
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
			},
			expectedResult: adapters.DemandResponse{
				Price:    0,
				Status:   204,
				DemandID: adapter.BidmachineKey,
			},
			expectedError: nil,
		},
		{
			name: "config matcher error",
			configMatcher: &mocks.ConfigMatcherMock{
				MatchFunc: func(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error) {
					return nil, errors.New("config matcher error")
				},
			},
			buildParams: &bidding.BuildParams{
				AppID:     1,
				SegmentID: 1,
			},
			expectedResult: adapters.DemandResponse{
				Price: 0,
			},
			expectedError: errors.New("config matcher error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &bidding.Builder{
				ConfigMatcher:       tt.configMatcher,
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

			if result[0].Status != tt.expectedResult.Status {
				t.Errorf("expected result: %+v, but got %+v", tt.expectedResult, result)
			}
		})
	}
}
