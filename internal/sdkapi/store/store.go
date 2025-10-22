package store

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"gorm.io/gorm"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
)

type AppFetcher struct {
	DB    *db.DB
	Cache cache[sdkapi.App]
}

type cache[T any] interface {
	Get(context.Context, []byte, func(ctx context.Context) (T, error)) (T, error)
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
		Select("id", "store_id", "store_url", "categories", "badv", "bcat", "bapp").
		Take(&dbApp, map[string]any{"app_key": appKey, "package_name": appBundle}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app, sdkapi.ErrAppNotValid
		}

		return app, fmt.Errorf("fetch app: %v", err)
	}

	app.ID = dbApp.ID
	app.StoreID = dbApp.StoreID.String
	app.StoreURL = dbApp.StoreURL.String
	app.Categories = dbApp.Categories
	app.Badv = dbApp.Badv.String
	app.Bcat = dbApp.Bcat.String
	app.Bapp = dbApp.Bapp.String

	return app, nil
}

type AdapterInitConfigsFetcher struct {
	DB               *db.DB
	ProfilesCache    cache[[]db.AppDemandProfile]
	AmazonSlotsCache cache[[]sdkapi.AmazonSlot]
	LineItemsCache   cache[[]db.LineItem]
}

func (f *AdapterInitConfigsFetcher) FetchAdapterInitConfigs(ctx context.Context, appID int64, adapterKeys []adapter.Key, setAmazonSlots bool, setOrder bool) ([]sdkapi.AdapterInitConfig, error) {
	dbProfiles, err := f.fetchEnabledAppDemandProfilesCached(ctx, appID, adapterKeys)
	if err != nil {
		return nil, fmt.Errorf("fetch profiles from cache or DB: %w", err)
	}

	configs := make([]sdkapi.AdapterInitConfig, 0, len(dbProfiles))
	for _, profile := range dbProfiles {
		adapterKey := adapter.Key(profile.Account.DemandSource.APIKey)

		config, err := sdkapi.NewAdapterInitConfig(adapterKey, setOrder)
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

		// Set AppKey to SDKKey for AppLovin (SDK API compatibility)
		if adapterKey == adapter.ApplovinKey {
			applovinConfig, ok := config.(*sdkapi.ApplovinInitConfig)
			if ok {
				applovinConfig.AppKey = applovinConfig.SDKKey
			}
		}

		if setAmazonSlots {
			amazonConfig, ok := config.(*sdkapi.AmazonInitConfig)
			if ok {
				amazonConfig.Slots, err = f.fetchAmazonSlotsCached(ctx, appID)
				if err != nil {
					return nil, fmt.Errorf("fetch amazon slots: %v", err)
				}
			}
		}

		if adapterKey == adapter.TaurusXKey {
			taurusxConfig, ok := config.(*sdkapi.TaurusXInitConfig)
			if ok {
				taurusxConfig.PlacementIDs, err = f.fetchTaurusXPlacements(ctx, appID)
				if err != nil {
					return nil, fmt.Errorf("fetch taurusx placements: %v", err)
				}
			}
		}

		configs = append(configs, config)
	}

	return configs, nil
}

func (f *AdapterInitConfigsFetcher) FetchEnabledAdapterKeys(ctx context.Context, appID int64, keys []adapter.Key) ([]adapter.Key, error) {
	profiles, err := f.fetchEnabledAppDemandProfilesCached(ctx, appID, keys)
	if err != nil {
		return nil, err
	}

	enabledKeys := make([]adapter.Key, 0, len(profiles))
	for _, profile := range profiles {
		enabledKeys = append(enabledKeys, adapter.Key(profile.Account.DemandSource.APIKey))
	}
	return enabledKeys, nil
}

func (f *AdapterInitConfigsFetcher) fetchEnabledAppDemandProfilesCached(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]db.AppDemandProfile, error) {
	cacheKey, err := f.profilesCacheKey(appID, adapterKeys)
	if err != nil {
		return nil, fmt.Errorf("generate profiles cache key: %w", err)
	}

	return f.ProfilesCache.Get(ctx, cacheKey, func(ctx context.Context) ([]db.AppDemandProfile, error) {
		return f.fetchEnabledAppDemandProfiles(ctx, appID, adapterKeys)
	})
}

