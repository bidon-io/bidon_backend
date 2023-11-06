package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"gorm.io/gorm"
)

type AppFetcher struct {
	DB    *db.DB
	Cache cache
}

type cache interface {
	Get(context.Context, []byte, func(ctx context.Context) (sdkapi.App, error)) (sdkapi.App, error)
}

func (f *AppFetcher) FetchCached(ctx context.Context, appKey, appBundle string) (app sdkapi.App, err error) {
	cacheKey := fmt.Sprintf("app:%s:%s", appKey, appBundle)

	return f.Cache.Get(ctx, []byte(cacheKey), func(ctx context.Context) (sdkapi.App, error) {
		return f.Fetch(ctx, appKey, appBundle)
	})
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

func (f *AdapterInitConfigsFetcher) FetchAdapterInitConfigs(ctx context.Context, appID int64, adapterKeys []adapter.Key, sdkVersion *semver.Version) ([]sdkapi.AdapterInitConfig, error) {
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

		// TODO: remove this block when we drop support for 0.4.x
		if !sdkapi.Version05GTEConstraint.Check(sdkVersion) {
			amazonConfig, ok := config.(*sdkapi.AmazonInitConfig)
			if ok {
				amazonConfig.Slots, err = f.fetchAmazonSlots(ctx, appID)
				if err != nil {
					return nil, fmt.Errorf("fetch amazon slots: %v", err)
				}
			}
		}

		configs = append(configs, config)
	}

	return configs, nil
}

// Deprecated: amazon slots moved to the auction as of 0.5.0
func (f *AdapterInitConfigsFetcher) fetchAmazonSlots(ctx context.Context, appID int64) ([]sdkapi.AmazonSlot, error) {
	var dbLineItems []db.LineItem

	err := f.DB.
		WithContext(ctx).
		Select("line_items.id, line_items.extra, line_items.ad_type, line_items.format").
		Where("app_id", appID).
		InnerJoins("Account", f.DB.Select("id")).
		InnerJoins("Account.DemandSource", f.DB.Select("api_key").Where("api_key", adapter.AmazonKey)).
		Order("line_items.id").
		Find(&dbLineItems).
		Error

	if err != nil {
		return nil, fmt.Errorf("find line items: %v", err)
	}

	slots := make([]sdkapi.AmazonSlot, 0, len(dbLineItems))
	for _, lineItem := range dbLineItems {
		slot := sdkapi.AmazonSlot{}

		slotUUID, ok := lineItem.Extra["slot_uuid"].(string)
		if !ok {
			return nil, fmt.Errorf("slot_uuid is either missing or not a string")
		}
		slot.SlotUUID = slotUUID

		format, ok := lineItem.Extra["format"].(string)
		if !ok {
			return nil, fmt.Errorf("format is either missing or not a string")
		}
		slot.Format = format

		slots = append(slots, slot)
	}

	return slots, nil
}
