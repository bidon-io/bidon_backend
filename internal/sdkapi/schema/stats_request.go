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
	AuctionID               string                `json:"auction_id" validate:"required"`
	AuctionPricefloor       float64               `json:"auction_pricefloor"`
	AuctionConfigurationID  int64                 `json:"auction_configuration_id" validate:"required_without=AuctionConfigurationUID"`
	AuctionConfigurationUID string                `json:"auction_configuration_uid" validate:"required_without=AuctionConfigurationID"`
	Result                  AuctionResult         `json:"result" validate:"required"`
	AdUnits                 []AuctionAdUnitResult `json:"ad_units" validate:"required"`
}

type AuctionResult struct {
	Status            string                `json:"status" validate:"required,oneof=SUCCESS FAIL AUCTION_CANCELLED"`
	WinnerDemandID    string                `json:"winner_demand_id"`
	WinnerAdUnitUID   string                `json:"winner_ad_unit_uid"`
	WinnerAdUnitLabel string                `json:"winner_ad_unit_label"`
	Price             float64               `json:"price"`
	BidType           BidType               `json:"bid_type" validate:"omitempty,oneof=RTB CPM"`
	AuctionStartTS    int64                 `json:"auction_start_ts"`
	AuctionFinishTS   int64                 `json:"auction_finish_ts"`
	Banner            *BannerAdObject       `json:"banner"`
	Interstitial      *InterstitialAdObject `json:"interstitial"`
	Rewarded          *RewardedAdObject     `json:"rewarded"`
}

func (s *AuctionResult) GetWinnerDemandID() string {
	return s.WinnerDemandID
}

func (s *AuctionResult) GetWinnerPrice() float64 {
	return s.Price
}

func (s *AuctionResult) GetWinnerAdUnitUID() int {
	adUnitUID, err := strconv.Atoi(s.WinnerAdUnitUID)
	if err != nil {
		return 0
	}
	return adUnitUID
}

func (s *AuctionResult) IsBidding() bool {
	return s.BidType == RTBBidType
}

func (s *AuctionResult) IsSuccess() bool {
	return s.Status == "SUCCESS"
}

func (s *AuctionResult) Format() ad.Format {
	if s.Banner != nil {
		return s.Banner.Format
	}

	return ad.EmptyFormat
}

type AuctionAdUnitResult struct {
	Price         float64 `json:"price"`
	TokenStartTS  int64   `json:"token_start_ts"`
	TokenFinishTS int64   `json:"token_finish_ts"`
	FillStartTS   int64   `json:"fill_start_ts"`
	FillFinishTS  int64   `json:"fill_finish_ts"`
	DemandID      string  `json:"demand_id" validate:"required"`
	BidType       BidType `json:"bid_type" validate:"omitempty,oneof=RTB CPM"`
	AdUnitUID     string  `json:"ad_unit_uid"`
	AdUnitLabel   string  `json:"ad_unit_label"`
	Status        string  `json:"status" validate:"required"`
	ErrorMessage  string  `json:"error_message"`
}

func (r AuctionAdUnitResult) GetDemandID() string {
	return r.DemandID
}

func (r AuctionAdUnitResult) GetPrice() float64 {
	return r.Price
}

func (r AuctionAdUnitResult) GetAdUnitUID() int {
	adUnitUID, err := strconv.Atoi(r.AdUnitUID)
	if err != nil {
		return 0
	}
	return adUnitUID
}

func (r AuctionAdUnitResult) IsFill() bool {
	return r.Status == "WIN" || r.Status == "LOSE"
}
