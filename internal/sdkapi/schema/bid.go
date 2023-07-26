package schema

type Bid struct {
	AuctionID              string                `json:"auction_id" validate:"required"`
	AuctionConfigurationID int                   `json:"auction_configuration_id" validate:"required"`
	ImpID                  string                `json:"imp_id" validate:"required"`
	DemandID               string                `json:"demand_id" validate:"required"`
	RoundID                string                `json:"round_id"`
	AdUnitID               string                `json:"ad_unit_id"`
	ECPM                   float64               `json:"ecpm" validate:"required"`
	Banner                 *BannerAdObject       `json:"banner"`
	Interstitial           *InterstitialAdObject `json:"interstitial"`
	Rewarded               *RewardedAdObject     `json:"rewarded"`
}

func (b Bid) Map() map[string]any {
	m := map[string]any{
		"auction_id":               b.AuctionID,
		"auction_configuration_id": b.AuctionConfigurationID,
		"imp_id":                   b.ImpID,
		"demand_id":                b.DemandID,
		"round_id":                 b.RoundID,
		"ad_unit_id":               b.AdUnitID,
		"ecpm":                     b.ECPM,
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
