package segment

import (
	"encoding/json"
	"github.com/bidon-io/bidon-backend/internal/admin"
)

type Params struct {
	Country     string `json:"country"`
	CustomProps string `json:"custom_props"`
}

func Match(segments []admin.Segment, params Params) *admin.Segment {
	for _, segment := range segments {
		isMatched := false

		for _, filter := range segment.Filters {
			switch filter.Type {
			case "country":
				isMatched = matchCountry(filter, params.Country)
			case "custom_string":
				isMatched = matchCustomString(filter, params.CustomProps)
			}

			if isMatched {
				return &segment
			}
		}
	}

	return nil
}

func matchCountry(filter admin.SegmentFilter, country string) bool {
	switch filter.Operator {
	case "IN":
		return containsString(filter.Values, country)
	case "NOT IN":
		return !containsString(filter.Values, country)
	default:
		return false
	}
}

func matchCustomString(filter admin.SegmentFilter, customProps string) bool {
	var parsedProps map[string]string

	if err := json.Unmarshal([]byte(customProps), &parsedProps); err != nil {
		return false
	}

	switch filter.Operator {
	case "==":
		return filter.Values[0] == parsedProps[filter.Name]
	case "!=":
		return filter.Values[0] != parsedProps[filter.Name]
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
