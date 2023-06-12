package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/google/go-cmp/cmp"
)

func TestAuctionConfigurationRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewAuctionConfigurationRepo(tx)

	configs := []admin.AuctionConfigurationAttrs{
		{
			Name:       "Config 1",
			AppID:      1,
			AdType:     ad.BannerType,
			Rounds:     []admin.AuctionRoundConfiguration{{ID: "1", Demands: []string{"demand1", "demand2"}, Timeout: 10}},
			Pricefloor: 0.5,
		},
		{
			Name:       "Config 2",
			AppID:      2,
			AdType:     ad.InterstitialType,
			Rounds:     []admin.AuctionRoundConfiguration{{ID: "2", Demands: []string{"demand3", "demand4"}, Timeout: 20}},
			Pricefloor: 0.75,
		},
		{
			Name:       "Config 3",
			AppID:      3,
			AdType:     ad.RewardedType,
			Rounds:     []admin.AuctionRoundConfiguration{{ID: "3", Demands: []string{"demand5", "demand6"}, Timeout: 30}},
			Pricefloor: 1.0,
		},
	}

	want := make([]admin.AuctionConfiguration, len(configs))
	for i, attrs := range configs {
		config, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, config, nil)
		}

		want[i] = *config
	}

	got, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestAuctionConfigurationRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewAuctionConfigurationRepo(tx)

	attrs := &admin.AuctionConfigurationAttrs{
		Name:       "Config 1",
		AppID:      1,
		AdType:     ad.BannerType,
		Rounds:     []admin.AuctionRoundConfiguration{{ID: "1", Demands: []string{"demand1", "demand2"}, Timeout: 10}},
		Pricefloor: 0.5,
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

func TestAuctionConfigurationRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewAuctionConfigurationRepo(tx)

	attrs := admin.AuctionConfigurationAttrs{
		Name:       "Config 1",
		AppID:      1,
		AdType:     ad.BannerType,
		Rounds:     []admin.AuctionRoundConfiguration{{ID: "1", Demands: []string{"demand1", "demand2"}, Timeout: 10}},
		Pricefloor: 0.5,
	}

	config, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, config, nil)
	}

	want := config
	want.AppID = 2

	updateParams := &admin.AuctionConfigurationAttrs{
		AppID: want.AppID,
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

	repo := store.NewAuctionConfigurationRepo(tx)

	attrs := &admin.AuctionConfigurationAttrs{
		Name:       "Config 1",
		AppID:      1,
		AdType:     ad.BannerType,
		Rounds:     []admin.AuctionRoundConfiguration{{ID: "1", Demands: []string{"demand1", "demand2"}, Timeout: 10}},
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
