package schema

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
)

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

	// Cache for parsed Ext data
	extData      map[string]any
	parseExtOnce sync.Once
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

func (r *BaseRequest) GetRegulations() Regulations {
	if r.Regulations != nil {
		return *r.Regulations
	}

	return Regulations{}
}

func (r *BaseRequest) NormalizeValues() {
	r.User.IDFA = strings.ToLower(r.User.IDFA)
	r.User.IDFV = strings.ToLower(r.User.IDFV)
	r.User.IDG = strings.ToLower(r.User.IDG)
	r.Session.ID = strings.ToLower(r.Session.ID)
}

func (r *BaseRequest) SetSDKVersion(version string) {
	r.App.SDKVersion = version
}

func (r *BaseRequest) GetSDKVersionSemver() (*semver.Version, error) {
	return semver.NewVersion(r.App.SDKVersion)
}

func (r *BaseRequest) GetAuctionConfigurationParams() (string, string) {
	return "", ""
}

func (r *BaseRequest) SetAuctionConfigurationParams(id int64, uid string) {
}

func (r *BaseRequest) GetExtData() map[string]any {
	r.parseExt()
	if r.extData == nil {
		return map[string]any{}
	}

	return r.extData
}

func (r *BaseRequest) GetMediationMode() string {
	ext := r.GetExtData()
	if mode, ok := ext["mediation_mode"].(string); ok {
		return mode
	}

	return ""
}

func (r *BaseRequest) parseExt() {
	r.parseExtOnce.Do(func() {
		if r.Ext == "" {
			return
		}
		_ = json.Unmarshal([]byte(r.Ext), &r.extData)
	})
}
