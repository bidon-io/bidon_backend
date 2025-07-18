package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type ConfigFetcher struct {
	DB    *db.DB
	Cache cache[*auction.Config]
}

func (m *ConfigFetcher) Match(ctx context.Context, appID int64, adType ad.Type, segmentID int64, version string) (*auction.Config, error) {
	dbConfig := &db.AuctionConfiguration{}

	query := m.DB.
		WithContext(ctx).
		Select("id", "public_uid", "external_win_notifications", "rounds", "demands", "bidding", "ad_unit_ids", "pricefloor", "timeout").
		Where(map[string]any{
			"app_id":  appID,
			"ad_type": db.AdTypeFromDomain(adType),
		}).
		Order("segment_id, is_default DESC, created_at DESC")

	if segmentID != 0 {
		query = query.Where("segment_id = ? OR segment_id IS NULL", segmentID)
	} else {
		query = query.Where("segment_id IS NULL")
	}

	if version == "v2" {
		query = query.Where("settings->>'v2' = ?", "true")
	} else {
		query = query.Where("settings->>'v2' IS NULL")
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
		Demands:                  db.StringArrayToAdapterKeys(&dbConfig.Demands),
		Bidding:                  db.StringArrayToAdapterKeys(&dbConfig.Bidding),
		AdUnitIDs:                dbConfig.AdUnitIds,
		PriceFloor:               dbConfig.Pricefloor,
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
		Select("id", "public_uid", "external_win_notifications", "rounds", "demands", "bidding", "ad_unit_ids", "pricefloor", "timeout").
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
		Demands:                  db.StringArrayToAdapterKeys(&dbConfig.Demands),
		Bidding:                  db.StringArrayToAdapterKeys(&dbConfig.Bidding),
		AdUnitIDs:                dbConfig.AdUnitIds,
		PriceFloor:               dbConfig.Pricefloor,
		Timeout:                  int(dbConfig.Timeout),
	}

	return config
}

// FetchBidMachinePlacements fetches auction configurations that include BidMachine in demands or bidding
// and returns a map of auction_key to placement_id from line_items
func (m *ConfigFetcher) FetchBidMachinePlacements(ctx context.Context, appID int64) (map[string]string, error) {
	// Get all auction configuration IDs that include BidMachine in demands or bidding
	var configIDs []int64
	err := m.DB.
		WithContext(ctx).
		Model(&db.AuctionConfiguration{}).
		Select("id").
		Where("app_id = ? AND (? = ANY(demands) OR ? = ANY(bidding)) AND auction_key IS NOT NULL AND auction_key != ''", appID, "bidmachine", "bidmachine").
		Pluck("id", &configIDs).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch auction configuration IDs: %w", err)
	}

	if len(configIDs) == 0 {
		return make(map[string]string), nil
	}

	// Get all line items that belong to BidMachine and are used in these auction configurations
	type result struct {
		AuctionKey string `gorm:"column:auction_key"`
		Placement  string `gorm:"column:placement"`
	}

	var results []result
	err = m.DB.
		WithContext(ctx).
		Table("line_items li").
		Select("ac.auction_key, li.extra->>'placement' as placement").
		Joins("JOIN auction_configurations ac ON li.id = ANY(ac.ad_unit_ids)").
		Joins("JOIN demand_source_accounts dsa ON li.account_id = dsa.id").
		Joins("JOIN demand_sources ds ON dsa.demand_source_id = ds.id").
		Where("ac.id IN (?) AND ds.api_key = ? AND li.extra->>'placement' IS NOT NULL AND li.extra->>'placement' != ''", configIDs, "bidmachine").
		Scan(&results).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch bidmachine placements: %w", err)
	}

	// Build the map, using the first placement found for each auction_key
	placements := make(map[string]string)
	for _, result := range results {
		if _, exists := placements[result.AuctionKey]; !exists {
			placements[result.AuctionKey] = result.Placement
		}
	}

	return placements, nil
}
