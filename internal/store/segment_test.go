package store_test

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/store"
	"github.com/google/go-cmp/cmp"
)

func TestSegmentRepo_List(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.SegmentRepo{DB: tx}

	segments := []admin.SegmentAttrs{
		{
			Name:        "Country Segment",
			Description: "Desc",
			AppID:       1,
			Filters:     []admin.SegmentFilter{{Type: "country", Name: "country", Operator: "in", Values: []string{"US", "UK"}}},
			Enabled:     ptr(true),
		},
		{
			Name:        "Custom String Segment",
			Description: "Desc",
			AppID:       1,
			Filters:     []admin.SegmentFilter{{Type: "string", Name: "custom_str", Operator: "==", Values: []string{"super"}}},
			Enabled:     ptr(false),
		},
		{
			Name:        "Custom Num Segment",
			Description: "Desc",
			AppID:       1,
			Filters:     []admin.SegmentFilter{{Type: "float", Name: "custom_num", Operator: ">=", Values: []string{"33"}}},
			Enabled:     ptr(true),
		},
	}

	want := make([]admin.Segment, len(segments))
	for i, attrs := range segments {
		segment, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, segment, nil)
		}

		want[i] = *segment
	}

	got, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestSegmentRepo_Find(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.SegmentRepo{DB: tx}

	attrs := &admin.SegmentAttrs{
		Name:        "Country Segment",
		Description: "Desc",
		AppID:       1,
		Filters:     []admin.SegmentFilter{{Type: "country", Name: "country", Operator: "in", Values: []string{"US", "UK"}}},
		Enabled:     ptr(true),
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestSegmentRepo_Update(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.SegmentRepo{DB: tx}

	attrs := admin.SegmentAttrs{
		Name:        "Country Segment",
		Description: "Desc",
		AppID:       1,
		Filters:     []admin.SegmentFilter{{Type: "country", Name: "country", Operator: "in", Values: []string{"US", "UK"}}},
		Enabled:     ptr(true),
	}

	segment, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, segment, nil)
	}

	want := segment
	want.AppID = 2

	want.Enabled = ptr(false)

	updateParams := &admin.SegmentAttrs{
		AppID:   want.AppID,
		Enabled: ptr(false),
	}
	got, err := repo.Update(context.Background(), segment.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestSegmentRepo_Delete(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.SegmentRepo{DB: tx}

	attrs := &admin.SegmentAttrs{
		Name:        "Country Segment",
		Description: "Desc",
		AppID:       1,
		Filters:     []admin.SegmentFilter{{Type: "country", Name: "country", Operator: "in", Values: []string{"US", "UK"}}},
		Enabled:     ptr(true),
	}
	segment, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, segment, nil)
	}

	err = repo.Delete(context.Background(), segment.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", segment.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), segment.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", segment.ID, got, err, nil, "record not found")
	}
}
