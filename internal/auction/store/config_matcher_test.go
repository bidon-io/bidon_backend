package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/google/go-cmp/cmp"
)

func TestConfigMatcher_Match(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	configs := []db.AuctionConfiguration{
		{AppID: 1, AdType: db.BannerAdType},
		{AppID: 2, AdType: db.BannerAdType},
		{AppID: 3, AdType: db.InterstitialAdType, Model: db.Model{CreatedAt: time.Now()}},
		{AppID: 3, AdType: db.InterstitialAdType, Model: db.Model{CreatedAt: time.Now().Add(-time.Hour)}},
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
			args: args{appID: 1, adType: ad.BannerType, segmentID: 0},
			want: &auction.Config{ID: app1BannerConfig.ID, Rounds: app1BannerConfig.Rounds},
		},
		{
			args: args{appID: 2, adType: ad.BannerType, segmentID: 0},
			want: &auction.Config{ID: app2BannerConfig.ID, Rounds: app2BannerConfig.Rounds},
		},
		{
			args: args{appID: 3, adType: ad.InterstitialType, segmentID: 0},
			want: &auction.Config{ID: latestConfig.ID, Rounds: latestConfig.Rounds},
		},
	}

	matcher := &store.ConfigMatcher{DB: tx}
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