func (f *AdapterInitConfigsFetcher) fetchEnabledAppDemandProfiles(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]db.AppDemandProfile, error) {
	var profiles []db.AppDemandProfile
	err := f.DB.
		WithContext(ctx).
		Select("app_demand_profiles.id, app_demand_profiles.data").
		Where("app_id = ? AND app_demand_profiles.enabled = ?", appID, true).
		InnerJoins("Account", f.DB.Select("id", "extra")).
		InnerJoins("Account.DemandSource", f.DB.Select("api_key").Where(map[string]any{"api_key": adapterKeys})).
		Find(&profiles).
		Error
	if err != nil {
		return nil, fmt.Errorf("find app demand profiles: %v", err)
	}
	return profiles, nil
}

func (f *AdapterInitConfigsFetcher) fetchAmazonSlotsCached(ctx context.Context, appID int64) ([]sdkapi.AmazonSlot, error) {
	cacheKey := f.amazonSlotsCacheKey(appID)

	return f.AmazonSlotsCache.Get(ctx, cacheKey, func(ctx context.Context) ([]sdkapi.AmazonSlot, error) {
		return f.fetchAmazonSlots(ctx, appID)
	})
}

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

func (f *AdapterInitConfigsFetcher) fetchTaurusXPlacements(ctx context.Context, appID int64) ([]sdkapi.TaurusXPlacement, error) {
	lineItems, err := f.fetchLineItemsCached(ctx, appID, adapter.TaurusXKey)
	if err != nil {
		return nil, fmt.Errorf("fetch line items: %v", err)
	}

	if len(lineItems) == 0 {
		return nil, nil
	}

	placements := make([]sdkapi.TaurusXPlacement, 0, len(lineItems))
	for _, lineItem := range lineItems {
		placementID, ok := lineItem.Extra["placement_id"].(string)
		if !ok {
			return nil, fmt.Errorf("placement_id is either missing or not a string")
		}

		format := f.convertAdTypeToFormat(lineItem)

		placements = append(placements, sdkapi.TaurusXPlacement{
			PlacementID: placementID,
			Format:      format,
		})
	}

	return placements, nil
}

func (f *AdapterInitConfigsFetcher) convertAdTypeToFormat(lineItem db.LineItem) string {
	domainType := lineItem.AdType.Domain()
	switch domainType {
	case ad.BannerType:
		if lineItem.Format.Valid && lineItem.Format.String == "MREC" {
			return "MREC"
		}
		return "BANNER"
	case ad.InterstitialType:
		return "INTERSTITIAL"
	case ad.RewardedType:
		return "REWARDED"
	default:
		return "BANNER"
	}
}

func (f *AdapterInitConfigsFetcher) fetchLineItemsCached(ctx context.Context, appID int64, adapterKey adapter.Key) ([]db.LineItem, error) {
	cacheKey := f.lineItemsCacheKey(appID, adapterKey)

	return f.LineItemsCache.Get(ctx, cacheKey, func(ctx context.Context) ([]db.LineItem, error) {
		return f.fetchLineItems(ctx, appID, adapterKey)
	})
}

func (f *AdapterInitConfigsFetcher) fetchLineItems(ctx context.Context, appID int64, adapterKey adapter.Key) ([]db.LineItem, error) {
	var lineItems []db.LineItem

	err := f.DB.
		WithContext(ctx).
		Select("line_items.id, line_items.public_uid, line_items.extra, line_items.bid_floor, line_items.ad_type, line_items.format").
		Where("app_id", appID).
		InnerJoins("Account", f.DB.Select("id")).
		InnerJoins("Account.DemandSource", f.DB.Select("api_key").Where("api_key", adapterKey)).
		Order("line_items.id").
		Find(&lineItems).
		Error

	if err != nil {
		return nil, fmt.Errorf("find line items: %v", err)
	}

	return lineItems, nil
}

func (f *AdapterInitConfigsFetcher) profilesCacheKey(appID int64, adapterKeys []adapter.Key) ([]byte, error) {
	// Sort adapter keys to get deterministic cache key
	sort.Slice(adapterKeys, func(i, j int) bool {
		return adapterKeys[i] < adapterKeys[j]
	})
	cacheKeyData := struct {
		AppID       int64         `json:"app_id"`
		AdapterKeys []adapter.Key `json:"adapter_keys"`
	}{
		AppID:       appID,
		AdapterKeys: adapterKeys,
	}
	jsonData, err := json.Marshal(cacheKeyData)

	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(jsonData)
	return hash[:], nil
}

func (f *AdapterInitConfigsFetcher) amazonSlotsCacheKey(appID int64) []byte {
	return []byte(fmt.Sprintf("amazon_slots:%d", appID))
}

func (f *AdapterInitConfigsFetcher) lineItemsCacheKey(appID int64, adapterKey adapter.Key) []byte {
	return []byte(fmt.Sprintf("line_items:%d:%s", appID, adapterKey))
}
