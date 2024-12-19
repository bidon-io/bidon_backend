package schema

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
)

// AuctionV2Request Has both the auction and bidding request fields
type AuctionV2Request struct {
	BaseRequest
	AdType   ad.Type         `param:"ad_type"`
	Adapters Adapters        `json:"adapters" validate:"required"`
	AdObject AdObjectV2      `json:"ad_object" validate:"required"`
	AdCache  []AdCacheObject `json:"ad_cache,omitempty"`
	Test     bool            `json:"test"` // Flag indicating that request is test
	TMax     int64           `json:"tmax"` // Max response time for server before timeout
}

func (r *AuctionV2Request) GetAuctionConfigurationParams() (string, string) {
	return "", r.AdObject.AuctionConfigurationUID
}

func (r *AuctionV2Request) SetAuctionConfigurationParams(id int64, uid string) {
	r.AdObject.AuctionConfigurationUID = uid
}

func (r *AuctionV2Request) ToBiddingRequest() BiddingRequest {
	return BiddingRequest{
		BaseRequest: r.BaseRequest,
		AdType:      r.AdType,
		Adapters:    r.Adapters,
		Imp:         r.AdObject.ToImp(),
		Test:        r.Test,
		TMax:        r.TMax,
	}
}
