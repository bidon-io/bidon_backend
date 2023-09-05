// Package admin implements an HTTP API handlers for managing entities.
package admin

import (
	"github.com/bidon-io/bidon-backend/internal/segment"
)

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

type SegmentService = ResourceService[Segment, SegmentAttrs]

func NewSegmentService(store Store) *SegmentService {
	return &SegmentService{
		repo: store.Segments(),
		policy: &segmentPolicy{
			repo: store.Segments(),
		},
	}
}

type SegmentRepo interface {
	AllResourceQuerier[Segment]
	OwnedResourceQuerier[Segment]
	ResourceManipulator[Segment, SegmentAttrs]
}

type segmentPolicy struct {
	repo SegmentRepo
}

func (p *segmentPolicy) scope(authCtx AuthContext) resourceScope[Segment] {
	return &ownedResourceScope[Segment]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}
