package store

import (
	"context"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/db"
)

type AdUnitLookup struct {
	DB    *db.DB
	Cache cache[*db.LineItem]
}

func (a *AdUnitLookup) GetByUIDCached(ctx context.Context, uid string) (*db.LineItem, error) {
	if uid == "" {
		return nil, nil
	}

	cacheKey := []byte("ad_unit_lookup:" + uid)

	return a.Cache.Get(ctx, cacheKey, func(ctx context.Context) (*db.LineItem, error) {
		return a.GetByUID(ctx, uid)
	})
}

func (a *AdUnitLookup) GetByUID(ctx context.Context, uid string) (*db.LineItem, error) {
	if uid == "" {
		return nil, nil
	}

	publicUID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return nil, nil
	}

	var lineItem db.LineItem

	err = a.DB.WithContext(ctx).
		Select("id", "extra").
		Where("public_uid = ?", publicUID).
		First(&lineItem).Error

	if err != nil {
		return nil, nil
	}

	return &lineItem, nil
}
