package store

import (
	"context"
	"errors"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"gorm.io/gorm"
)

type AppFetcher struct {
	DB *db.DB
}

func (f *AppFetcher) Fetch(ctx context.Context, appKey, appBundle string) (*sdkapi.App, error) {
	var dbApp db.App
	err := f.DB.
		WithContext(ctx).
		Select("id").
		Find(&dbApp, map[string]any{"app_key": appKey, "package_name": appBundle}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = sdkapi.ErrAppNotValid
		}

		return nil, err
	}

	return &sdkapi.App{ID: dbApp.ID}, nil
}
