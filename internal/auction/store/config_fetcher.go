package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/db"
	"gorm.io/gorm"
)

type ConfigFetcher struct {
	DB    *db.DB
	Cache cache[*auction.Config]
}

func (m *ConfigFetcher) Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64) (*auction.Config, error) {
	dbConfig := &db.AuctionConfiguration{}

	query := m.DB.
		WithContext(ctx).
		Select("id", "public_uid", "external_win_notifications", "rounds", "demands", "bidding", "ad_unit_ids", "timeout").
		Where(map[string]any{
			"app_id":  appID,
			"ad_type": db.AdTypeFromDomain(adType),
		}).
		Order("created_at DESC")

	if segmentID != 0 {
		query = query.Where("segment_id = ? OR segment_id IS NULL", segmentID)
	} else {
		query = query.Where("segment_id IS NULL")
	}

	err := query.Take(dbConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = auction.ErrNoAdsFound
		}
		return nil, err
	}

	config := &auction.Config{
		ID:                       dbConfig.ID,
		UID:                      strconv.FormatInt(dbConfig.PublicUID.Int64, 10),
		ExternalWinNotifications: *dbConfig.ExternalWinNotifications,
		Rounds:                   dbConfig.Rounds,
		Demands:                  db.StringArrayToAdapterKeys(&dbConfig.Demands),
		Bidding:                  db.StringArrayToAdapterKeys(&dbConfig.Bidding),
		AdUnitIDs:                dbConfig.AdUnitIds,
		Timeout:                  int(dbConfig.Timeout),
	}

	return config, nil
}

func (m *ConfigFetcher) FetchByUIDCached(ctx context.Context, appID int64, id, uid string) *auction.Config {
	if id == "" && uid == "" {
		return nil
	}

	key := fmt.Sprintf("%d:%s:%s", appID, id, uid)

	res, err := m.Cache.Get(ctx, []byte(key), func(ctx context.Context) (*auction.Config, error) {
		auc := m.FetchByUID(ctx, appID, id, uid)
		if auc == nil {
			return nil, fmt.Errorf("no config found for app_id:%d, id: %s, uid: %s", appID, id, uid)
		} else {
			return auc, nil
		}
	})

	if err != nil {
		return nil
	} else {
		return res
	}
}

// FetchByUID fetches an auction configuration by its public UID or ID
// If both id and uid are empty, returns nil
// If both id and uid are provided, uid takes precedence
// If no configuration is found, returns nil
func (m *ConfigFetcher) FetchByUID(ctx context.Context, appID int64, id, uid string) *auction.Config {
	if id == "" && uid == "" {
		return nil
	}

	dbConfig := &db.AuctionConfiguration{}

	filter := map[string]any{
		"app_id": appID,
	}
	if uid != "" {
		filter["public_uid"] = uid
	} else {
		filter["id"] = id
	}

	err := m.DB.
		WithContext(ctx).
		Select("id", "public_uid", "external_win_notifications", "rounds", "demands", "bidding", "ad_unit_ids", "timeout").
		Where(filter).
		Order("created_at DESC").
		Take(dbConfig).
		Error
	if err != nil {
		return nil
	}

	config := &auction.Config{
		ID:                       dbConfig.ID,
		UID:                      strconv.FormatInt(dbConfig.PublicUID.Int64, 10),
		ExternalWinNotifications: *dbConfig.ExternalWinNotifications,
		Rounds:                   dbConfig.Rounds,
		Demands:                  db.StringArrayToAdapterKeys(&dbConfig.Demands),
		Bidding:                  db.StringArrayToAdapterKeys(&dbConfig.Bidding),
		AdUnitIDs:                dbConfig.AdUnitIds,
		Timeout:                  int(dbConfig.Timeout),
	}

	return config
}
