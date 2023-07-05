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

func (o *AdObject) Format() ad.Format {
	if o.Banner != nil {
		return o.Banner.Format
	}

	return ad.EmptyFormat
}

type BannerAdObject struct {
	Format ad.Format `json:"format"`
}

type InterstitialAdObject struct{}

type RewardedAdObject struct{}

