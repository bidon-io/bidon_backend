package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"gorm.io/gorm"
)

type AppFetcher struct {
	DB *db.DB
}

func (f *AppFetcher) Fetch(ctx context.Context, appKey, appBundle string) (app sdkapi.App, err error) {
	var dbApp db.App
	err = f.DB.
		WithContext(ctx).
		Select("id").
		Take(&dbApp, map[string]any{"app_key": appKey, "package_name": appBundle}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app, sdkapi.ErrAppNotValid
		}

		return app, fmt.Errorf("fetch app: %v", err)
	}

	app.ID = dbApp.ID

	return app, nil
}

type AdapterInitConfigsFetcher struct {
	DB *db.DB
}

func (f *AdapterInitConfigsFetcher) FetchAdapterInitConfigs(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]sdkapi.AdapterInitConfig, error) {
	var dbProfiles []db.AppDemandProfile

	err := f.DB.
		WithContext(ctx).
		Select("app_demand_profiles.id, app_demand_profiles.data").
		Where("app_id", appID).
		InnerJoins("Account", f.DB.Select("id", "extra")).
		InnerJoins("Account.DemandSource", f.DB.Select("api_key").Where(map[string]any{"api_key": adapterKeys})).
		Find(&dbProfiles).
		Error
	if err != nil {
		return nil, fmt.Errorf("find app demand profiles: %v", err)
	}

	configs := make([]sdkapi.AdapterInitConfig, 0, len(dbProfiles))
	for _, profile := range dbProfiles {
		adapterKey := adapter.Key(profile.Account.DemandSource.APIKey)
		config, err := sdkapi.NewAdapterInitConfig(adapterKey)
		if err != nil {
			return nil, fmt.Errorf("new AdapterInitConfig: %w", err)
		}

		err = json.Unmarshal(profile.Account.Extra, config)
		if err != nil {
			return nil, fmt.Errorf("unmarshal account extra: %v", err)
		}

		err = json.Unmarshal(profile.Data, config)
		if err != nil {
			return nil, fmt.Errorf("unmarshal profile data: %v", err)
		}

		applovinConfig, ok := config.(*sdkapi.ApplovinInitConfig)
		if ok {
			applovinConfig.AppKey = applovinConfig.SDKKey
		}

		configs = append(configs, config)
	}

	return configs, nil
}
