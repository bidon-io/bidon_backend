package store

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
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
		return nil, err
	}

	return &sdkapi.App{ID: dbApp.ID}, nil
}
