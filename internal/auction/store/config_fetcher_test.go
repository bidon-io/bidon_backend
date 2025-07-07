package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/go-cmp/cmp"
	"github.com/lib/pq"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
)

func TestConfigFetcher_Match(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	apps := make([]db.App, 5)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx)
	}

	app4segment := dbtest.CreateSegment(t, tx, func(segment *db.Segment) {
		segment.App = apps[4]
	})

	isDefault := true
	configs := []db.AuctionConfiguration{
		{
			AppID:     apps[0].ID,
			PublicUID: sql.NullInt64{Int64: 1, Valid: true},
			AdType:    db.BannerAdType,
		},
		{
			AppID:     apps[1].ID,
			PublicUID: sql.NullInt64{Int64: 2, Valid: true},
			AdType:    db.BannerAdType,
		},
		{
			AppID:     apps[2].ID,
			PublicUID: sql.NullInt64{Int64: 3, Valid: true},
			AdType:    db.InterstitialAdType,
			CreatedAt: time.Now(),
		},
		{
			AppID:     apps[3].ID,
			PublicUID: sql.NullInt64{Int64: 4, Valid: true},
			AdType:    db.InterstitialAdType,
			CreatedAt: time.Now(),
			Demands:   pq.StringArray{"gam", "dtexchange"},
			Bidding:   pq.StringArray{"bidmachine", "mintegral"},
			AdUnitIds: pq.Int64Array{1, 2, 3},
		},
		{
			AppID:     apps[4].ID,
			PublicUID: sql.NullInt64{Int64: 5, Valid: true},
			AdType:    db.InterstitialAdType,
			CreatedAt: time.Now(),
			IsDefault: &isDefault,
		},
		{
			AppID:     apps[4].ID,
			PublicUID: sql.NullInt64{Int64: 6, Valid: true},
			AdType:    db.InterstitialAdType,
			CreatedAt: time.Now(),
		},
		{
			AppID:     apps[4].ID,
			PublicUID: sql.NullInt64{Int64: 7, Valid: true},
			AdType:    db.InterstitialAdType,
			CreatedAt: time.Now(),
		},
		{
			AppID:     apps[4].ID,
			PublicUID: sql.NullInt64{Int64: 8, Valid: true},
			AdType:    db.InterstitialAdType,
			CreatedAt: time.Now(),
			SegmentID: &sql.NullInt64{Int64: app4segment.ID, Valid: true},
		},
	}
	if err := tx.Create(&configs).Error; err != nil {
		t.Fatalf("Error creating configs: %v", err)
	}
	app1BannerConfig := &configs[0]
	app2BannerConfig := &configs[1]
	app2InterstitialConfig := &configs[2]
	app3InterstitialConfig := &configs[3]
	app4DefaultInterstitialConfig := &configs[4]
	app4SegmentInterstitialConfig := &configs[7]

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
				ID:        app1BannerConfig.ID,
				UID:       strconv.FormatInt(app1BannerConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app1BannerConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app1BannerConfig.Bidding),
				AdUnitIDs: app1BannerConfig.AdUnitIds,
			},
		},
		{
			args: args{appID: apps[1].ID, adType: ad.BannerType, segmentID: 0},
			want: &auction.Config{
				ID:        app2BannerConfig.ID,
				UID:       strconv.FormatInt(app2BannerConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app2BannerConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app2BannerConfig.Bidding),
				AdUnitIDs: app2BannerConfig.AdUnitIds,
			},
		},
		{
			args: args{appID: apps[2].ID, adType: ad.InterstitialType, segmentID: 0},
			want: &auction.Config{
				ID:        app2InterstitialConfig.ID,
				UID:       strconv.FormatInt(app2InterstitialConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app2InterstitialConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app2InterstitialConfig.Bidding),
				AdUnitIDs: app2InterstitialConfig.AdUnitIds,
			},
		},
		{
			args: args{appID: apps[3].ID, adType: ad.InterstitialType, segmentID: 0},
			want: &auction.Config{
				ID:        app3InterstitialConfig.ID,
				UID:       strconv.FormatInt(app3InterstitialConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app3InterstitialConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app3InterstitialConfig.Bidding),
				AdUnitIDs: app3InterstitialConfig.AdUnitIds,
				Timeout:   int(app3InterstitialConfig.Timeout),
			},
		},
		{
			args: args{appID: apps[4].ID, adType: ad.InterstitialType, segmentID: 0},
			want: &auction.Config{
				ID:        app4DefaultInterstitialConfig.ID,
				UID:       strconv.FormatInt(app4DefaultInterstitialConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app4DefaultInterstitialConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app4DefaultInterstitialConfig.Bidding),
				AdUnitIDs: app4DefaultInterstitialConfig.AdUnitIds,
				Timeout:   int(app4DefaultInterstitialConfig.Timeout),
			},
		},
		{
			args: args{appID: apps[4].ID, adType: ad.InterstitialType, segmentID: app4segment.ID},
			want: &auction.Config{
				ID:        app4SegmentInterstitialConfig.ID,
				UID:       strconv.FormatInt(app4SegmentInterstitialConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app4SegmentInterstitialConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app4SegmentInterstitialConfig.Bidding),
				AdUnitIDs: app4SegmentInterstitialConfig.AdUnitIds,
				Timeout:   int(app4SegmentInterstitialConfig.Timeout),
			},
		},
	}

	matcher := &store.ConfigFetcher{DB: tx}
	for _, tC := range testCases {
		got, err := matcher.Match(context.Background(), tC.args.appID, tC.args.adType, tC.args.segmentID, "v1")
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

	apps := make([]db.App, 4)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx)
	}

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
			CreatedAt: time.Now(),
		},
		{
			AppID:     apps[3].ID,
			PublicUID: sql.NullInt64{Int64: 4444444444444444444, Valid: true},
			AdType:    db.InterstitialAdType,
			CreatedAt: time.Now(),
			Demands:   pq.StringArray{"gam", "dtexchange"},
			Bidding:   pq.StringArray{"bidmachine", "mintegral"},
			AdUnitIds: pq.Int64Array{1, 2, 3},
			Timeout:   1500,
		},
	}
	if err := tx.Create(&configs).Error; err != nil {
		t.Fatalf("Error creating configs: %v", err)
	}
	app1BannerConfig := &configs[0]
	app2BannerConfig := &configs[1]
	app2InterstitialConfig := &configs[2]
	app3InterstitialConfig := &configs[3]

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
				ID:        app1BannerConfig.ID,
				UID:       strconv.FormatInt(app1BannerConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app1BannerConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app1BannerConfig.Bidding),
				AdUnitIDs: app1BannerConfig.AdUnitIds,
			},
		},
		{
			args: args{appID: apps[1].ID, uid: fmt.Sprint(app2BannerConfig.PublicUID.Int64), id: ""},
			want: &auction.Config{
				ID:        app2BannerConfig.ID,
				UID:       strconv.FormatInt(app2BannerConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app2BannerConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app2BannerConfig.Bidding),
				AdUnitIDs: app2BannerConfig.AdUnitIds,
			},
		},
		{
			args: args{appID: apps[2].ID, uid: "", id: fmt.Sprint(app2InterstitialConfig.ID)},
			want: &auction.Config{
				ID:        app2InterstitialConfig.ID,
				UID:       strconv.FormatInt(app2InterstitialConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app2InterstitialConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app2InterstitialConfig.Bidding),
				AdUnitIDs: app2InterstitialConfig.AdUnitIds,
			},
		},
		{
			args: args{appID: apps[3].ID, uid: "", id: fmt.Sprint(app3InterstitialConfig.ID)},
			want: &auction.Config{
				ID:        app3InterstitialConfig.ID,
				UID:       strconv.FormatInt(app3InterstitialConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app3InterstitialConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app3InterstitialConfig.Bidding),
				AdUnitIDs: app3InterstitialConfig.AdUnitIds,
				Timeout:   int(app3InterstitialConfig.Timeout),
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

	rdb, _ := redismock.NewClusterMock()

	apps := make([]db.App, 3)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx)
	}

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
			CreatedAt: time.Now(),
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
				ID:        app1BannerConfig.ID,
				UID:       strconv.FormatInt(app1BannerConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app1BannerConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app1BannerConfig.Bidding),
				AdUnitIDs: app1BannerConfig.AdUnitIds,
			},
		},
		{
			args: args{appID: apps[1].ID, uid: fmt.Sprint(app2BannerConfig.PublicUID.Int64), id: ""},
			want: &auction.Config{
				ID:        app2BannerConfig.ID,
				UID:       strconv.FormatInt(app2BannerConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&app2BannerConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&app2BannerConfig.Bidding),
				AdUnitIDs: app2BannerConfig.AdUnitIds,
			},
		},
		{
			args: args{appID: apps[2].ID, uid: "", id: fmt.Sprint(latestConfig.ID)},
			want: &auction.Config{
				ID:        latestConfig.ID,
				UID:       strconv.FormatInt(latestConfig.PublicUID.Int64, 10),
				Demands:   db.StringArrayToAdapterKeys(&latestConfig.Demands),
				Bidding:   db.StringArrayToAdapterKeys(&latestConfig.Bidding),
				AdUnitIDs: latestConfig.AdUnitIds,
			},
		},
		{
			args: args{appID: apps[2].ID, uid: "", id: ""},
			want: nil,
		},
	}

	configCache := config.NewRedisCacheOf[*auction.Config](rdb, 10*time.Minute, "auction_configs")
	matcher := &store.ConfigFetcher{DB: tx, Cache: configCache}
	for _, tC := range testCases {
		got := matcher.FetchByUIDCached(context.Background(), tC.args.appID, tC.args.id, tC.args.uid)

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("matcher.FetchByUIDCached(ctx, %d, %q, %q) mismatch (-want, +got):\n%s", tC.args.appID, tC.args.id, tC.args.uid, diff)
		}
	}
}
