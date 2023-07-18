package store

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	builder "github.com/bidon-io/bidon-backend/internal/bidding/adapters_builder"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type AppDemandProfileFetcher struct {
	DB *db.DB
}

func (f *AppDemandProfileFetcher) Fetch(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]builder.AppDemandProfile, error) {
	var dbProfiles []db.AppDemandProfile

	err := f.DB.
		WithContext(ctx).
		Select("app_demand_profiles.id").
		Where("app_id", appID).
		InnerJoins("Account", f.DB.Select("id", "extra")).
		InnerJoins("Account.DemandSource", f.DB.Select("api_key").Where(map[string]any{"api_key": adapterKeys})).
		Find(&dbProfiles).
		Error
	if err != nil {
		return nil, err
	}

	profiles := make([]builder.AppDemandProfile, len(dbProfiles))
	for i, dbProfile := range dbProfiles {
		profile := &profiles[i]

		profile.AdapterKey = adapter.Key(dbProfile.Account.DemandSource.APIKey)
		profile.AccountExtra = dbProfile.Account.Extra
	}

	return profiles, nil
}
