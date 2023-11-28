package segment

import (
	"context"
	"encoding/json"
	"strconv"
)

type Ext struct {
	Gender            string                 `json:"gender"`
	TotalInAppsAmount float64                `json:"total_in_apps_amount"`
	IsPaying          bool                   `json:"is_paying"`
	GameLevel         int                    `json:"game_level"`
	Age               int                    `json:"age"`
	CustomAttributes  map[string]interface{} `json:"custom_attributes"`
}

type Params struct {
	Country string `json:"country"`
	Ext     string `json:"ext"`
	AppID   int64  `json:"app_id"`
}

type Segment struct {
	ID      int64    `json:"id"`
	UID     string   `json:"uid"`
	Filters []Filter `json:"filters"`
}

func (s Segment) StringID() string {
	var segmentID string
	if s.ID != 0 {
		segmentID = strconv.FormatInt(s.ID, 10)
	}

	return segmentID
}

type Filter struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type Matcher struct {
	Fetcher Fetcher
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . Fetcher

type Fetcher interface {
	FetchCached(ctx context.Context, appID int64) ([]Segment, error)
}

func (m *Matcher) Match(ctx context.Context, params *Params) Segment {
	sgmnts, err := m.Fetcher.FetchCached(ctx, params.AppID)
	if err != nil {
		return Segment{ID: 0}
	}

	for _, sgmnt := range sgmnts {
		if isSegmentMatch(sgmnt, params) {
			return sgmnt
		}
	}

	return Segment{ID: 0}
}

func isSegmentMatch(sgmnt Segment, params *Params) bool {
	if len(sgmnt.Filters) == 0 {
		return false
	}

	for _, filter := range sgmnt.Filters {
		switch filter.Type {
		case "country":
			if !matchCountry(filter, params.Country) {
				return false
			}
		case "custom_string":
			if !matchCustomString(filter, params.Ext) {
				return false
			}
		default:
			return false
		}
	}

	return true
}

func matchCountry(filter Filter, country string) bool {
	switch filter.Operator {
	case "IN":
		return containsString(filter.Values, country)
	case "NOT IN":
		return !containsString(filter.Values, country)
	default:
		return false
	}
}

func matchCustomString(filter Filter, ext string) bool {
	var parsedExt Ext

	if err := json.Unmarshal([]byte(ext), &parsedExt); err != nil {
		return false
	}

	switch filter.Operator {
	case "==":
		return filter.Values[0] == parsedExt.CustomAttributes[filter.Name]
	case "!=":
		return filter.Values[0] != parsedExt.CustomAttributes[filter.Name]
	default:
		return false
	}
}

func containsString(values []string, str string) bool {
	for _, v := range values {
		if v == str {
			return true
		}
	}
	return false
}
