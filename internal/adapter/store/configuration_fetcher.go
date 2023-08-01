package store

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type ConfigurationFetcher struct {
	DB *db.DB
}

func (f *ConfigurationFetcher) Fetch(ctx context.Context, appID int64, adapterKeys []adapter.Key) (adapter.RawConfigsMap, error) {
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
		return nil, err
	}

	configs := adapter.RawConfigsMap{}
	for _, dbProfile := range dbProfiles {
		key := adapter.Key(dbProfile.Account.DemandSource.APIKey)
		configs[key] = adapter.Config{
			AccountExtra: dbProfile.Account.Extra,
			AppData:      dbProfile.Data,
		}
	}

	return configs, nil
}
