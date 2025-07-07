package schema

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
)

type AdObject struct {
	AuctionID               string                         `json:"auction_id" validate:"required"`
	AuctionKey              string                         `json:"auction_key"`
	AuctionConfigurationID  int64                          `json:"auction_configuration_id"`
	AuctionConfigurationUID string                         `json:"auction_configuration_uid"`
	PriceFloor              float64                        `json:"auction_pricefloor" validate:"gte=0"`
	Orientation             string                         `json:"orientation" validate:"oneof=PORTRAIT LANDSCAPE"`
	Demands                 map[adapter.Key]map[string]any `json:"demands"`
	Banner                  *BannerAdObject                `json:"banner"`
	Interstitial            *InterstitialAdObject          `json:"interstitial"`
	Rewarded                *RewardedAdObject              `json:"rewarded"`
}

func (o *AdObject) Format() ad.Format {
	if o.Banner != nil {
		return o.Banner.Format
	}

	return ad.EmptyFormat
}

type BannerAdObject struct {
	Format ad.Format `json:"format" validate:"oneof=BANNER LEADERBOARD MREC ADAPTIVE"`
}

type InterstitialAdObject struct{}

type RewardedAdObject struct{}

func (o *AdObject) GetBidFloor() float64 {
	return o.PriceFloor
}

const MinBidFloor = 0.000001

// GetBidFloorForBidding returns bidfloor increased by fraction of cent, so the bid is always higher than the bidfloor
func (o *AdObject) GetBidFloorForBidding() float64 {
	return o.PriceFloor + MinBidFloor
}

func (o *AdObject) IsAdaptive() bool {
	return o.Format() == ad.AdaptiveFormat
}

func (o *AdObject) Type() ad.Type {
	if o.Rewarded != nil {
		return ad.RewardedType
	}

	if o.Interstitial != nil {
		return ad.InterstitialType
	}

	return ad.BannerType
}

func (o *AdObject) IsPortrait() bool {
	return o.Orientation == "PORTRAIT"
}
