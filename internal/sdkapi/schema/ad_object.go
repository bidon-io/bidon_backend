package schema

import "github.com/bidon-io/bidon-backend/internal/ad"

type AdObject struct {
	PlacementID  string                `json:"placement_id"`
	AuctionID    string                `json:"auction_id"`
	Orientation  string                `json:"orientation" validate:"omitempty,oneof=PORTRAIT LANDSCAPE"`
	PriceFloor   float64               `json:"pricefloor"`
	Banner       *BannerAdObject       `json:"banner"`
	Interstitial *InterstitialAdObject `json:"interstitial"`
	Rewarded     *RewardedAdObject     `json:"rewarded"`
}

func (o *AdObject) Map() map[string]any {
	m := map[string]any{
		"placement_id": o.PlacementID,
		"auction_id":   o.AuctionID,
		"orientation":  o.Orientation,
		"pricefloor":   o.PriceFloor,
	}

	if o.Banner != nil {
		m["banner"] = o.Banner.Map()
	}
	if o.Interstitial != nil {
		m["interstitial"] = o.Interstitial.Map()
	}
	if o.Rewarded != nil {
		m["rewarded"] = o.Rewarded.Map()
	}

	return m
}

func (o *AdObject) Format() ad.Format {
	if o.Banner != nil {
		return o.Banner.Format
	}

	return ad.EmptyFormat
}

type BannerAdObject struct {
	Format ad.Format `json:"format"`
}

func (o BannerAdObject) Map() map[string]any {
	return map[string]any{
		"format": o.Format,
	}
}

type InterstitialAdObject struct{}

func (o InterstitialAdObject) Map() map[string]any {
	return map[string]any{}
}

type RewardedAdObject struct{}

func (o RewardedAdObject) Map() map[string]any {
	return map[string]any{}
}
