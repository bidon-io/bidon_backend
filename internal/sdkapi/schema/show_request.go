package schema

import "github.com/bidon-io/bidon-backend/internal/ad"

type ShowRequest struct {
	BaseRequest
	AdType ad.Type `param:"ad_type"`
	Bid    *Bid    `json:"bid" validate:"required_without=Show"`
	Show   *Bid    `json:"show" validate:"required_without=Bid"`
}

func (r *ShowRequest) Map() map[string]any {
	m := r.BaseRequest.Map()

	m["ad_type"] = r.AdType

	if r.Bid != nil {
		m["show"] = r.Bid.Map()
	}
	if r.Show != nil {
		m["bid"] = r.Show.Map()
	}

	return m
}
