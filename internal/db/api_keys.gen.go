// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package db

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const TableNameAPIKey = "api_keys"

// APIKey mapped from table <api_keys>
type APIKey struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Value          string    `gorm:"column:value;type:character varying;not null" json:"value"`
	UserID         int64     `gorm:"column:user_id;type:bigint;not null;index:api_keys_user_id_idx,priority:1" json:"user_id"`
	LastAccessedAt time.Time `gorm:"column:last_accessed_at;type:timestamp without time zone" json:"last_accessed_at"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp without time zone;not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp without time zone;not null" json:"updated_at"`
	User           User      `json:"user"`
}

// TableName APIKey's table name
func (*APIKey) TableName() string {
	return TableNameAPIKey
}
