package schema

import (
	"strconv"
	"strings"

	"github.com/bidon-io/bidon-backend/internal/ad"
)

type BiddingRequest struct {
	BaseRequest
	AdType   ad.Type  `param:"ad_type"`
	Adapters Adapters `json:"adapters" validate:"required"`
	AdObject AdObject `json:"imp" validate:"required"`
	Test     bool     `json:"test"` // Flag indicating that request is test
	TMax     int64    `json:"tmax"` // Max response time for server before timeout
}

func (b *BiddingRequest) NormalizeValues() {
	b.BaseRequest.NormalizeValues()
	b.AdObject.AuctionID = strings.ToLower(b.AdObject.AuctionID)
}

func (b *BiddingRequest) GetAuctionConfigurationParams() (string, string) {
	return strconv.FormatInt(b.AdObject.AuctionConfigurationID, 10), b.AdObject.AuctionConfigurationUID
}

func (b *BiddingRequest) SetAuctionConfigurationParams(id int64, uid string) {
	b.AdObject.AuctionConfigurationID = id
	b.AdObject.AuctionConfigurationUID = uid
}
