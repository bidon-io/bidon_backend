package schema

import "strconv"

type Bid struct {
	AuctionID               string                `json:"auction_id" validate:"required"`
	AuctionConfigurationID  int                   `json:"auction_configuration_id" validate:"required"`
	AuctionConfigurationUID string                `json:"auction_configuration_uid"`
	ImpID                   string                `json:"imp_id" validate:"required"`
	DemandID                string                `json:"demand_id" validate:"required"`
	RoundID                 string                `json:"round_id"`
	RoundIndex              int                   `json:"round_idx"`
	AdUnitID                string                `json:"ad_unit_id"`
	LineItemUID             string                `json:"line_item_uid"`
	ECPM                    float64               `json:"ecpm" validate:"required"`
	BidType                 string                `json:"bid_type" validate:"omitempty,oneof=rtb cpm"`
	Banner                  *BannerAdObject       `json:"banner"`
	Interstitial            *InterstitialAdObject `json:"interstitial"`
	Rewarded                *RewardedAdObject     `json:"rewarded"`
}

func (b Bid) IsBidding() bool {
	return b.BidType == "rtb"
}

func (b Bid) Map() map[string]any {
	auctionConfigurationUID, err := strconv.Atoi(b.AuctionConfigurationUID)
	if err != nil {
		auctionConfigurationUID = 0
	}
	lineItemUID, err := strconv.Atoi(b.LineItemUID)
	if err != nil {
		lineItemUID = 0
	}

	m := map[string]any{
		"auction_id":                b.AuctionID,
		"auction_configuration_id":  b.AuctionConfigurationID,
		"auction_configuration_uid": auctionConfigurationUID,
		"imp_id":                    b.ImpID,
		"demand_id":                 b.DemandID,
		"round_id":                  b.RoundID,
		"round_number":              b.RoundIndex,
		"ad_unit_id":                b.AdUnitID,
		"line_item_uid":             lineItemUID,
		"ecpm":                      b.ECPM,
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
