package notification_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/notification/mocks"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/prebid/openrtb/v19/openrtb3"
)

func TestHandler_HandleRound(t *testing.T) {
	ctx := context.Background()
	floor := float64(2)
	imp := &schema.Imp{ID: "imp-1", BidFloor: &floor}
	responses := []adapters.DemandResponse{
		{Bid: &adapters.BidDemandResponse{ID: "bid-1", ImpID: "imp-1", Price: 1.23}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-2", ImpID: "imp-1", Price: 4.56}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-3", ImpID: "imp-1", Price: 7.89}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-4", ImpID: "imp-1", Price: 0.12}},
		{Error: fmt.Errorf("error-1")},
	}
	expectedBids := []notification.Bid{
		{ID: "bid-2", ImpID: "imp-1", Price: 4.56},
		{ID: "bid-3", ImpID: "imp-1", Price: 7.89},
	}
	mockRepo := &mocks.AuctionResultRepoMock{}
	mockRepo.CreateOrUpdateFunc = func(ctx context.Context, imp *schema.Imp, bids []notification.Bid) error {
		if diff := cmp.Diff(expectedBids, bids); diff != "" {
			t.Errorf("CreateOrUpdate() mismatched arguments (-want, +got)\n%s", diff)
		}
		return nil
	}

	handler := notification.Handler{AuctionResultRepo: mockRepo}

	err := handler.HandleRound(ctx, imp, responses)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandler_HandleStats(t *testing.T) {
	ctx := context.Background()
	imp := schema.Stats{Result: schema.StatsResult{Status: "SUCCESS", ECPM: 7.89}}
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

	handler := notification.Handler{AuctionResultRepo: mockRepo}

	err := handler.HandleStats(ctx, imp, config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandler_HandleStats_notification_true(t *testing.T) {
	ctx := context.Background()
	imp := schema.Stats{}

	config := auction.Config{ExternalWinNotifications: true}

	handler := notification.Handler{}

	err := handler.HandleStats(ctx, imp, config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandler_HandleShow(t *testing.T) {
	ctx := context.Background()
	imp := &schema.Imp{ID: "imp-1"}
	adapters := []*adapters.DemandResponse{
		{Bid: &adapters.BidDemandResponse{ID: "bid-1", ImpID: "imp-1", Price: 1.23, BURL: "http://example.com/burl"}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-2", ImpID: "imp-1", Price: 4.56, BURL: "http://example.com/burl"}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-3", ImpID: "imp-2", Price: 7.89, BURL: "http://example.com/burl"}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-4", ImpID: "imp-1", Price: 0.12, BURL: "http://example.com/burl"}},
	}
	handler := notification.Handler{}

	err := handler.HandleShow(ctx, imp, adapters)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandler_HandleWin(t *testing.T) {
	ctx := context.Background()
	imp := &schema.Imp{ID: "imp-1"}
	adapters := []*adapters.DemandResponse{
		{Bid: &adapters.BidDemandResponse{ID: "bid-1", ImpID: "imp-1", Price: 1.23, NURL: "http://example.com/win"}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-2", ImpID: "imp-1", Price: 4.56, NURL: "http://example.com/win"}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-3", ImpID: "imp-2", Price: 7.89, NURL: "http://example.com/win"}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-4", ImpID: "imp-1", Price: 0.12, NURL: "http://example.com/win"}},
	}
	handler := notification.Handler{}

	err := handler.HandleWin(ctx, imp, adapters)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandler_HandleLoss(t *testing.T) {
	ctx := context.Background()
	imp := &schema.Imp{ID: "imp-1"}
	adapters := []*adapters.DemandResponse{
		{Bid: &adapters.BidDemandResponse{ID: "bid-1", ImpID: "imp-1", Price: 1.23, LURL: "http://example.com/loss"}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-2", ImpID: "imp-1", Price: 4.56, LURL: "http://example.com/loss"}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-3", ImpID: "imp-2", Price: 7.89, LURL: "http://example.com/loss"}},
		{Bid: &adapters.BidDemandResponse{ID: "bid-4", ImpID: "imp-1", Price: 0.12, LURL: "http://example.com/loss"}},
	}
	handler := notification.Handler{}

	err := handler.HandleLoss(ctx, imp, adapters)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandler_SendLoss(t *testing.T) {
	// Create a test context and input data
	ctx := context.Background()

	// Create a mock HTTP server to handle the LURL request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		params := r.URL.Query()
		if diff := cmp.Diff("request-1", params.Get("id")); diff != "" {
			t.Errorf("mismatched id (-want, +got)\n%s", diff)
		}
		if diff := cmp.Diff("4.56", params.Get("auction_price")); diff != "" {
			t.Errorf("mismatched winprice (-want, +got)\n%s", diff)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	bid := notification.Bid{
		RequestID: "request-1",
		ImpID:     "imp-1",
		Price:     1.23,
		LURL:      fmt.Sprintf("%s/lurl?auction_price=${AUCTION_PRICE}&id=${AUCTION_ID}", server.URL),
	}
	lossReason := openrtb3.LossBelowAuctionFloor
	winPrice := 4.56
	secondPrice := 3.00

	// Create a Handler instance with the mock HTTP client and server
	handler := notification.Handler{HttpClient: server.Client()}

	// Call the SendLoss method with the test context and input data
	handler.SendLoss(ctx, bid, lossReason, winPrice, secondPrice)
}

func TestAuctionResult_WinningBid(t *testing.T) {
	auctionResult := &notification.AuctionResult{
		AuctionID: "auction-1",
		Rounds: []notification.Round{
			{
				RoundID: "round-1",
				Bids: []notification.Bid{
					{ID: "bid-1", ImpID: "imp-1", Price: 1.23},
					{ID: "bid-2", ImpID: "imp-1", Price: 4.56},
					{ID: "bid-3", ImpID: "imp-2", Price: 7.89},
					{ID: "bid-4", ImpID: "imp-1", Price: 0.12},
				},
				BidFloor: 0.5,
			},
			{
				RoundID: "round-2",
				Bids: []notification.Bid{
					{ID: "bid-5", ImpID: "imp-1", Price: 2.34},
					{ID: "bid-6", ImpID: "imp-1", Price: 5.67},
					{ID: "bid-7", ImpID: "imp-2", Price: 8.9},
					{ID: "bid-8", ImpID: "imp-1", Price: 0.23},
				},
				BidFloor: 0.5,
			},
		},
	}

	winningBid := auctionResult.WinningBid()

	if winningBid != 8.9 {
		t.Errorf("expected winningBid 8.9, got %f", winningBid)
	}
}
