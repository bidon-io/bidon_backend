package schema

import (
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/ad"
)

type Bid struct {
	AuctionID               string                `json:"auction_id" validate:"required"`
	AuctionConfigurationID  int64                 `json:"auction_configuration_id" validate:"required_without=AuctionConfigurationUID"`
	AuctionConfigurationUID string                `json:"auction_configuration_uid" validate:"required_without=AuctionConfigurationID"`
	ImpID                   string                `json:"imp_id"`
	DemandID                string                `json:"demand_id" validate:"required"`
	RoundID                 string                `json:"round_id"`
	RoundIndex              int                   `json:"round_idx"`
	AdUnitUID               string                `json:"ad_unit_uid"`
	AdUnitLabel             string                `json:"ad_unit_label"`
	Price                   float64               `json:"price"`
	BidType                 BidType               `json:"bid_type" validate:"omitempty,oneof=RTB CPM"`
	AuctionPriceFloor       float64               `json:"auction_pricefloor"`
	Banner                  *BannerAdObject       `json:"banner"`
	Interstitial            *InterstitialAdObject `json:"interstitial"`
	Rewarded                *RewardedAdObject     `json:"rewarded"`
}

func (b *Bid) IsBidding() bool {
	return b.BidType == RTBBidType
}

func (b *Bid) GetAdUnitUID() int {
	adUnitUID, err := strconv.Atoi(b.AdUnitUID)
	if err != nil {
		return 0
	}
	return adUnitUID
}

func (b *Bid) GetPrice() float64 {
	return b.Price
}

func (b *Bid) Format() ad.Format {
	if b.Banner != nil {
		return b.Banner.Format
	}

	return ad.EmptyFormat
}
