package schema

import (
	"strconv"
)

type BidType string

const (
	EmptyBidType BidType = ""
	RTBBidType   BidType = "RTB"
	CPMBidType   BidType = "CPM"
)

func (b BidType) String() string {
	return string(b)
}

type Bid struct {
	AuctionID               string                `json:"auction_id" validate:"required"`
	AuctionConfigurationID  int                   `json:"auction_configuration_id" validate:"required"`
	AuctionConfigurationUID string                `json:"auction_configuration_uid"`
	ImpID                   string                `json:"imp_id" validate:"required"`
	DemandID                string                `json:"demand_id" validate:"required"`
	RoundID                 string                `json:"round_id"`
	RoundIndex              int                   `json:"round_idx"`
	AdUnitID                string                `json:"ad_unit_id"`    // Deprecated: use AdUnitUID instead
	LineItemUID             string                `json:"line_item_uid"` // Deprecated: use AdUnitUID instead
	AdUnitUID               string                `json:"ad_unit_uid"`
	AdUnitLabel             string                `json:"ad_unit_label"`
	ECPM                    float64               `json:"ecpm" validate:"required_without=Price"` // Deprecated: use Price instead
	Price                   float64               `json:"price" validate:"required_without=ECPM"`
	BidType                 BidType               `json:"bid_type" validate:"omitempty,oneof=RTB CPM"`
	RoundPriceFloor         float64               `json:"round_price_floor"`
	AuctionPriceFloor       float64               `json:"auction_price_floor"`
	Banner                  *BannerAdObject       `json:"banner"`
	Interstitial            *InterstitialAdObject `json:"interstitial"`
	Rewarded                *RewardedAdObject     `json:"rewarded"`
}

func (b Bid) IsBidding() bool {
	return b.BidType == RTBBidType
}

func (b Bid) GetAdUnitUID() int {
	var adUnitUIDStr string
	if b.AdUnitUID != "" {
		adUnitUIDStr = b.AdUnitUID
	} else {
		adUnitUIDStr = b.LineItemUID
	}

	adUnitUID, err := strconv.Atoi(adUnitUIDStr)
	if err != nil {
		return 0
	}
	return adUnitUID
}

func (b Bid) GetPrice() float64 {
	if b.Price != 0 {
		return b.Price
	}
	return b.ECPM
}

func (b Bid) Map() map[string]any {
	auctionConfigurationUID, err := strconv.Atoi(b.AuctionConfigurationUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	m := map[string]any{
		"auction_id":                b.AuctionID,
		"auction_configuration_id":  b.AuctionConfigurationID,
		"auction_configuration_uid": auctionConfigurationUID,
		"imp_id":                    b.ImpID,
		"demand_id":                 b.DemandID,
		"round_id":                  b.RoundID,
		"round_number":              b.RoundIndex,
		"ad_unit_id":                b.AdUnitID, // Deprecated: use AdUnitUID instead
		"line_item_uid":             b.GetAdUnitUID(),
		"line_item_label":           b.AdUnitLabel,
		"ecpm":                      b.GetPrice(),
		"round_price_floor":         b.RoundPriceFloor,
		"auction_price_floor":       b.AuctionPriceFloor,
		"bid_type":                  b.BidType,
		"bidding":                   b.IsBidding(),
	}

	if b.Banner != nil {
		m["banner"] = b.Banner.Map()
	}
	if b.Interstitial != nil {
		m["interstitial"] = b.Interstitial.Map()
	}
	if b.Rewarded != nil {
		m["rewarded"] = b.Rewarded.Map()
	}

	return m
}
