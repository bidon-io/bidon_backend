package adminstore_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/resource"
	adminstore "github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
)

func TestAuctionConfigurationRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAuctionConfigurationRepo(tx)

	apps := make([]db.App, 3)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx)
	}

	segments := make([]db.Segment, 3)
	segments[0] = dbtest.CreateSegment(t, tx, func(segment *db.Segment) {
		segment.App = apps[0]
	})
	configs := []admin.AuctionConfigurationAttrs{
		{
			Name:       "Config 1",
			AppID:      apps[0].ID,
			AdType:     ad.BannerType,
			Pricefloor: 0.5,
			SegmentID:  &segments[0].ID,
		},
		{
			Name:       "Config 2",
			AppID:      apps[1].ID,
			AdType:     ad.InterstitialType,
			Pricefloor: 0.75,
		},
		{
			Name:       "Config 3",
			AppID:      apps[2].ID,
			AdType:     ad.RewardedType,
			Pricefloor: 1.0,
		},
	}

	wantItems := make([]admin.AuctionConfiguration, len(configs))
	for i, attrs := range configs {
		config, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; wantItems %T, %v", &attrs, nil, err, config, nil)
		}

		wantItems[i] = *config
		wantItems[i].App = adminstore.AppAttrsWithId(&apps[i])
		if segments[i].ID != 0 {
			wantItems[i].Segment = adminstore.SegmentAttrsWithId(&segments[i])
		}
	}

	want := &resource.Collection[admin.AuctionConfiguration]{
		Items: wantItems,
		Meta:  resource.CollectionMeta{TotalCount: int64(len(wantItems))},
	}

	got, err := repo.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; wantItems %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAuctionConfigurationRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAuctionConfigurationRepo(tx)

	app := dbtest.CreateApp(t, tx)
	attrs := &admin.AuctionConfigurationAttrs{
		Name:       "Config 1",
		AppID:      app.ID,
		AdType:     ad.BannerType,
		Pricefloor: 0.5,
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

func TestAuctionConfigurationRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAuctionConfigurationRepo(tx)

	app := dbtest.CreateApp(t, tx)
	attrs := admin.AuctionConfigurationAttrs{
		Name:       "Config 1",
		AppID:      app.ID,
		AdType:     ad.BannerType,
		Pricefloor: 0.5,
	}

	config, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, config, nil)
	}

	want := config
	want.Name = "New Name"

	updateParams := &admin.AuctionConfigurationAttrs{
		Name: want.Name,
	}
	got, err := repo.Update(context.Background(), config.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAuctionConfigurationRepo_Delete(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewAuctionConfigurationRepo(tx)

	app := dbtest.CreateApp(t, tx)
	attrs := &admin.AuctionConfigurationAttrs{
		Name:       "Config 1",
		AppID:      app.ID,
		AdType:     ad.BannerType,
		Pricefloor: 0.5,
	}
	config, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, config, nil)
	}

	err = repo.Delete(context.Background(), config.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", config.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), config.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", config.ID, got, err, nil, "record not found")
	}
}
