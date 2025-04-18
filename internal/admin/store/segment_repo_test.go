package adminstore_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/resource"
	adminstore "github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/bidon-io/bidon-backend/internal/segment"
)

func TestSegmentRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewSegmentRepo(tx)

	apps := make([]db.App, 3)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx)
	}

	segments := []admin.SegmentAttrs{
		{
			Name:        "Country Segment",
			Description: "Desc",
			AppID:       apps[0].ID,
			Filters:     []segment.Filter{{Type: "country", Name: "country", Operator: "in", Values: []string{"US", "UK"}}},
			Enabled:     ptr(true),
			Priority:    1,
		},
		{
			Name:        "Custom String Segment",
			Description: "Desc",
			AppID:       apps[1].ID,
			Filters:     []segment.Filter{{Type: "string", Name: "custom_str", Operator: "==", Values: []string{"super"}}},
			Enabled:     ptr(false),
			Priority:    1,
		},
		{
			Name:        "Custom Num Segment",
			Description: "Desc",
			AppID:       apps[2].ID,
			Filters:     []segment.Filter{{Type: "float", Name: "custom_num", Operator: ">=", Values: []string{"33"}}},
			Enabled:     ptr(true),
			Priority:    0,
		},
	}

	items := make([]admin.Segment, len(segments))
	for i, attrs := range segments {
		segment, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; items %T, %v", &attrs, nil, err, segment, nil)
		}

		items[i] = *segment
		items[i].App = adminstore.AppAttrsWithId(&apps[i])
	}

	want := &resource.Collection[admin.Segment]{
		Items: items,
		Meta:  resource.CollectionMeta{TotalCount: int64(len(items))},
	}

	got, err := repo.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestSegmentRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewSegmentRepo(tx)

	app := dbtest.CreateApp(t, tx)
	attrs := &admin.SegmentAttrs{
		Name:        "Country Segment",
		Description: "Desc",
		AppID:       app.ID,
		Filters:     []segment.Filter{{Type: "country", Name: "country", Operator: "in", Values: []string{"US", "UK"}}},
		Enabled:     ptr(true),
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}
	want.App = adminstore.AppAttrsWithId(&app)

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestSegmentRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewSegmentRepo(tx)

	app := dbtest.CreateApp(t, tx)
	attrs := admin.SegmentAttrs{
		Name:        "Country Segment",
		Description: "Desc",
		AppID:       app.ID,
		Filters:     []segment.Filter{{Type: "country", Name: "country", Operator: "in", Values: []string{"US", "UK"}}},
		Enabled:     ptr(true),
		Priority:    1,
	}

	segment, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, segment, nil)
	}

	want := segment
	want.Description = "New Desc"

	want.Enabled = ptr(false)

	updateParams := &admin.SegmentAttrs{
		Description: want.Description,
		Enabled:     ptr(false),
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
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewSegmentRepo(tx)

	app := dbtest.CreateApp(t, tx)
	attrs := &admin.SegmentAttrs{
		Name:        "Country Segment",
		Description: "Desc",
		AppID:       app.ID,
		Filters:     []segment.Filter{{Type: "country", Name: "country", Operator: "in", Values: []string{"US", "UK"}}},
		Enabled:     ptr(true),
		Priority:    2,
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
