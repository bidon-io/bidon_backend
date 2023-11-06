package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/google/go-cmp/cmp"
)

func TestConfigFetcher_Match(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	apps := dbtest.CreateAppsList(t, tx, 3)
	configs := []db.AuctionConfiguration{
		{
			AppID:     apps[0].ID,
			PublicUID: sql.NullInt64{Int64: 1111111111111111111, Valid: true},
			AdType:    db.BannerAdType,
		},
		{
			AppID:     apps[1].ID,
			PublicUID: sql.NullInt64{Int64: 2222222222222222222, Valid: true},
			AdType:    db.BannerAdType,
		},
		{
			AppID:     apps[2].ID,
			PublicUID: sql.NullInt64{Int64: 3333333333333333333, Valid: true},
			AdType:    db.InterstitialAdType,
			Model:     db.Model{CreatedAt: time.Now()},
		},
	}
	if err := tx.Create(&configs).Error; err != nil {
		t.Fatalf("Error creating configs: %v", err)
	}
	app1BannerConfig := &configs[0]
	app2BannerConfig := &configs[1]
	latestConfig := &configs[2]

	type args struct {
		appID     int64
		adType    ad.Type
		segmentID int64
	}
	testCases := []struct {
		args args
		want *auction.Config
	}{
		{
			args: args{appID: apps[0].ID, adType: ad.BannerType, segmentID: 0},
			want: &auction.Config{
				ID:     app1BannerConfig.ID,
				UID:    strconv.FormatInt(app1BannerConfig.PublicUID.Int64, 10),
				Rounds: app1BannerConfig.Rounds,
			},
		},
		{
			args: args{appID: apps[1].ID, adType: ad.BannerType, segmentID: 0},
			want: &auction.Config{
				ID:     app2BannerConfig.ID,
				UID:    strconv.FormatInt(app2BannerConfig.PublicUID.Int64, 10),
				Rounds: app2BannerConfig.Rounds,
			},
		},
		{
			args: args{appID: apps[2].ID, adType: ad.InterstitialType, segmentID: 0},
			want: &auction.Config{
				ID:     latestConfig.ID,
				UID:    strconv.FormatInt(latestConfig.PublicUID.Int64, 10),
				Rounds: latestConfig.Rounds,
			},
		},
	}

	matcher := &store.ConfigFetcher{DB: tx}
	for _, tC := range testCases {
		got, err := matcher.Match(context.Background(), tC.args.appID, tC.args.adType, tC.args.segmentID)
		if err != nil {
			t.Errorf("Error matching config: %v", err)
		}

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("matcher.Match(ctx, %d, %q) mismatch (-want, +got):\n%s", tC.args.appID, tC.args.adType, diff)
		}
	}
}

func TestConfigFetcher_FetchByUID(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	apps := dbtest.CreateAppsList(t, tx, 3)
	configs := []db.AuctionConfiguration{
		{
			AppID:     apps[0].ID,
			PublicUID: sql.NullInt64{Int64: 1111111111111111111, Valid: true},
			AdType:    db.BannerAdType,
		},
		{
			AppID:     apps[1].ID,
			PublicUID: sql.NullInt64{Int64: 2222222222222222222, Valid: true},
			AdType:    db.BannerAdType,
		},
		{
			AppID:     apps[2].ID,
			PublicUID: sql.NullInt64{Int64: 3333333333333333333, Valid: true},
			AdType:    db.InterstitialAdType,
			Model:     db.Model{CreatedAt: time.Now()},
		},
	}
	if err := tx.Create(&configs).Error; err != nil {
		t.Fatalf("Error creating configs: %v", err)
	}
	app1BannerConfig := &configs[0]
	app2BannerConfig := &configs[1]
	latestConfig := &configs[2]

	type args struct {
		appID int64
		id    string
		uid   string
	}
	testCases := []struct {
		args args
		want *auction.Config
	}{
		{
			args: args{appID: apps[0].ID, uid: "", id: fmt.Sprint(app1BannerConfig.ID)},
			want: &auction.Config{
				ID:     app1BannerConfig.ID,
				UID:    strconv.FormatInt(app1BannerConfig.PublicUID.Int64, 10),
				Rounds: app1BannerConfig.Rounds,
			},
		},
		{
			args: args{appID: apps[1].ID, uid: fmt.Sprint(app2BannerConfig.PublicUID.Int64), id: ""},
			want: &auction.Config{
				ID:     app2BannerConfig.ID,
				UID:    strconv.FormatInt(app2BannerConfig.PublicUID.Int64, 10),
				Rounds: app2BannerConfig.Rounds,
			},
		},
		{
			args: args{appID: apps[2].ID, uid: "", id: fmt.Sprint(latestConfig.ID)},
			want: &auction.Config{
				ID:     latestConfig.ID,
				UID:    strconv.FormatInt(latestConfig.PublicUID.Int64, 10),
				Rounds: latestConfig.Rounds,
			},
		},
		{
			args: args{appID: apps[2].ID, uid: "", id: ""},
			want: nil,
		},
	}

	matcher := &store.ConfigFetcher{DB: tx}
	for _, tC := range testCases {
		got := matcher.FetchByUID(context.Background(), tC.args.appID, tC.args.id, tC.args.uid)

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("matcher.FetchByUID(ctx, %d, %q, %q) mismatch (-want, +got):\n%s", tC.args.appID, tC.args.id, tC.args.uid, diff)
		}
	}
}

