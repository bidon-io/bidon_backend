package store

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type ConfigMatcher struct {
	DB *db.DB
}

func (m *ConfigMatcher) Match(ctx context.Context, appID int64, adType ad.Type) (*auction.Config, error) {
	dbConfig := &db.AuctionConfiguration{}
	err := m.DB.
		WithContext(ctx).
		Select("id", "rounds").
		Where(map[string]any{
			"app_id":  appID,
			"ad_type": db.AdTypeFromDomain(adType),
		}).
		Order("created_at DESC").
		Take(dbConfig).
		Error
	if err != nil {
		return nil, err
	}

	config := &auction.Config{
		ID:     dbConfig.ID,
		Rounds: dbConfig.Rounds,
	}

	return config, nil
}
