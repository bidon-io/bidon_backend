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

func TestConfigFetcher_FetchBidMachinePlacements(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	// Create test apps
	apps := make([]db.App, 3)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx)
	}

	// Create BidMachine demand source
	bidmachineDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = "bidmachine"
		source.HumanName = "BidMachine"
	})

	// Create other demand source for negative testing
	otherDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = "other"
		source.HumanName = "Other"
	})

	// Create demand source accounts
	bidmachineAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = bidmachineDemandSource
	})

	otherAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = otherDemandSource
	})

	// Create line items with BidMachine placements
	lineItem1 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = apps[0]
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"placement": "placement-1",
		}
	})

	lineItem2 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = apps[0]
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"placement": "placement-2",
		}
	})

	// Line item without placement
	lineItem3 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = apps[1]
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"other_field": "value",
		}
	})

	// Line item from different demand source
	lineItem4 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = apps[1]
		item.Account = otherAccount
		item.Extra = map[string]any{
			"placement": "other-placement",
		}
	})

	// Create auction configurations
	configs := []db.AuctionConfiguration{
		// Config with BidMachine in demands
		{
			AppID:      apps[0].ID,
			PublicUID:  sql.NullInt64{Int64: 1, Valid: true},
			AdType:     db.BannerAdType,
			AuctionKey: "auction-key-1",
			Demands:    pq.StringArray{"bidmachine", "gam"},
			AdUnitIds:  pq.Int64Array{lineItem1.ID},
		},
		// Config with BidMachine in bidding
		{
			AppID:      apps[0].ID,
			PublicUID:  sql.NullInt64{Int64: 2, Valid: true},
			AdType:     db.InterstitialAdType,
			AuctionKey: "auction-key-2",
			Bidding:    pq.StringArray{"bidmachine", "mintegral"},
			AdUnitIds:  pq.Int64Array{lineItem2.ID},
		},
		// Config without BidMachine
		{
			AppID:      apps[1].ID,
			PublicUID:  sql.NullInt64{Int64: 3, Valid: true},
			AdType:     db.BannerAdType,
			AuctionKey: "auction-key-3",
			Demands:    pq.StringArray{"gam", "dtexchange"},
			AdUnitIds:  pq.Int64Array{lineItem4.ID},
		},
		// Config with BidMachine but line item without placement
		{
			AppID:      apps[1].ID,
			PublicUID:  sql.NullInt64{Int64: 4, Valid: true},
			AdType:     db.RewardedAdType,
			AuctionKey: "auction-key-4",
			Demands:    pq.StringArray{"bidmachine"},
			AdUnitIds:  pq.Int64Array{lineItem3.ID},
		},

	}

	if err := tx.Create(&configs).Error; err != nil {
		t.Fatalf("Error creating configs: %v", err)
	}

	testCases := []struct {
		name   string
		appID  int64
		want   map[string]string
		hasErr bool
	}{
		{
			name:  "App with BidMachine placements",
			appID: apps[0].ID,
			want: map[string]string{
				"auction-key-1": "placement-1",
				"auction-key-2": "placement-2",
			},
			hasErr: false,
		},
		{
			name:   "App without BidMachine placements",
			appID:  apps[1].ID,
			want:   map[string]string{},
			hasErr: false,
		},
		{
			name:   "App with no auction configurations",
			appID:  apps[2].ID,
			want:   map[string]string{},
			hasErr: false,
		},
		{
			name:   "Non-existent app",
			appID:  99999,
			want:   map[string]string{},
			hasErr: false,
		},
	}

	fetcher := &store.ConfigFetcher{DB: tx}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fetcher.FetchBidMachinePlacements(context.Background(), tc.appID)

			if tc.hasErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.hasErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("FetchBidMachinePlacements() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestConfigFetcher_FetchBidMachinePlacements_MultipleLineItemsSameAuctionKey(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	// Create test app
	app := dbtest.CreateApp(t, tx)

	// Create BidMachine demand source
	bidmachineDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = "bidmachine"
		source.HumanName = "BidMachine"
	})

	// Create demand source account
	bidmachineAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = bidmachineDemandSource
	})

	// Create multiple line items with different placements
	lineItem1 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"placement": "first-placement",
		}
	})

	lineItem2 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"placement": "second-placement",
		}
	})

	// Create auction configuration that uses both line items
	config := db.AuctionConfiguration{
		AppID:      app.ID,
		PublicUID:  sql.NullInt64{Int64: 1, Valid: true},
		AdType:     db.BannerAdType,
		AuctionKey: "same-auction-key",
		Demands:    pq.StringArray{"bidmachine"},
		AdUnitIds:  pq.Int64Array{lineItem1.ID, lineItem2.ID},
	}

	if err := tx.Create(&config).Error; err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	fetcher := &store.ConfigFetcher{DB: tx}
	got, err := fetcher.FetchBidMachinePlacements(context.Background(), app.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should return only one placement (the first one found)
	if len(got) != 1 {
		t.Errorf("Expected 1 placement, got %d: %v", len(got), got)
	}

	placement, exists := got["same-auction-key"]
	if !exists {
		t.Errorf("Expected auction key 'same-auction-key' to exist in result")
	}

	// Should be one of the two placements (order depends on database query result)
	if placement != "first-placement" && placement != "second-placement" {
		t.Errorf("Expected placement to be 'first-placement' or 'second-placement', got %s", placement)
	}
}

