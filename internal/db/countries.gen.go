// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"
)

const TableNameCountry = "countries"

// Country mapped from table <countries>
type Country struct {
	ID         int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	Alpha2Code string         `gorm:"column:alpha_2_code;type:character varying;not null;uniqueIndex:index_countries_on_alpha_2_code,priority:1" json:"alpha_2_code"`
	Alpha3Code string         `gorm:"column:alpha_3_code;type:character varying;not null;uniqueIndex:index_countries_on_alpha_3_code,priority:1" json:"alpha_3_code"`
	HumanName  sql.NullString `gorm:"column:human_name;type:character varying" json:"human_name"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp(6) without time zone;not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;type:timestamp(6) without time zone;not null" json:"updated_at"`
}

// TableName Country's table name
func (*Country) TableName() string {
	return TableNameCountry
}
