package notification_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/notification/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/google/go-cmp/cmp"
)

func TestHandler_HandleBiddingRound(t *testing.T) {
	ctx := context.Background()
	floor := float64(2)
	imp := &schema.Imp{BidFloor: &floor}
	responses := bidding.AuctionResult{
		Bids: []adapters.DemandResponse{
			{Bid: &adapters.BidDemandResponse{ID: "bid-1", ImpID: "imp-1", Price: 1.23}},
			{Bid: &adapters.BidDemandResponse{ID: "bid-2", ImpID: "imp-1", Price: 4.56}},
			{Bid: &adapters.BidDemandResponse{ID: "bid-3", ImpID: "imp-1", Price: 7.89}},
			{Bid: &adapters.BidDemandResponse{ID: "bid-4", ImpID: "imp-1", Price: 0.12}},
			{Error: fmt.Errorf("error-1")},
		},
	}
	expectedBids := []notification.Bid{
		{ID: "bid-2", ImpID: "imp-1", Price: 4.56},
		{ID: "bid-3", ImpID: "imp-1", Price: 7.89},
	}

	mockRepo := &mocks.AuctionResultRepoMock{
		CreateOrUpdateFunc: func(ctx context.Context, imp *schema.Imp, bids []notification.Bid) error {
			if diff := cmp.Diff(expectedBids, bids); diff != "" {
				t.Errorf("CreateOrUpdate() mismatched arguments (-want, +got)\n%s", diff)
			}
			return nil
		},
	}
	wg := &sync.WaitGroup{}
	wg.Add(2) // We have 2 bids lower than floor, send 2 events

	sender := &mocks.SenderMock{SendEventFunc: func(ctx context.Context, p notification.Params) {
		wg.Done()
	}}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	err := handler.HandleBiddingRound(ctx, imp, responses, "bundle-1", "banner")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events, sent event lower than expected")
	}
}

func TestHandler_HandleStats_WinBid(t *testing.T) {
	ctx := context.Background()
	imp := schema.Stats{
		Result: schema.StatsResult{Status: "SUCCESS", ECPM: 7.89},
		Rounds: []schema.StatsRound{{
			ID:         "round-1",
			PriceFloor: 7.0,
			Demands: []schema.StatsDemand{{
				Price:  7.30,
				ID:     "bid-1",
				Status: "WIN",
			}},
			Bidding: schema.StatsBidding{
				Bids: []schema.StatsBid{
					{ID: "bid-1", Status: "LOSS", Price: 1.23},
					{ID: "bid-2", Status: "LOSS", Price: 4.56},
					{ID: "bid-3", Status: "WIN", Price: 7.89},
					{ID: "bid-4", Status: "NO_BID", Price: 0.12},
				},
			},
		}},
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

	handler := notification.Handler{AuctionResultRepo: mockRepo, Sender: sender}

	handler.HandleStats(ctx, imp, &config, "bundle-1", "banner")

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events, sent event lower than expected")
	}
}

func TestHandler_HandleStats_Loss(t *testing.T) {
	ctx := context.Background()
	imp := schema.Stats{
		Result: schema.StatsResult{
			Status: "SUCCESS",
			Price:  7.89,
		},
		Rounds: []schema.StatsRound{{
			ID:         "round-1",
			PriceFloor: 7.0,
			Demands: []schema.StatsDemand{{
				Price:  7.89,
				ID:     "bid-1",
				Status: "WIN",
			}},
			Bidding: schema.StatsBidding{
				Bids: []schema.StatsBid{{
					ID:     "bid-2",
					Status: "LOSS",
					Price:  7.6,
				}},
			},
		}},
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

	handler := notification.Handler{
		AuctionResultRepo: repoMock,
		Sender:            sender,
	}

	handler.HandleStats(ctx, imp, &config, "bundle-1", "banner")

	if waitTimeout(wg, 1*time.Second) {
		t.Errorf("timeout waiting for events, sent event lower than expected")
	}
}

func TestHandler_HandleShow_BiddingImpression(t *testing.T) {
	ctx := context.Background()
	impression := &schema.Bid{AuctionID: "auction-1", Price: 1.23, BidType: schema.RTBBidType}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return &notification.AuctionResult{
				Rounds: []notification.Round{{
					Bids: []notification.Bid{{
						ID:    "bid-1",
						ImpID: "imp-1",
						Price: 1.23,
					}},
				}},
			}, nil
		},
	}

	sender := &mocks.SenderMock{SendEventFunc: func(ctx context.Context, p notification.Params) {
		if p.NotificationType != "BURL" {
			t.Errorf("expected BURL notification, got %s", p.NotificationType)
		}
		if p.Bid.Price != impression.GetPrice() {
			t.Errorf("expected price %f, got %f", impression.GetPrice(), p.Bid.Price)
		}
	}}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	handler.HandleShow(ctx, impression, "bundle-1", "banner")
}

func TestHandler_HandleShow_NonBiddingImpression(t *testing.T) {
	ctx := context.Background()
	impression := &schema.Bid{AuctionID: "auction-1", Price: 1.23, BidType: schema.CPMBidType}

	mockRepo := &mocks.AuctionResultRepoMock{}
	sender := &mocks.SenderMock{
		SendEventFunc: func(ctx context.Context, p notification.Params) {
			t.Errorf("expected no notification, got %s", p.NotificationType)
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	handler.HandleShow(ctx, impression, "bundle-1", "banner")
}

func TestHandler_HandleShow_AuctionResultNotFound(t *testing.T) {
	ctx := context.Background()
	impression := &schema.Bid{AuctionID: "auction-1", Price: 1.23, BidType: schema.RTBBidType}

	mockRepo := &mocks.AuctionResultRepoMock{
		FindFunc: func(ctx context.Context, auctionID string) (*notification.AuctionResult, error) {
			return nil, nil
		},
	}

	sender := &mocks.SenderMock{
		SendEventFunc: func(ctx context.Context, p notification.Params) {
			t.Errorf("expected no notification, got %s", p.NotificationType)
		},
	}

	handler := notification.Handler{
		AuctionResultRepo: mockRepo,
		Sender:            sender,
	}

	handler.HandleShow(ctx, impression, "bundle-1", "banner")
}

func TestHandler_HandleWin(t *testing.T) {
	ctx := context.Background()
	imp := &schema.Bid{}

	handler := notification.Handler{}

	err := handler.HandleWin(ctx, imp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandler_HandleLoss(t *testing.T) {
	ctx := context.Background()
	imp := &schema.Bid{}

	handler := notification.Handler{}

	err := handler.HandleLoss(ctx, imp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