func TestConfigFetcher_FetchBidMachinePlacements_EdgeCases(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	// Create test app
	app := dbtest.CreateApp(t, tx)

	// Create BidMachine demand source
	bidmachineDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = "bidmachine"
		source.HumanName = "BidMachine"
	})

	// Create demand source account
	bidmachineAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = bidmachineDemandSource
	})

	// Create line item with empty placement
	lineItem1 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"placement": "",
		}
	})

	// Create line item with null placement
	lineItem2 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"placement": nil,
		}
	})

	// Create line item with valid placement
	lineItem3 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"placement": "valid-placement",
		}
	})

	// Create auction configurations
	configs := []db.AuctionConfiguration{
		// Config with line item having empty placement (should be ignored)
		{
			AppID:      app.ID,
			PublicUID:  sql.NullInt64{Int64: 1, Valid: true},
			AdType:     db.BannerAdType,
			AuctionKey: "auction-key-empty",
			Demands:    pq.StringArray{"bidmachine"},
			AdUnitIds:  pq.Int64Array{lineItem1.ID},
		},
		// Config with line item having null placement (should be ignored)
		{
			AppID:      app.ID,
			PublicUID:  sql.NullInt64{Int64: 2, Valid: true},
			AdType:     db.InterstitialAdType,
			AuctionKey: "auction-key-null",
			Demands:    pq.StringArray{"bidmachine"},
			AdUnitIds:  pq.Int64Array{lineItem2.ID},
		},
		// Config with valid placement
		{
			AppID:      app.ID,
			PublicUID:  sql.NullInt64{Int64: 3, Valid: true},
			AdType:     db.RewardedAdType,
			AuctionKey: "auction-key-valid",
			Demands:    pq.StringArray{"bidmachine"},
			AdUnitIds:  pq.Int64Array{lineItem3.ID},
		},
	}

	if err := tx.Create(&configs).Error; err != nil {
		t.Fatalf("Error creating configs: %v", err)
	}

	fetcher := &store.ConfigFetcher{DB: tx}
	got, err := fetcher.FetchBidMachinePlacements(context.Background(), app.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only return the valid placement
	want := map[string]string{
		"auction-key-valid": "valid-placement",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FetchBidMachinePlacements() mismatch (-want, +got):\n%s", diff)
	}
}

func TestConfigFetcher_FetchBidMachinePlacements_NoBidMachineConfigs(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	// Create test app
	app := dbtest.CreateApp(t, tx)

	// Create non-BidMachine demand sources
	gamDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = "gam"
		source.HumanName = "Google Ad Manager"
	})

	dtexchangeDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = "dtexchange"
		source.HumanName = "DT Exchange"
	})

	// Create demand source accounts
	gamAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = gamDemandSource
	})

	dtexchangeAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = dtexchangeDemandSource
	})

	// Create line items with placements (but not BidMachine)
	lineItem1 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = gamAccount
		item.Extra = map[string]any{
			"placement": "gam-placement",
		}
	})

	lineItem2 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = dtexchangeAccount
		item.Extra = map[string]any{
			"placement": "dtexchange-placement",
		}
	})

	// Create auction configurations without BidMachine
	configs := []db.AuctionConfiguration{
		{
			AppID:      app.ID,
			PublicUID:  sql.NullInt64{Int64: 1, Valid: true},
			AdType:     db.BannerAdType,
			Demands:    pq.StringArray{"gam", "dtexchange"},
			AdUnitIds:  pq.Int64Array{lineItem1.ID},
		},
		{
			AppID:     app.ID,
			PublicUID: sql.NullInt64{Int64: 2, Valid: true},
			AdType:    db.InterstitialAdType,
			Bidding:   pq.StringArray{"mintegral", "unityads"},
			AdUnitIds: pq.Int64Array{lineItem2.ID},
		},
	}

	if err := tx.Create(&configs).Error; err != nil {
		t.Fatalf("Error creating configs: %v", err)
	}

	fetcher := &store.ConfigFetcher{DB: tx}
	got, err := fetcher.FetchBidMachinePlacements(context.Background(), app.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should return empty map since no configurations contain BidMachine
	want := map[string]string{}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FetchBidMachinePlacements() mismatch (-want, +got):\n%s", diff)
	}
}

