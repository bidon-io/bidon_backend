// Package admin implements an HTTP API handlers for managing entities.
package admin

import (
	"context"

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
		repo:   store.Segments(),
		policy: newSegmentPolicy(store),
	}
}

type SegmentRepo interface {
	AllResourceQuerier[Segment]
	OwnedResourceQuerier[Segment]
	ResourceManipulator[Segment, SegmentAttrs]
}

type segmentPolicy struct {
	repo SegmentRepo

	appPolicy *appPolicy
}

func newSegmentPolicy(store Store) *segmentPolicy {
	return &segmentPolicy{
		repo: store.Segments(),

		appPolicy: newAppPolicy(store),
	}
}

func (p *segmentPolicy) getReadScope(authCtx AuthContext) resourceScope[Segment] {
	return &ownedResourceScope[Segment]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *segmentPolicy) getManageScope(authCtx AuthContext) resourceScope[Segment] {
	return &ownedResourceScope[Segment]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *segmentPolicy) authorizeCreate(ctx context.Context, authCtx AuthContext, attrs *SegmentAttrs) error {
	// Check if user can manage the app.
	_, err := p.appPolicy.getManageScope(authCtx).find(ctx, attrs.AppID)
	if err != nil {
		return err
	}

	return nil
}

func (p *segmentPolicy) authorizeUpdate(ctx context.Context, authCtx AuthContext, segment *Segment, attrs *SegmentAttrs) error {
	// If user tries to change the app and app is not the same as before, check if user can manage the new app.
	if attrs.AppID != 0 && attrs.AppID != segment.AppID {
		_, err := p.appPolicy.getManageScope(authCtx).find(ctx, attrs.AppID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *segmentPolicy) authorizeDelete(_ context.Context, _ AuthContext, _ *Segment) error {
	return nil
}
