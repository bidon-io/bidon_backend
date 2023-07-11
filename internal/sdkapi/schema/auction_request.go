package schema

import "github.com/bidon-io/bidon-backend/internal/ad"

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

	for key, adapter := range r.Adapters {
		m[string(key)] = adapter.Map()
	}

	return m
}
