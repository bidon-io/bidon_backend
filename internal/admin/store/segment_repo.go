package adminstore

import (
	"context"
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
	"gorm.io/gorm"
)

type SegmentRepo struct {
	*resourceRepo[admin.Segment, admin.SegmentAttrs, db.Segment]
}

func NewSegmentRepo(d *db.DB) *SegmentRepo {
	return &SegmentRepo{
		resourceRepo: &resourceRepo[admin.Segment, admin.SegmentAttrs, db.Segment]{
			db:           d,
			mapper:       segmentMapper{db: d},
			associations: []string{"App"},
		},
	}
}

func (r *SegmentRepo) ListOwnedByUser(ctx context.Context, userID int64) ([]admin.Segment, error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Select("user_id").Where(map[string]any{"user_id": userID}))
	})
}

func (r *SegmentRepo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.Segment, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Select("user_id").Where(map[string]any{"user_id": userID}))
	})
}

type segmentMapper struct {
	db *db.DB
}

//lint:ignore U1000 this method is used by generic struct
func (m segmentMapper) dbModel(s *admin.SegmentAttrs, id int64) *db.Segment {
	publicUID := sql.NullInt64{}
	if id == 0 {
		publicUID.Int64 = m.db.GenerateSnowflakeID()
		publicUID.Valid = true
	}

	return &db.Segment{
		Model:       db.Model{ID: id},
		Name:        s.Name,
		Description: s.Description,
		Filters:     s.Filters,
		Enabled:     s.Enabled,
		AppID:       s.AppID,
		Priority:    s.Priority,
		PublicUID:   publicUID,
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
