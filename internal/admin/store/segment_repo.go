package store

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type SegmentRepo = resourceRepo[admin.Segment, admin.SegmentAttrs, db.Segment]

func NewSegmentRepo(db *db.DB) *SegmentRepo {
	return &SegmentRepo{
		db:           db,
		mapper:       segmentMapper{},
		associations: []string{"App"},
	}
}

type segmentMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m segmentMapper) dbModel(s *admin.SegmentAttrs, id int64) *db.Segment {
	return &db.Segment{
		Model:       db.Model{ID: id},
		Name:        s.Name,
		Description: s.Description,
		Filters:     s.Filters,
		Enabled:     s.Enabled,
		AppID:       s.AppID,
		Priority:    s.Priority,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m segmentMapper) resource(s *db.Segment) admin.Segment {
	return admin.Segment{
		ID:           s.ID,
		SegmentAttrs: m.resourceAttrs(s),
		App: admin.App{
			ID:       s.App.ID,
			AppAttrs: appMapper{}.resourceAttrs(&s.App),
		},
	}
}

func (m segmentMapper) resourceAttrs(s *db.Segment) admin.SegmentAttrs {
	return admin.SegmentAttrs{
		Name:        s.Name,
		Description: s.Description,
		Filters:     s.Filters,
		Enabled:     s.Enabled,
		AppID:       s.AppID,
		Priority:    s.Priority,
	}
}