func TestConfigFetcher_FetchBidMachinePlacements_MissingPlacement(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	// Create test app
	app := dbtest.CreateApp(t, tx)

	// Create BidMachine demand source
	bidmachineDemandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = "bidmachine"
		source.HumanName = "BidMachine"
	})

	// Create demand source account
	bidmachineAccount := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = bidmachineDemandSource
	})

	// Create line items with different placement scenarios
	lineItemWithPlacement := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"placement": "valid-placement",
		}
	})

	lineItemWithoutPlacement := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = bidmachineAccount
		item.Extra = map[string]any{
			"other_field": "some_value",
			// no placement field
		}
	})

	lineItemWithEmptyExtra := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = bidmachineAccount
		item.Extra = map[string]any{} // empty extra
	})

	// Create auction configurations
	configs := []db.AuctionConfiguration{
		// Config with line item that has placement
		{
			AppID:      app.ID,
			PublicUID:  sql.NullInt64{Int64: 1, Valid: true},
			AdType:     db.BannerAdType,
			Demands:    pq.StringArray{"bidmachine"},
			AdUnitIds:  pq.Int64Array{lineItemWithPlacement.ID},
		},
		// Config with line item that doesn't have placement (should be ignored)
		{
			AppID:      app.ID,
			PublicUID:  sql.NullInt64{Int64: 2, Valid: true},
			AdType:     db.InterstitialAdType,
			Demands:    pq.StringArray{"bidmachine"},
			AdUnitIds:  pq.Int64Array{lineItemWithoutPlacement.ID},
		},
		// Config with line item that has empty extra (should be ignored)
		{
			AppID:      app.ID,
			PublicUID:  sql.NullInt64{Int64: 3, Valid: true},
			AdType:     db.RewardedAdType,
			Demands:    pq.StringArray{"bidmachine"},
			AdUnitIds:  pq.Int64Array{lineItemWithEmptyExtra.ID},
		},
		// Config with multiple line items (mixed placement availability)
		{
			AppID:      app.ID,
			PublicUID:  sql.NullInt64{Int64: 4, Valid: true},
			AdType:     db.BannerAdType,
			Bidding:    pq.StringArray{"bidmachine"},
			AdUnitIds:  pq.Int64Array{lineItemWithPlacement.ID, lineItemWithoutPlacement.ID, lineItemWithEmptyExtra.ID},
		},
	}

	if err := tx.Create(&configs).Error; err != nil {
		t.Fatalf("Error creating configs: %v", err)
	}

	fetcher := &store.ConfigFetcher{DB: tx}
	got, err := fetcher.FetchBidMachinePlacements(context.Background(), app.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only return placements for line items that actually have placement field
	// The auction_key values will be generated based on PublicUID, so we need to get them
	var dbConfigs []db.AuctionConfiguration
	if err := tx.Where("app_id = ? AND public_uid IN (?, ?)", app.ID, 1, 4).Find(&dbConfigs).Error; err != nil {
		t.Fatalf("Error fetching configs: %v", err)
	}

	want := map[string]string{}
	for _, config := range dbConfigs {
		want[config.AuctionKey] = "valid-placement"
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FetchBidMachinePlacements() mismatch (-want, +got):\n%s", diff)
	}

	// Verify that we got exactly 2 results (from configs with PublicUID 1 and 4)
	if len(got) != 2 {
		t.Errorf("Expected 2 placements, got %d: %v", len(got), got)
	}
}
