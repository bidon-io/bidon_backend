package schema

import (
	"strconv"
	"strings"

	"github.com/bidon-io/bidon-backend/internal/ad"
)

type StatsRequest struct {
	BaseRequest
	AdType ad.Type `param:"ad_type"`
	Stats  Stats   `json:"stats" validate:"required"`
}

func (r *StatsRequest) NormalizeValues() {
	r.BaseRequest.NormalizeValues()

	// Some SDK versions can send lower case bid_type
	r.Stats.Result.BidType = BidType(strings.ToUpper(r.Stats.Result.BidType.String()))
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

type StatsResult struct {
	Status            string                `json:"status" validate:"required,oneof=SUCCESS FAIL AUCTION_CANCELLED"`
	WinnerID          string                `json:"winner_id"` // Deprecated: WinnerDemandID instead
	WinnerDemandID    string                `json:"winner_demand_id"`
	WinnerAdUnitUID   string                `json:"winner_ad_unit_uid"`
	WinnerAdUnitLabel string                `json:"winner_ad_unit_label"`
	RoundID           string                `json:"round_id"`
	ECPM              float64               `json:"ecpm"`
	Price             float64               `json:"price"`
	BidType           BidType               `json:"bid_type" validate:"omitempty,oneof=RTB CPM"`
	AuctionStartTS    int64                 `json:"auction_start_ts"`
	AuctionFinishTS   int64                 `json:"auction_finish_ts"`
	Banner            *BannerAdObject       `json:"banner"`
	Interstitial      *InterstitialAdObject `json:"interstitial"`
	Rewarded          *RewardedAdObject     `json:"rewarded"`
}

func (s *StatsResult) GetWinnerDemandID() string {
	if s.WinnerDemandID != "" {
		return s.WinnerDemandID
	}
	return s.WinnerID
}

func (s *StatsResult) GetWinnerPrice() float64 {
	if s.Price != 0 {
		return s.Price
	}
	return s.ECPM
}

func (s *StatsResult) GetWinnerAdUnitUID() int {
	adUnitUID, err := strconv.Atoi(s.WinnerAdUnitUID)
	if err != nil {
		return 0
	}
	return adUnitUID
}

func (s *StatsResult) IsBidding() bool {
	return s.BidType == RTBBidType
}

func (s *StatsResult) IsSuccess() bool {
	return s.Status == "SUCCESS"
}

func (s *StatsResult) Format() ad.Format {
	if s.Banner != nil {
		return s.Banner.Format
	}

	return ad.EmptyFormat
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

type StatsDemand struct {
	ID           string  `json:"id" validate:"required"`
	Status       string  `json:"status" validate:"required"`
	AdUnitID     string  `json:"ad_unit_id"`
	LineItemUID  string  `json:"line_item_uid"` // Deprecated: use AdUnitID instead
	AdUnitUID    string  `json:"ad_unit_uid"`
	AdUnitLabel  string  `json:"ad_unit_label"`
	ECPM         float64 `json:"ecpm"` // Deprecated: use Price instead
	Price        float64 `json:"price"`
	FillStartTS  int64   `json:"fill_start_ts"`
	FillFinishTS int64   `json:"fill_finish_ts"`
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

type StatsBidding struct {
	BidStartTS  int64      `json:"bid_start_ts" validate:"required"`
	BidFinishTS int64      `json:"bid_finish_ts"`
	Bids        []StatsBid `json:"bids" validate:"required"`
}

type StatsBid struct {
	ID            string  `json:"id" validate:"required"`
	Status        string  `json:"status" validate:"required"`
	ECPM          float64 `json:"ecpm" validate:"required_without=Price"` // Deprecated: use Price instead
	Price         float64 `json:"price" validate:"required_without=ECPM"`
	AdUnitUID     string  `json:"ad_unit_uid"`
	AdUnitLabel   string  `json:"ad_unit_label"`
	FillStartTS   int64   `json:"fill_start_ts"`
	FillFinishTS  int64   `json:"fill_finish_ts"`
	TokenStartTS  int64   `json:"token_start_ts"`
	TokenFinishTS int64   `json:"token_finish_ts"`
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
