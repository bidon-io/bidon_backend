package schema

import "github.com/bidon-io/bidon-backend/internal/ad"

type StatsRequest struct {
	BaseRequest
	AdType ad.Type `param:"ad_type"`
	Stats  Stats   `json:"stats" validate:"required"`
}

type Stats struct {
	AuctionID              string       `json:"auction_id" validate:"required"`
	AuctionConfigurationID int          `json:"auction_configuration_id" validate:"required"`
	Result                 StatsResult  `json:"result" validate:"required"`
	Rounds                 []StatsRound `json:"rounds" validate:"required"`
}

type StatsResult struct {
	Status          string  `json:"status" validate:"required,oneof=SUCCESS FAIL AUCTION_CANCELLED"`
	WinnerID        string  `json:"winner_id"`
	ECPM            float64 `json:"ecpm"`
	AuctionStartTS  int     `json:"auction_start_ts"`
	AuctionFinishTS int     `json:"auction_finish_ts"`
}

func (s StatsResult) IsSuccess() bool {
	return s.Status == "SUCCESS"
}

type StatsRound struct {
	ID         string        `json:"id" validate:"required"`
	PriceFloor float64       `json:"pricefloor" validate:"required"`
	Demands    []StatsDemand `json:"demands" validate:"required"`
	WinnerID   string        `json:"winner_id"`
	WinnerECPM float64       `json:"winner_ecpm"`
}

type StatsDemand struct {
	ID           string  `json:"id" validate:"required"`
	Status       string  `json:"status" validate:"required"`
	AdUnitID     string  `json:"ad_unit_id"`
	ECPM         float64 `json:"ecpm"`
	BidStartTS   int     `json:"bid_start_ts"`
	BidFinishTS  int     `json:"bid_finish_ts"`
	FillStartTS  int     `json:"fill_start_ts"`
	FillFinishTS int     `json:"fill_finish_ts"`
}
