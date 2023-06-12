package store

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type SegmentRepo = resourceRepo[admin.Segment, admin.SegmentAttrs, db.Segment]

func NewSegmentRepo(db *db.DB) *SegmentRepo {
	return &SegmentRepo{
		db:     db,
		mapper: segmentMapper{},
	}
}

type segmentMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m segmentMapper) dbModel(s *admin.SegmentAttrs) *db.Segment {
	return &db.Segment{
		Name:        s.Name,
		Description: s.Description,
		Filters:     s.Filters,
		Enabled:     s.Enabled,
		AppID:       s.AppID,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m segmentMapper) resource(s *db.Segment) admin.Segment {
	return admin.Segment{
		ID: s.ID,
		SegmentAttrs: admin.SegmentAttrs{
			Name:        s.Name,
			Description: s.Description,
			Filters:     s.Filters,
			Enabled:     s.Enabled,
			AppID:       s.AppID,
		},
	}
}
