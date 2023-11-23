package dbtest

import (
	"fmt"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/segment"
)

func segmentDefaults(n uint32) func(*db.Segment) {
	return func(seg *db.Segment) {
		if seg.AppID == 0 && seg.App.ID == 0 {
			seg.App = BuildApp(func(app *db.App) {
				*app = seg.App
			})
		}
		if seg.Name == "" {
			seg.Name = fmt.Sprintf("Test Segment %d", n)
		}
		if seg.Description == "" {
			seg.Description = fmt.Sprintf("Test Segment %d Description", n)
		}
		if seg.Filters == nil {
			seg.Filters = []segment.Filter{}
		}
	}
}

func BuildSegment(opts ...func(*db.Segment)) db.Segment {
	var seg db.Segment

	n := counter.get("segment")

	opts = append(opts, segmentDefaults(n))
	for _, opt := range opts {
		opt(&seg)
	}

	return seg
}

func CreateSegment(t testing.TB, tx *db.DB, opts ...func(*db.Segment)) db.Segment {
	t.Helper()

	seg := BuildSegment(opts...)
	if err := tx.Create(&seg).Error; err != nil {
		t.Fatalf("Failed to create segment: %v", err)
	}

	return seg
}
