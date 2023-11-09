package notification

import (
	"encoding/json"

	"github.com/bidon-io/bidon-backend/internal/adapter"
)

type AuctionResult struct {
	AuctionID string  `json:"auction_id"`
	Rounds    []Round `json:"rounds"`
}

func (a *AuctionResult) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *AuctionResult) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *AuctionResult) WinningBid() float64 {
	maxPrice := 0.0
	for _, round := range a.Rounds {
		for _, bid := range round.Bids {
			if bid.Price > maxPrice {
				maxPrice = bid.Price
			}
		}
	}
	return maxPrice
}

type Round struct {
	RoundID  string  `json:"round_id"`
	Bids     []Bid   `json:"bids"`
	BidFloor float64 `json:"bidfloor"`
}

type Bid struct {
	ID        string      `json:"id"`
	ImpID     string      `json:"impid"`
	Price     float64     `json:"price"`
	DemandID  adapter.Key `json:"demand_id"`
	AdID      string      `json:"adid"`
	SeatID    string      `json:"seatid"`
	LURL      string      `json:"lurl"`
	NURL      string      `json:"nurl"`
	BURL      string      `json:"burl"`
	RequestID string      `json:"request_id"`
}
