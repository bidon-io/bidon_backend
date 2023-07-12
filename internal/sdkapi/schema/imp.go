package schema

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
)

type Imp struct {
	ID              string                         `json:"id" validate:"required,uuid4"`
	AuctionID       string                         `json:"auction_id" validate:"required,uuid4"`
	AuctionConfigID int64                          `json:"auction_configuration_id" validate:"required"`
	RoundID         string                         `json:"round_id" validate:"required"`
	Bidfloor        float64                        `json:"bidfloor" validate:"required"`
	Orientation     string                         `json:"orientation" validate:"oneof=PORTRAIT LANDSCAPE"`
	Demands         map[adapter.Key]map[string]any `json:"demands"`
	Banner          *BannerAdObject                `json:"banner"`
	Interstitial    *InterstitialAdObject          `json:"interstitial"`
	Rewarded        *RewardedAdObject              `json:"rewarded"`
}

func (o *Imp) Format() ad.Format {
	if o.Banner != nil {
		return o.Banner.Format
	}

	return ad.EmptyFormat
}