func TestConfigFetcher_FetchByUIDCached(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	apps := dbtest.CreateAppsList(t, tx, 3)
	configs := []db.AuctionConfiguration{
		{
			AppID:     apps[0].ID,
			PublicUID: sql.NullInt64{Int64: 1111111111111111111, Valid: true},
			AdType:    db.BannerAdType,
		},
		{
			AppID:     apps[1].ID,
			PublicUID: sql.NullInt64{Int64: 2222222222222222222, Valid: true},
			AdType:    db.BannerAdType,
		},
		{
			AppID:     apps[2].ID,
			PublicUID: sql.NullInt64{Int64: 3333333333333333333, Valid: true},
			AdType:    db.InterstitialAdType,
			Model:     db.Model{CreatedAt: time.Now()},
		},
	}
	if err := tx.Create(&configs).Error; err != nil {
		t.Fatalf("Error creating configs: %v", err)
	}
	app1BannerConfig := &configs[0]
	app2BannerConfig := &configs[1]
	latestConfig := &configs[2]

	type args struct {
		appID int64
		id    string
		uid   string
	}
	testCases := []struct {
		args args
		want *auction.Config
	}{
		{
			args: args{appID: apps[0].ID, uid: "", id: fmt.Sprint(app1BannerConfig.ID)},
			want: &auction.Config{
				ID:     app1BannerConfig.ID,
				UID:    strconv.FormatInt(app1BannerConfig.PublicUID.Int64, 10),
				Rounds: app1BannerConfig.Rounds,
			},
		},
		{
			args: args{appID: apps[1].ID, uid: fmt.Sprint(app2BannerConfig.PublicUID.Int64), id: ""},
			want: &auction.Config{
				ID:     app2BannerConfig.ID,
				UID:    strconv.FormatInt(app2BannerConfig.PublicUID.Int64, 10),
				Rounds: app2BannerConfig.Rounds,
			},
		},
		{
			args: args{appID: apps[2].ID, uid: "", id: fmt.Sprint(latestConfig.ID)},
			want: &auction.Config{
				ID:     latestConfig.ID,
				UID:    strconv.FormatInt(latestConfig.PublicUID.Int64, 10),
				Rounds: latestConfig.Rounds,
			},
		},
		{
			args: args{appID: apps[2].ID, uid: "", id: ""},
			want: nil,
		},
	}

	configCache := config.NewMemoryCacheOf[*auction.Config](30*time.Second, 1*time.Second, 1*time.Hour)

	matcher := &store.ConfigFetcher{DB: tx, Cache: configCache}
	for _, tC := range testCases {
		got := matcher.FetchByUIDCached(context.Background(), tC.args.appID, tC.args.id, tC.args.uid)

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("matcher.FetchByUIDCached(ctx, %d, %q, %q) mismatch (-want, +got):\n%s", tC.args.appID, tC.args.id, tC.args.uid, diff)
		}
	}
}
