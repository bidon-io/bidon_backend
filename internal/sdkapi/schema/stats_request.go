package schema

import "github.com/bidon-io/bidon-backend/internal/ad"

type StatsRequest struct {
	BaseRequest
	AdType ad.Type `param:"ad_type"`
	Stats  Stats   `json:"stats" validate:"required"`
}

func (r *StatsRequest) Map() map[string]any {
	m := r.BaseRequest.Map()

	m["ad_type"] = r.AdType
	m["stats"] = r.Stats.Map()

	return m
}

type Stats struct {
	AuctionID              string       `json:"auction_id" validate:"required"`
	AuctionConfigurationID int          `json:"auction_configuration_id" validate:"required"`
	Result                 StatsResult  `json:"result" validate:"required"`
	Rounds                 []StatsRound `json:"rounds" validate:"required"`
}

func (s Stats) Map() map[string]any {
	m := map[string]any{
		"auction_id":               s.AuctionID,
		"auction_configuration_id": s.AuctionConfigurationID,
		"result":                   s.Result.Map(),
		"rounds":                   sliceMap(s.Rounds),
	}

	return m
}

type StatsResult struct {
	Status          string  `json:"status" validate:"required,oneof=SUCCESS FAIL AUCTION_CANCELLED"`
	WinnerID        string  `json:"winner_id"`
	RoundID         string  `json:"round_id"`
	ECPM            float64 `json:"ecpm"`
	AuctionStartTS  int     `json:"auction_start_ts"`
	AuctionFinishTS int     `json:"auction_finish_ts"`
}

func (s StatsResult) Map() map[string]any {
	m := map[string]any{
		"status":            s.Status,
		"winner_id":         s.WinnerID,
		"ecpm":              s.ECPM,
		"auction_start_ts":  s.AuctionStartTS,
		"auction_finish_ts": s.AuctionFinishTS,
	}

	return m
}

func (s StatsResult) IsSuccess() bool {
	return s.Status == "SUCCESS"
}

type StatsRound struct {
	ID         string        `json:"id" validate:"required"`
	PriceFloor float64       `json:"pricefloor" validate:"required"`
	Demands    []StatsDemand `json:"demands" validate:"required"`
	Bidding    StatsBidding  `json:"bidding"`
	WinnerID   string        `json:"winner_id"`
	WinnerECPM float64       `json:"winner_ecpm"`
}

func (r StatsRound) Map() map[string]any {
	m := map[string]any{
		"id":          r.ID,
		"pricefloor":  r.PriceFloor,
		"demands":     sliceMap(r.Demands),
		"bidding":     r.Bidding.Map(),
		"winner_id":   r.WinnerID,
		"winner_ecpm": r.WinnerECPM,
	}

	return m
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

func (d StatsDemand) Map() map[string]any {
	m := map[string]any{
		"id":             d.ID,
		"status":         d.Status,
		"ad_unit_id":     d.AdUnitID,
		"ecpm":           d.ECPM,
		"bid_start_ts":   d.BidStartTS,
		"bid_finish_ts":  d.BidFinishTS,
		"fill_start_ts":  d.FillStartTS,
		"fill_finish_ts": d.FillFinishTS,
	}

	return m
}

type StatsBidding struct {
	BidStartTS  int        `json:"bid_start_ts" validate:"required"`
	BidFinishTS int        `json:"bid_finish_ts"`
	Bids        []StatsBid `json:"bids" validate:"required"`
}

func (b StatsBidding) Map() map[string]any {
	m := map[string]any{
		"bid_start_ts":  b.BidStartTS,
		"bid_finish_ts": b.BidFinishTS,
		"bids":          sliceMap(b.Bids),
	}

	return m
}

type StatsBid struct {
	ID           string  `json:"id" validate:"required"`
	Status       string  `json:"status" validate:"required"`
	ECPM         float64 `json:"ecpm" validate:"required"`
	FillStartTS  int     `json:"fill_start_ts"`
	FillFinishTS int     `json:"fill_finish_ts"`
}

func (b StatsBid) Map() map[string]any {
	m := map[string]any{
		"id":             b.ID,
		"status":         b.Status,
		"ecpm":           b.ECPM,
		"fill_start_ts":  b.FillStartTS,
		"fill_finish_ts": b.FillFinishTS,
	}

	return m
}

type mapper interface {
	Map() map[string]any
}

func sliceMap[E mapper](s []E) []map[string]any {
	m := make([]map[string]any, len(s))
	for i, e := range s {
		m[i] = e.Map()
	}

	return m
}
