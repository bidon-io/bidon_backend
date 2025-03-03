package adminstore

import (
	"context"
	"strconv"

	"gorm.io/gorm"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/resource"
	"github.com/bidon-io/bidon-backend/internal/db"
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

func (r *SegmentRepo) ListOwnedByUser(ctx context.Context, userID int64, _ map[string][]string) (*resource.Collection[admin.Segment], error) {
	return r.list(ctx, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Table("App").Where(map[string]any{"user_id": userID}))
	}, nil)
}

func (r *SegmentRepo) FindOwnedByUser(ctx context.Context, userID int64, id int64) (*admin.Segment, error) {
	return r.find(ctx, id, func(db *gorm.DB) *gorm.DB {
		s := db.Session(&gorm.Session{NewDB: true})
		return db.InnerJoins("App", s.Table("App").Where(map[string]any{"user_id": userID}))
	})
}

type segmentMapper struct {
	db *db.DB
}

//lint:ignore U1000 this method is used by generic struct
func (m segmentMapper) dbModel(s *admin.SegmentAttrs, id int64) *db.Segment {
	return &db.Segment{
		ID:          id,
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
		PublicUID:    strconv.FormatInt(s.PublicUID.Int64, 10),
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
