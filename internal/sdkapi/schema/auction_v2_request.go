package schema

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
)

// AuctionV2Request Has both the auction and bidding request fields
type AuctionV2Request struct {
	BaseRequest
	AdType     ad.Type    `param:"ad_type"`
	Adapters   Adapters   `json:"adapters" validate:"required"`
	AdObjectV2 AdObjectV2 `json:"imp" validate:"required"`
	Test       bool       `json:"test"` // Flag indicating that request is test
	TMax       int64      `json:"tmax"` // Max response time for server before timeout
}

func (r *AuctionV2Request) GetAuctionConfigurationParams() (string, string) {
	return "", r.AdObjectV2.AuctionConfigurationUID
}

func (r *AuctionV2Request) SetAuctionConfigurationParams(id int64, uid string) {
	r.AdObjectV2.AuctionConfigurationUID = uid
}

func (r *AuctionV2Request) ToAuctionRequest() AuctionRequest {
	return AuctionRequest{
		BaseRequest: r.BaseRequest,
		AdType:      r.AdType,
		Adapters:    r.Adapters,
		AdObject: AdObject{
			AuctionID:               r.AdObjectV2.AuctionID,
			AuctionConfigurationUID: r.AdObjectV2.AuctionConfigurationUID,
			PriceFloor:              r.AdObjectV2.PriceFloor,
			Banner:                  r.AdObjectV2.Banner,
			Interstitial:            r.AdObjectV2.Interstitial,
			Rewarded:                r.AdObjectV2.Rewarded,
		},
	}
}

func (r *AuctionV2Request) ToBiddingRequest(roundID string) BiddingRequest {
	return BiddingRequest{
		BaseRequest: r.BaseRequest,
		AdType:      r.AdType,
		Adapters:    r.Adapters,
		Imp:         r.AdObjectV2.ToImp(roundID),
		Test:        r.Test,
		TMax:        r.TMax,
	}
}
