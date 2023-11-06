package schema

import (
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/ad"
)

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

func (r *StatsRequest) GetAuctionConfigurationParams() (string, string) {
	return strconv.FormatInt(r.Stats.AuctionConfigurationID, 10), r.Stats.AuctionConfigurationUID
}

func (r *StatsRequest) SetAuctionConfigurationParams(id int64, uid string) {
	r.Stats.AuctionConfigurationID = id
	r.Stats.AuctionConfigurationUID = uid
}

type Stats struct {
	AuctionID               string       `json:"auction_id" validate:"required"`
	AuctionConfigurationID  int64        `json:"auction_configuration_id" validate:"required_without=AuctionConfigurationUID"`
	AuctionConfigurationUID string       `json:"auction_configuration_uid" validate:"required_without=AuctionConfigurationID"`
	Result                  StatsResult  `json:"result" validate:"required"`
	Rounds                  []StatsRound `json:"rounds" validate:"required"`
}

func (s Stats) Map() map[string]any {
	auctionConfigurationUID, err := strconv.Atoi(s.AuctionConfigurationUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	m := map[string]any{
		"auction_id":                s.AuctionID,
		"auction_configuration_id":  s.AuctionConfigurationID,
		"auction_configuration_uid": auctionConfigurationUID,
		"result":                    s.Result.Map(),
		"rounds":                    sliceMap(s.Rounds),
	}

	return m
}

type StatsResult struct {
	Status            string  `json:"status" validate:"required,oneof=SUCCESS FAIL AUCTION_CANCELLED"`
	WinnerID          string  `json:"winner_id"` // Deprecated: WinnerDemandID instead
	WinnerDemandID    string  `json:"winner_demand_id"`
	WinnerAdUnitUID   string  `json:"winner_ad_unit_uid"`
	WinnerAdUnitLabel string  `json:"winner_ad_unit_label"`
	RoundID           string  `json:"round_id"`
	ECPM              float64 `json:"ecpm"`
	Price             float64 `json:"price"`
	BidType           string  `json:"bid_type" validate:"omitempty,oneof=rtb cpm"`
	AuctionStartTS    int     `json:"auction_start_ts"`
	AuctionFinishTS   int     `json:"auction_finish_ts"`
}

func (s StatsResult) Map() map[string]any {
	m := map[string]any{
		"status":                 s.Status,
		"winner_id":              s.GetWinnerDemandID(),
		"winner_line_item_uid":   s.WinnerAdUnitUID,
		"winner_line_item_label": s.WinnerAdUnitLabel,
		"round_id":               s.RoundID,
		"ecpm":                   s.GetWinnerPrice(),
		"bid_type":               s.BidType,
		"bidding":                s.IsBidding(),
		"auction_start_ts":       s.AuctionStartTS,
		"auction_finish_ts":      s.AuctionFinishTS,
	}

	return m
}

func (s StatsResult) GetWinnerDemandID() string {
	if s.WinnerDemandID != "" {
		return s.WinnerDemandID
	}
	return s.WinnerID
}

func (s StatsResult) GetWinnerPrice() float64 {
	if s.Price != 0 {
		return s.Price
	}
	return s.ECPM
}

func (s StatsResult) GetWinnerAdUnitUID() int {
	adUnitUID, err := strconv.Atoi(s.WinnerAdUnitUID)
	if err != nil {
		return 0
	}
	return adUnitUID
}

func (s StatsResult) IsBidding() bool {
	return s.BidType == "rtb"
}

func (s StatsResult) IsSuccess() bool {
	return s.Status == "SUCCESS"
}

type StatsRound struct {
	ID                string        `json:"id" validate:"required"`
	PriceFloor        float64       `json:"pricefloor" validate:"required"`
	Demands           []StatsDemand `json:"demands" validate:"required"`
	Bidding           StatsBidding  `json:"bidding"`
	WinnerID          string        `json:"winner_id"`
	WinnerDemandID    string        `json:"winner_demand_id"`
	WinnerAdUnitUID   string        `json:"winner_ad_unit_uid"`
	WinnerAdUnitLabel string        `json:"winner_ad_unit_label"`
	WinnerECPM        float64       `json:"winner_ecpm"` // Deprecated: use WinnerPrice instead
	WinnerPrice       float64       `json:"winner_price"`
}

func (r StatsRound) GetWinnerDemandID() string {
	if r.WinnerDemandID != "" {
		return r.WinnerDemandID
	}
	return r.WinnerID
}

func (r StatsRound) GetWinnerPrice() float64 {
	if r.WinnerPrice != 0 {
		return r.WinnerPrice
	}
	return r.WinnerECPM
}

func (r StatsRound) GetWinnerAdUnitUID() int {
	adUnitUID, err := strconv.Atoi(r.WinnerAdUnitUID)
	if err != nil {
		return 0
	}
	return adUnitUID
}

func (r StatsRound) Map() map[string]any {
	m := map[string]any{
		"id":                     r.ID,
		"pricefloor":             r.PriceFloor,
		"demands":                sliceMap(r.Demands),
		"bidding":                r.Bidding.Map(),
		"winner_id":              r.GetWinnerDemandID(),
		"winner_ecpm":            r.GetWinnerPrice(),
		"winner_line_item_uid":   r.WinnerAdUnitUID,
		"winner_line_item_label": r.WinnerAdUnitLabel,
	}

	return m
}

type StatsDemand struct {
	ID           string  `json:"id" validate:"required"`
	Status       string  `json:"status" validate:"required"`
	AdUnitID     string  `json:"ad_unit_id"`
	LineItemUID  string  `json:"line_item_uid"` // Deprecated: use AdUnitID instead
	AdUnitUID    string  `json:"ad_unit_uid"`
	AdUnitLabel  string  `json:"ad_unit_label"`
	ECPM         float64 `json:"ecpm"` // Deprecated: use Price instead
	Price        float64 `json:"price"`
	BidStartTS   int     `json:"bid_start_ts"`
	BidFinishTS  int     `json:"bid_finish_ts"`
	FillStartTS  int     `json:"fill_start_ts"`
	FillFinishTS int     `json:"fill_finish_ts"`
}

func (d StatsDemand) GetPrice() float64 {
	if d.Price != 0 {
		return d.Price
	}
	return d.ECPM
}

func (d StatsDemand) GetAdUnitUID() int {
	var adUnitUIDStr string
	if d.AdUnitUID != "" {
		adUnitUIDStr = d.AdUnitUID
	} else {
		adUnitUIDStr = d.LineItemUID
	}

	lineItemUID, err := strconv.Atoi(adUnitUIDStr)
	if err != nil {
		return 0
	}
	return lineItemUID
}

func (d StatsDemand) Map() map[string]any {
	lineItemUID, err := strconv.Atoi(d.LineItemUID)
	if err != nil {
		lineItemUID = 0
	}

	m := map[string]any{
		"id":             d.ID,
		"status":         d.Status,
		"ad_unit_id":     d.AdUnitID,
		"line_item_uid":  lineItemUID,
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
	ECPM         float64 `json:"ecpm" validate:"required_without=Price"` // Deprecated: use Price instead
	Price        float64 `json:"price" validate:"required_without=ECPM"`
	AdUnitUID    string  `json:"ad_unit_uid"`
	AdUnitLabel  string  `json:"ad_unit_label"`
	FillStartTS  int     `json:"fill_start_ts"`
	FillFinishTS int     `json:"fill_finish_ts"`
}

func (b StatsBid) GetPrice() float64 {
	if b.Price != 0 {
		return b.Price
	}
	return b.ECPM
}

func (b StatsBid) GetAdUnitUID() int {
	adUnitUID, err := strconv.Atoi(b.AdUnitUID)
	if err != nil {
		return 0
	}
	return adUnitUID
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
