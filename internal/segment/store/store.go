package store

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/segment"
)

type SegmentFetcher struct {
	DB *db.DB
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
			UID:     dbSegment.PublicUID.Int64,
			Filters: dbSegment.Filters,
		})
	}

	return sgmnts, nil
}
