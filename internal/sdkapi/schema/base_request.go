package schema

type BaseRequest struct {
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

func (r BaseRequest) Map() map[string]any {
	m := map[string]any{
		"device":  r.Device.Map(),
		"session": r.Session.Map(),
		"app":     r.App.Map(),
		"user":    r.User.Map(),
		"ext":     r.Ext,
		"token":   r.Token,
		"segment": r.Segment.Map(),
	}

	if r.Geo != nil {
		m["geo"] = r.Geo.Map()
	}
	if r.Regulations != nil {
		m["regs"] = r.Regulations.Map()
	}

	return m
}

func (r BaseRequest) GetApp() App {
	return r.App
}

func (r BaseRequest) GetGeo() Geo {
	if r.Device.Geo != nil {
		return *r.Device.Geo
	} else if r.Geo != nil {
		return *r.Geo
	}

	return Geo{}
}

func (r BaseRequest) GetRegulations() Regulations {
	if r.Regulations != nil {
		return *r.Regulations
	}

	return Regulations{}
}
