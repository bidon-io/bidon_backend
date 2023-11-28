package store

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"strconv"
)

type SegmentFetcher struct {
	DB    *db.DB
	Cache cache
}

type cache interface {
	Get(context.Context, []byte, func(ctx context.Context) ([]segment.Segment, error)) ([]segment.Segment, error)
}

func (f *SegmentFetcher) FetchCached(ctx context.Context, appID int64) ([]segment.Segment, error) {
	return f.Cache.Get(ctx, []byte(strconv.FormatInt(appID, 10)), func(ctx context.Context) ([]segment.Segment, error) {
		return f.Fetch(ctx, appID)
	})
}

func (f *SegmentFetcher) Fetch(ctx context.Context, appID int64) ([]segment.Segment, error) {
	var dbSegments []db.Segment
	var sgmnts []segment.Segment

	err := f.DB.WithContext(ctx).
		Where("app_id = ? AND enabled", appID).
		Order("priority ASC").
		Find(&dbSegments).Error
	if err != nil {
		return nil, err
	}

	for _, dbSegment := range dbSegments {
		sgmnts = append(sgmnts, segment.Segment{
			ID:      dbSegment.ID,
			UID:     strconv.FormatInt(dbSegment.PublicUID.Int64, 10),
			Filters: dbSegment.Filters,
		})
	}

	return sgmnts, nil
}
