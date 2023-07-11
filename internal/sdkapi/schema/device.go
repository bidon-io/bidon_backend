package schema

import "github.com/bidon-io/bidon-backend/internal/device"

type Device struct {
	Geo             *Geo        `json:"geo"`
	UserAgent       string      `json:"ua" validate:"required"`
	Manufacturer    string      `json:"make" validate:"required"`
	Model           string      `json:"model" validate:"required"`
	OS              string      `json:"os" validate:"required"`
	OSVersion       string      `json:"osv" validate:"required"`
	HardwareVersion string      `json:"hwv" validate:"required"`
	Height          int         `json:"h" validate:"required"`
	Width           int         `json:"w" validate:"required"`
	PPI             int         `json:"ppi" validate:"required"`
	PXRatio         float64     `json:"pxratio" validate:"required"`
	JS              *int        `json:"js" validate:"required"`
	Language        string      `json:"language" validate:"required"`
	Carrier         string      `json:"carrier"`
	MCCMNC          string      `json:"mccmnc"`
	ConnectionType  string      `json:"connection_type" validate:"oneof=ETHERNET WIFI CELLULAR CELLULAR_UNKNOWN CELLULAR_2_G CELLULAR_3_G CELLULAR_4_G CELLULAR_5_G"`
	Type            device.Type `json:"type" validate:"oneof=PHONE TABLET"` // TODO: add Marshal/Unmarshal to device.Type
}

func (d Device) Map() map[string]any {
	m := map[string]any{
		"ua":              d.UserAgent,
		"make":            d.Manufacturer,
		"model":           d.Model,
		"os":              d.OS,
		"osv":             d.OSVersion,
		"hwv":             d.HardwareVersion,
		"h":               d.Height,
		"w":               d.Width,
		"ppi":             d.PPI,
		"pxratio":         d.PXRatio,
		"js":              d.JS,
		"language":        d.Language,
		"carrier":         d.Carrier,
		"mccmnc":          d.MCCMNC,
		"connection_type": d.ConnectionType,
		"type":            d.Type,
	}

	if d.Geo != nil {
		m["geo"] = d.Geo.Map()
	}

	return m
}
