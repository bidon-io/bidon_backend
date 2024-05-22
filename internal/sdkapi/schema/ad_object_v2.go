package schema

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
)

type AdObjectV2 struct {
	AuctionID               string                         `json:"auction_id" validate:"required,uuid4"`
	AuctionConfigurationID  int64                          `json:"auction_configuration_id"`
	AuctionConfigurationUID string                         `json:"auction_configuration_uid"`
	PriceFloor              float64                        `json:"pricefloor" validate:"required,gte=0"`
	Orientation             string                         `json:"orientation" validate:"oneof=PORTRAIT LANDSCAPE"`
	Demands                 map[adapter.Key]map[string]any `json:"demands"`
	Banner                  *BannerAdObject                `json:"banner"`
	Interstitial            *InterstitialAdObject          `json:"interstitial"`
	Rewarded                *RewardedAdObject              `json:"rewarded"`
}

func (o *AdObjectV2) ToImp(roundID string) Imp {
	return Imp{
		AuctionID:               o.AuctionID,
		AuctionConfigurationID:  o.AuctionConfigurationID,
		AuctionConfigurationUID: o.AuctionConfigurationUID,
		RoundID:                 roundID,
		BidFloor:                &o.PriceFloor,
		Orientation:             o.Orientation,
		Demands:                 o.Demands,
		Banner:                  o.Banner,
		Interstitial:            o.Interstitial,
		Rewarded:                o.Rewarded,
	}
}

func (o *AdObjectV2) Format() ad.Format {
	if o.Banner != nil {
		return o.Banner.Format
	}

	return ad.EmptyFormat
}
