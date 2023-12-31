// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package db

import (
	"time"
)

const TableNameDemandSource = "demand_sources"

// DemandSource mapped from table <demand_sources>
type DemandSource struct {
	ID        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	APIKey    string    `gorm:"column:api_key;type:character varying;not null;uniqueIndex:index_demand_sources_on_api_key,priority:1" json:"api_key"`
	HumanName string    `gorm:"column:human_name;type:character varying;not null" json:"human_name"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp(6) without time zone;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp(6) without time zone;not null" json:"updated_at"`
}

// TableName DemandSource's table name
func (*DemandSource) TableName() string {
	return TableNameDemandSource
}
