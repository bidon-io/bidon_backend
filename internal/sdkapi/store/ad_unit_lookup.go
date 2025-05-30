package store

import (
	"context"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/db"
)

type AdUnitLookup struct {
	DB    *db.DB
	Cache cache[int64]
}

func (a *AdUnitLookup) GetInternalIDByUIDCached(ctx context.Context, uid string) (int64, error) {
	if uid == "" {
		return 0, nil
	}

	cacheKey := []byte("ad_unit_uid:" + uid)

	return a.Cache.Get(ctx, cacheKey, func(ctx context.Context) (int64, error) {
		return a.GetInternalIDByUID(ctx, uid)
	})
}

func (a *AdUnitLookup) GetInternalIDByUID(ctx context.Context, uid string) (int64, error) {
	if uid == "" {
		return 0, nil
	}

	publicUID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return 0, nil
	}

	var lineItem db.LineItem

	err = a.DB.WithContext(ctx).
		Select("id").
		Where("public_uid = ?", publicUID).
		First(&lineItem).Error

	if err != nil {
		return 0, nil
	}

	return lineItem.ID, nil
}
