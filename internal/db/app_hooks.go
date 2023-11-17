package db

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

func (a *App) BeforeCreate(tx *gorm.DB) error {
	if a.PublicUID == (sql.NullInt64{}) {
		snowflakeID, err := generateSnowflakeID(tx)
		if err != nil {
			return fmt.Errorf("generate snowflake id: %v", err)
		}

		a.PublicUID = sql.NullInt64{
			Int64: snowflakeID,
			Valid: true,
		}
	}

	return nil
}
