package notification_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/notification/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

func TestHandlerV2_HandleStats_WinBid(t *testing.T) {
	ctx := context.Background()
	imp := schema.StatsV2{
		AuctionID:               "f26af577-869e-41cb-909e-4d3eba57a28b",
		AuctionPricefloor:       7.0,
		AuctionConfigurationID:  10,
		AuctionConfigurationUID: "1701972528521547776",
		Result: schema.AuctionResult{
			Status:            "SUCCESS",
			BidType:           schema.RTBBidType,
			Price:             7.89,
			WinnerDemandID:    "vungle",
			WinnerAdUnitUID:   "1633824270331281408",
			WinnerAdUnitLabel: "vungle_inter_mergeblock_ios_3",
		},
		AdUnits: []schema.AuctionAdUnitResult{
			{
				Price:       7.30,
				DemandID:    "applovin",
				BidType:     schema.CPMBidType,
				AdUnitUID:   "1633833116256829440",
				AdUnitLabel: "applovin_inter_mergeblock_ios_6",
				Status:      "WIN",
			},
			{
				Price:       1.23,
				DemandID:    "bigoads",
				BidType:     schema.RTBBidType,
				AdUnitUID:   "1633833116256829140",
				AdUnitLabel: "bigoads_inter_mergeblock_ios_5",
				Status:      "LOSS",
			},
			{
				DemandID: "meta",
				BidType:  schema.RTBBidType,
				Status:   "NO_BID",
			},
			{
				Price:       7.89,
				DemandID:    "vungle",
				BidType:     schema.RTBBidType,
				AdUnitUID:   "1633824270331281408",
				AdUnitLabel: "vungle_inter_mergeblock_ios_3",
				Status:      "WIN",
			},
		},
	}
	result := notification.AuctionResult{
		Rounds: []notification.Round{{
			Bids: []notification.Bid{
				{ID: "bid-1", ImpID: "imp-1", Price: 1.23},
				{ID: "bid-2", ImpID: "imp-1", Price: 4.56},
				{ID: "bid-3", ImpID: "imp-2", Price: 7.89},
				{ID: "bid-4", ImpID: "imp-1", Price: 0.12},
			},
		}},
	}
	config := auction.Config{ExternalWinNotifications: false}
	mockRepo := &mocks.AuctionResultRepoMock{}
	mockRepo.FindFunc = func(ctx context.Context, id string) (*notification.AuctionResult, error) {
		return &result, nil
	}
	wg := &sync.WaitGroup{}
	wg.Add(4) // We have 4 bid in the auction result, wait all

	sender := &mocks.SenderMock{SendEventFunc: func(ctx context.Context, p notification.Params) {
		defer wg.Done()

		if p.FirstPrice != 7.89 {
			t.Errorf("expected first price 7.89, got %f", p.FirstPrice)
		}

		if p.SecondPrice != 7.30 {
			t.Errorf("expected second price 7.30, got %f", p.SecondPrice)
		}
	}}

	handlerV2 := notification.HandlerV2{AuctionResultRepo: mockRepo, Sender: sender}

	handlerV2.HandleStats(ctx, imp, &config, "bundle-1", "banner")

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events, sent event lower than expected")
	}
}

func TestHandlerV2_HandleStats_Loss(t *testing.T) {
	ctx := context.Background()

	imp := schema.StatsV2{
		AuctionID:               "f26af577-869e-41cb-909e-4d3eba57a28b",
		AuctionPricefloor:       7.0,
		AuctionConfigurationID:  10,
		AuctionConfigurationUID: "1701972528521547776",
		Result: schema.AuctionResult{
			Status:            "SUCCESS",
			BidType:           schema.CPMBidType,
			Price:             7.89,
			WinnerDemandID:    "vungle",
			WinnerAdUnitUID:   "1633824270331281408",
			WinnerAdUnitLabel: "vungle_inter_mergeblock_ios_3",
		},
		AdUnits: []schema.AuctionAdUnitResult{
			{
				Price:       7.89,
				DemandID:    "dtexchange",
				BidType:     schema.CPMBidType,
				AdUnitUID:   "1633833116256829440",
				AdUnitLabel: "dtexchange_inter_mergeblock_ios_6",
				Status:      "WIN",
			},
			{
				Price:       7.6,
				DemandID:    "vungle",
				BidType:     schema.RTBBidType,
				AdUnitUID:   "1633824270331281408",
				AdUnitLabel: "vungle_inter_mergeblock_ios_3",
				Status:      "LOSS",
			},
		},
	}

	config := auction.Config{ExternalWinNotifications: false}
	repoMock := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return &notification.AuctionResult{
				Rounds: []notification.Round{{
					RoundID:  "round-1",
					BidFloor: 7.0,
					Bids: []notification.Bid{{
						ID:    "bid-2",
						ImpID: "imp-1",
						Price: 7.6,
					}},
				}},
			}, nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1) // We have 1 bid in the auction result, should wait single goroutine to finish

	sender := &mocks.SenderMock{SendEventFunc: func(ctx context.Context, p notification.Params) {
		wg.Done()
	}}

	handlerV2 := notification.HandlerV2{
		AuctionResultRepo: repoMock,
		Sender:            sender,
	}

	handlerV2.HandleStats(ctx, imp, &config, "bundle-1", "banner")

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events, sent event lower than expected")
	}
}
