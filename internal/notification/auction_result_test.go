package notification_test

import (
	"testing"

	"github.com/bidon-io/bidon-backend/internal/notification"
)

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
