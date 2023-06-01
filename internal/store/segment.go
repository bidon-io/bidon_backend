package store

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
)

type segment struct {
	Model
	Name        string                `gorm:"column:name;type:varchar;not null"`
	Description string                `gorm:"column:description;type:text;not null"`
	Filters     []admin.SegmentFilter `gorm:"column:filters;type:jsonb;not null;default:'[]';serializer:json"`
	Enabled     *bool                 `gorm:"column:enabled;type:bool;not null;default:true"`
	AppID       int64                 `gorm:"column:app_id;type:bigint;not null"`
}

//lint:ignore U1000 this method is used by generic struct
func (s *segment) initFromResourceAttrs(attrs *admin.SegmentAttrs) {
	s.Name = attrs.Name
	s.Description = attrs.Description
	s.Filters = attrs.Filters
	s.Enabled = attrs.Enabled
	s.AppID = attrs.AppID
}

//lint:ignore U1000 this method is used by generic struct
func (s *segment) toResource() admin.Segment {
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

type SegmentRepo = resourceRepo[admin.Segment, admin.SegmentAttrs, segment, *segment]
