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
	AdUnitID                string                `json:"ad_unit_id"`    // Deprecated: use AdUnitUID instead
	LineItemUID             string                `json:"line_item_uid"` // Deprecated: use AdUnitUID instead
	AdUnitUID               string                `json:"ad_unit_uid"`
	AdUnitLabel             string                `json:"ad_unit_label"`
	ECPM                    float64               `json:"ecpm" validate:"required_without=Price"` // Deprecated: use Price instead
	Price                   float64               `json:"price" validate:"required_without=ECPM"`
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

func (b *Bid) GetPrice() float64 {
	if b.Price != 0 {
		return b.Price
	}
	return b.ECPM
}

func (b *Bid) Format() ad.Format {
	if b.Banner != nil {
		return b.Banner.Format
	}

	return ad.EmptyFormat
}
