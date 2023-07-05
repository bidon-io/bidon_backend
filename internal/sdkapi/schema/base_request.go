package schema

import "github.com/bidon-io/bidon-backend/internal/ad"

type BaseRequest struct {
	AdType      ad.Type      `param:"ad_type"`
	Device      Device       `json:"device" validate:"required"`
	Session     Session      `json:"session" validate:"required"`
	App         App          `json:"app" validate:"required"`
	User        User         `json:"user" validate:"required"`
	Geo         *Geo         `json:"geo"`
	Regulations *Regulations `json:"regs"`
	Ext         string       `json:"ext"`
	Token       string       `json:"token"`
	Segment     Segment      `json:"segment"`
}

func (r *BaseRequest) GetApp() App {
	return r.App
}

func (r *BaseRequest) GetGeo() Geo {
	if r.Device.Geo != nil {
		return *r.Device.Geo
	} else if r.Geo != nil {
		return *r.Geo
	}

	return Geo{}
}
