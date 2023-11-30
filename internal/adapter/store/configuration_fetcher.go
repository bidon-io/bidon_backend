package store

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
	"sort"
	"strings"
)

type ConfigurationFetcher struct {
	DB    *db.DB
	Cache cache
}

type cache interface {
	Get(context.Context, []byte, func(ctx context.Context) (adapter.RawConfigsMap, error)) (adapter.RawConfigsMap, error)
}

func (f *ConfigurationFetcher) FetchCached(ctx context.Context, appID int64, adapterKeys []adapter.Key) (adapter.RawConfigsMap, error) {
	key := f.cacheKey(appID, adapterKeys)
	return f.Cache.Get(ctx, key, func(ctx context.Context) (adapter.RawConfigsMap, error) {
		return f.Fetch(ctx, appID, adapterKeys)
	})
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
		return nil, fmt.Errorf("cannot load adapter config from DB: %w", err)
	}

	configs := adapter.RawConfigsMap{}
	for _, dbProfile := range dbProfiles {
		var extra map[string]any
		err = json.Unmarshal(dbProfile.Account.Extra, &extra)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal account extra: %v", err)
		}

		var data map[string]any
		err = json.Unmarshal(dbProfile.Data, &data)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal profile data: %v", err)
		}

		key := adapter.Key(dbProfile.Account.DemandSource.APIKey)
		configs[key] = adapter.Config{
			AccountExtra: extra,
			AppData:      data,
		}
	}

	return configs, nil
}

func (f *ConfigurationFetcher) cacheKey(appID int64, adapterKeys []adapter.Key) []byte {
	stringKeys := make([]string, len(adapterKeys))
	for i, key := range adapterKeys {
		stringKeys[i] = string(key)
	}
	// Sort adapter keys to get deterministic cache key
	sort.Strings(stringKeys)
	key := fmt.Sprintf("%d:%s", appID, strings.Join(stringKeys, ":"))
	hash := sha256.Sum256([]byte(key))

	return hash[:]
}
