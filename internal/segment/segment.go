package segment

import (
	"context"
	"encoding/json"
)

type Ext struct {
	Gender            string                 `json:"gender"`
	TotalInAppsAmount int                    `json:"total_in_apps_amount"`
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
	Filters []Filter `json:"filters"`
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

type Fetcher interface {
	Fetch(ctx context.Context, appID int64) ([]Segment, error)
}

func (m *Matcher) Match(ctx context.Context, params *Params) Segment {
	dbSgmnts, err := m.Fetcher.Fetch(ctx, params.AppID)
	if err != nil {
		return Segment{ID: 0}
	}

	for _, dbSgmnt := range dbSgmnts {
		isMatched := false

		for _, filter := range dbSgmnt.Filters {
			switch filter.Type {
			case "country":
				isMatched = matchCountry(filter, params.Country)
			case "custom_string":
				isMatched = matchCustomString(filter, params.Ext)
			}

			if isMatched {
				return Segment{ID: dbSgmnt.ID}
			}
		}
	}

	return Segment{ID: 0}
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
