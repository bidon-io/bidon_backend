// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"

	"github.com/bidon-io/bidon-backend/internal/segment"
)

const TableNameSegment = "segments"

// Segment mapped from table <segments>
type Segment struct {
	ID          int64            `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	Name        string           `gorm:"column:name;type:character varying;not null" json:"name"`
	Description string           `gorm:"column:description;type:text;not null" json:"description"`
	Filters     []segment.Filter `gorm:"column:filters;type:jsonb;not null;default:[];serializer:json" json:"filters"`
	Enabled     *bool            `gorm:"column:enabled;type:boolean;not null;default:true" json:"enabled"`
	AppID       int64            `gorm:"column:app_id;type:bigint;not null;index:index_segments_on_app_id,priority:1" json:"app_id"`
	CreatedAt   time.Time        `gorm:"column:created_at;type:timestamp(6) without time zone;not null" json:"created_at"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;type:timestamp(6) without time zone;not null" json:"updated_at"`
	Priority    int32            `gorm:"column:priority;type:integer;not null" json:"priority"`
	PublicUID   sql.NullInt64    `gorm:"column:public_uid;type:bigint;uniqueIndex:index_segments_on_public_uid,priority:1" json:"public_uid"`
	App         App              `json:"app"`
}

// TableName Segment's table name
func (*Segment) TableName() string {
	return TableNameSegment
}
