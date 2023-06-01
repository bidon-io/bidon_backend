package store

import (
  "context"
  "github.com/bidon-io/bidon-backend/internal/admin"
  "gorm.io/gorm"
  "gorm.io/gorm/clause"
)

type segment struct {
	Model
	Name        string                `gorm:"column:name;type:varchar;not null"`
	Description string                `gorm:"column:description;type:text;not null"`
	Filters     []admin.SegmentFilter `gorm:"column:filters;type:jsonb;not null;default:'[]';serializer:json"`
	Enabled     *bool                 `gorm:"column:enabled;type:bool;not null;default:true"`
	AppID       int64                 `gorm:"column:app_id;type:bigint;not null"`
}

func fromAttrs(attrs *admin.SegmentAttrs) segment {
	return segment{
		Name:        attrs.Name,
		Description: attrs.Description,
		Filters:     attrs.Filters,
		Enabled:     attrs.Enabled,
		AppID:       attrs.AppID,
	}
}

func (s *segment) segment() admin.Segment {
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

type SegmentRepo struct {
	DB *gorm.DB
}

func (r *SegmentRepo) List(ctx context.Context) ([]admin.Segment, error) {
	var dbSegments []segment
	if err := r.DB.WithContext(ctx).Find(&dbSegments).Error; err != nil {
		return nil, err
	}

	segments := make([]admin.Segment, len(dbSegments))
	for i, segmentDbModel := range dbSegments {
		segments[i] = segmentDbModel.segment()
	}

	return segments, nil
}

func (r *SegmentRepo) Find(ctx context.Context, id int64) (*admin.Segment, error) {
	var dbSegment segment
	if err := r.DB.WithContext(ctx).First(&dbSegment, id).Error; err != nil {
		return nil, err
	}

	segment := dbSegment.segment()
	return &segment, nil
}

func (r *SegmentRepo) Create(ctx context.Context, attrs *admin.SegmentAttrs) (*admin.Segment, error) {
	dbSegment := fromAttrs(attrs)
	if err := r.DB.WithContext(ctx).Create(&dbSegment).Error; err != nil {
		return nil, err
	}

	segment := dbSegment.segment()
	return &segment, nil
}

func (r *SegmentRepo) Update(ctx context.Context, id int64, attrs *admin.SegmentAttrs) (*admin.Segment, error) {
	dbSegment := fromAttrs(attrs)
	dbSegment.ID = id

	if err := r.DB.WithContext(ctx).Model(&dbSegment).Clauses(clause.Returning{}).Updates(&dbSegment).Error; err != nil {
		return nil, err
	}

	segment := dbSegment.segment()
	return &segment, nil
}

func (r *SegmentRepo) Delete(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&segment{}, id).Error
}
