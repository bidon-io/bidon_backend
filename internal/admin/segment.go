// Package admin implements an HTTP API handlers for managing entities.
package admin

import "github.com/bidon-io/bidon-backend/internal/segment"

type Segment struct {
	ID int64 `json:"id"`
	SegmentAttrs
	App `json:"app"`
}

type SegmentAttrs struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Filters     []segment.Filter `json:"filters"`
	Enabled     *bool            `json:"enabled"`
	AppID       int64            `json:"app_id"`
	Priority    int32            `json:"priority"`
}

type SegmentService = resourceService[Segment, SegmentAttrs]
