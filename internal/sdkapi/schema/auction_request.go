package schema

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
)

// AuctionRequest Has both the auction and bidding request fields
type AuctionRequest struct {
	BaseRequest
	AdType   ad.Type         `param:"ad_type"`
	Adapters Adapters        `json:"adapters" validate:"required"`
	AdObject AdObject        `json:"ad_object" validate:"required"`
	AdCache  []AdCacheObject `json:"ad_cache,omitempty"`
	Test     bool            `json:"test"` // Flag indicating that request is test
	TMax     int64           `json:"tmax"` // Max response time for server before timeout
}

func (r *AuctionRequest) GetAuctionConfigurationParams() (string, string) {
	return "", r.AdObject.AuctionConfigurationUID
}

func (r *AuctionRequest) SetAuctionConfigurationParams(_ int64, uid string) {
	r.AdObject.AuctionConfigurationUID = uid
}
