package db

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.PublicUID == (sql.NullInt64{}) {
		snowflakeID, err := generateSnowflakeID(tx)
		if err != nil {
			return fmt.Errorf("generate snowflake id: %v", err)
		}

		u.PublicUID = sql.NullInt64{
			Int64: snowflakeID,
			Valid: true,
		}
	}

	return nil
}
