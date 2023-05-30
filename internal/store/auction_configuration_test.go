package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/store"
	"github.com/google/go-cmp/cmp"
)

func TestAuctionConfigurationRepo_List(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AuctionConfigurationRepo{DB: tx}

	configs := []auction.ConfigurationAttrs{
		{
			Name:       "Config 1",
			AppID:      1,
			AdType:     auction.BannerAdType,
			Rounds:     []auction.RoundConfiguration{{ID: "1", Demands: []string{"demand1", "demand2"}, Timeout: 10}},
			Pricefloor: 0.5,
		},
		{
			Name:       "Config 2",
			AppID:      2,
			AdType:     auction.InterstitialAdType,
			Rounds:     []auction.RoundConfiguration{{ID: "2", Demands: []string{"demand3", "demand4"}, Timeout: 20}},
			Pricefloor: 0.75,
		},
		{
			Name:       "Config 3",
			AppID:      3,
			AdType:     auction.RewardedAdType,
			Rounds:     []auction.RoundConfiguration{{ID: "3", Demands: []string{"demand5", "demand6"}, Timeout: 30}},
			Pricefloor: 1.0,
		},
	}

	want := make([]auction.Configuration, len(configs))
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
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AuctionConfigurationRepo{DB: tx}

	attrs := &auction.ConfigurationAttrs{
		Name:       "Config 1",
		AppID:      1,
		AdType:     auction.BannerAdType,
		Rounds:     []auction.RoundConfiguration{{ID: "1", Demands: []string{"demand1", "demand2"}, Timeout: 10}},
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
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AuctionConfigurationRepo{DB: tx}

	attrs := auction.ConfigurationAttrs{
		Name:       "Config 1",
		AppID:      1,
		AdType:     auction.BannerAdType,
		Rounds:     []auction.RoundConfiguration{{ID: "1", Demands: []string{"demand1", "demand2"}, Timeout: 10}},
		Pricefloor: 0.5,
	}

	config, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, config, nil)
	}

	want := config
	want.AppID = 2

	updateParams := &auction.ConfigurationAttrs{
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
	tx := db.Begin()
	defer tx.Rollback()

	repo := &store.AuctionConfigurationRepo{DB: tx}

	attrs := &auction.ConfigurationAttrs{
		Name:       "Config 1",
		AppID:      1,
		AdType:     auction.BannerAdType,
		Rounds:     []auction.RoundConfiguration{{ID: "1", Demands: []string{"demand1", "demand2"}, Timeout: 10}},
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
