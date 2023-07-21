package schema

import (
	"strings"

	"github.com/bidon-io/bidon-backend/internal/ad"
)

type BiddingRequest struct {
	BaseRequest
	AdType   ad.Type  `param:"ad_type"`
	Adapters Adapters `json:"adapters" validate:"required"`
	Imp      Imp      `json:"imp" validate:"required"`
	Test     bool     `json:"test"` // Flag indicating that request is test
	TMax     int      `json:"tmax"` // Max response time for server before timeout
}

func (b *BiddingRequest) NormalizeValues() {
	b.BaseRequest.NormalizeValues()
	b.Imp.ID = strings.ToLower(b.Imp.ID)
	b.Imp.AuctionID = strings.ToLower(b.Imp.AuctionID)
}
