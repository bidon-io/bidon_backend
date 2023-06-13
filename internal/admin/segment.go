// Package admin implements an HTTP API handlers for managing entities.
package admin

type Segment struct {
	ID int64 `json:"id"`
	SegmentAttrs
}

type SegmentAttrs struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Filters     []SegmentFilter `json:"filters"`
	Enabled     *bool           `json:"enabled"`
	AppID       int64           `json:"app_id"`
	Priority    int32           `json:"priority"`
}

type SegmentFilter struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type SegmentService = resourceService[Segment, SegmentAttrs]
