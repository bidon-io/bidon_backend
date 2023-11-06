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

func (r *AuctionRequest) Map() map[string]any {
	m := r.BaseRequest.Map()

	m["ad_type"] = r.AdType
	m["ad_object"] = r.AdObject.Map()

	if r.Adapters != nil {
		m["adapters"] = r.Adapters.Map()
	}

	return m
}

func (r *AuctionRequest) GetAuctionConfigurationParams() (string, string) {
	return "", r.AdObject.AuctionConfigurationUID
}

func (r *AuctionRequest) SetAuctionConfigurationParams(id int64, uid string) {
	r.AdObject.AuctionConfigurationUID = uid
}
