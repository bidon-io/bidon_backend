package schema

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
)

type AuctionRequest struct {
	BaseRequest
	AdType   ad.Type  `param:"ad_type"`
	Adapters Adapters `json:"adapters" validate:"required"`
	AdObject AdObject `json:"ad_object" validate:"required"`
}

func (r *AuctionRequest) GetAuctionConfigurationParams() (string, string) {
	return "", r.AdObject.AuctionConfigurationUID
}

func (r *AuctionRequest) SetAuctionConfigurationParams(id int64, uid string) {
	r.AdObject.AuctionConfigurationUID = uid
}
