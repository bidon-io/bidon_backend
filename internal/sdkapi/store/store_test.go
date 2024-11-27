package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"os"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/config"

	"github.com/bidon-io/bidon-backend/internal/ad"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/google/go-cmp/cmp"
)

var testDB *db.DB

func TestMain(m *testing.M) {
	testDB = dbtest.Prepare()

	os.Exit(m.Run())
}

func TestAppFetcher_Fetch(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	user := dbtest.CreateUser(t, tx)
	app := &db.App{
		UserID:      user.ID,
		AppKey:      sql.NullString{String: "asdf", Valid: true},
		PackageName: sql.NullString{String: "com.example.app", Valid: true},
	}
	if err := tx.Create(app).Error; err != nil {
		t.Fatalf("Error creating app: %v", err)
	}

	fetcher := &AppFetcher{DB: tx, Cache: config.NewMemoryCacheOf[sdkapi.App](10 * time.Minute)}

	testCases := []struct {
		name      string
		appKey    string
		appBundle string
		want      any
	}{
		{
			name:      "App matches",
			appKey:    app.AppKey.String,
			appBundle: app.PackageName.String,
			want:      sdkapi.App{ID: app.ID},
		},
		{
			name:      "App key does not match",
			appKey:    "fdsa",
			appBundle: app.PackageName.String,
			want:      sdkapi.ErrAppNotValid,
		},
		{
			name:      "App bundle does not match",
			appKey:    app.AppKey.String,
			appBundle: "not.found",
			want:      sdkapi.ErrAppNotValid,
		},
		{
			name:      "Nothing matches",
			appKey:    "fdsa",
			appBundle: "not.found",
			want:      sdkapi.ErrAppNotValid,
		},
	}

	for _, tC := range testCases {
		app, _ := fetcher.Fetch(context.Background(), tC.appKey, tC.appBundle)
		appCached, err := fetcher.FetchCached(context.Background(), tC.appKey, tC.appBundle)

		if diff := cmp.Diff(app, appCached); diff != "" {
			t.Errorf("fetcher.FetchCached -> %v mismatch (-want +got):\n%s", tC.name, diff)
		}

		var got any
		switch tC.want.(type) {
		case sdkapi.App:
			got = app
		case error:
			got = err
		}

		if diff := cmp.Diff(tC.want, got); diff != "" {
			t.Errorf("fetcher.Fetch -> %v mismatch (-want +got):\n%s", tC.name, diff)
		}
	}
}

func TestAdapterInitConfigsFetcher_FetchAdapterInitConfigs_Valid(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	rdb, _ := redismock.NewClientMock()

	keys := adapter.Keys

	demandSources := make([]db.DemandSource, len(keys))
	for i, key := range keys {
		demandSources[i] = dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
			source.APIKey = string(key)
		})
	}

	accounts := make([]db.DemandSourceAccount, len(demandSources))
	for i, source := range demandSources {
		accounts[i] = dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
			account.DemandSource = source
			account.Extra = dbtest.ValidDemandSourceAccountExtra(t, adapter.Key(source.APIKey))
		})
	}

	apps := make([]db.App, 2)
	for i := range apps {
		apps[i] = dbtest.CreateApp(t, tx)
	}

	selectApp := func(i int) db.App {
		firstAppAccounts := []adapter.Key{adapter.AdmobKey, adapter.AmazonKey, adapter.ApplovinKey, adapter.BidmachineKey, adapter.BigoAdsKey, adapter.DTExchangeKey, adapter.GAMKey}
		for _, key := range firstAppAccounts {
			if key == adapter.Key(accounts[i].DemandSource.APIKey) {
				return apps[0]
			}
		}
		return apps[1]
	}
	for i, account := range accounts {
		dbtest.CreateAppDemandProfile(t, tx, func(profile *db.AppDemandProfile) {
			profile.App = selectApp(i)
			profile.Account = account
			profile.Data = dbtest.ValidAppDemandProfileData(t, adapter.Key(account.DemandSource.APIKey), profile.App.ID)
		})
	}

	profilesCache := config.NewRedisCacheOf[[]db.AppDemandProfile](rdb, 10*time.Minute, "app_demand_profiles")
	amazonSlotsCache := config.NewRedisCacheOf[[]sdkapi.AmazonSlot](rdb, 10*time.Minute, "amazon_slots")
	fetcher := &AdapterInitConfigsFetcher{DB: tx, ProfilesCache: profilesCache, AmazonSlotsCache: amazonSlotsCache}

	tests := []struct {
		name           string
		appID          int64
		adapterKeys    []adapter.Key
		setAmazonSlots bool
		setOrder       bool
		want           []sdkapi.AdapterInitConfig
	}{
		{
			name:           "first app with all adapters",
			appID:          apps[0].ID,
			adapterKeys:    adapter.Keys,
			setAmazonSlots: true,
			setOrder:       false,
			want: []sdkapi.AdapterInitConfig{
				&sdkapi.AdmobInitConfig{
					AppID: fmt.Sprintf("admob_app_%d", apps[0].ID),
				},
				&sdkapi.AmazonInitConfig{
					AppKey: fmt.Sprintf("amazon_app_%d", apps[0].ID),
					Slots:  []sdkapi.AmazonSlot{},
				},
				&sdkapi.ApplovinInitConfig{
					AppKey: "applovin",
					SDKKey: "applovin",
				},
				&sdkapi.BidmachineInitConfig{
					SellerID:        "1",
					Endpoint:        "x.appbaqend.com",
					MediationConfig: []string{"one", "two"},
				},
				&sdkapi.BigoAdsInitConfig{
					AppID: fmt.Sprintf("bigoads_app_%d", apps[0].ID),
				},
				&sdkapi.DTExchangeInitConfig{
					AppID: fmt.Sprintf("dtexchange_app_%d", apps[0].ID),
				},
				&sdkapi.GAMInitConfig{
					NetworkCode: "111",
					AppID:       fmt.Sprintf("gam_app_%d", apps[0].ID),
					Order:       0,
				},
			},
		},
		{
			name:           "second app with all adapters",
			appID:          apps[1].ID,
			adapterKeys:    adapter.Keys,
			setAmazonSlots: true,
			setOrder:       false,
			want: []sdkapi.AdapterInitConfig{
				&sdkapi.ChartboostInitConfig{
					AppID:        fmt.Sprintf("chartboost_app_%d", apps[1].ID),
					AppSignature: "123",
				},
				&sdkapi.InmobiInitConfig{
					AccountID: "inmobi",
					AppKey:    fmt.Sprintf("inmobi_app_%d", apps[1].ID),
				},
				&sdkapi.IronSourceInitConfig{
					AppKey: fmt.Sprintf("ironsource_app_%d", apps[1].ID),
				},
				&sdkapi.MetaInitConfig{
					AppID:     fmt.Sprintf("meta_app_%d", apps[1].ID),
					AppSecret: fmt.Sprintf("meta_app_%d_secret", apps[1].ID),
				},
				&sdkapi.MintegralInitConfig{
					AppID:  fmt.Sprintf("mintegral_app_%d", apps[1].ID),
					AppKey: "mintegral",
				},
				&sdkapi.MobileFuseInitConfig{},
				&sdkapi.UnityAdsInitConfig{
					GameID: fmt.Sprintf("unityads_game_%d", apps[1].ID),
				},
				&sdkapi.VKAdsInitConfig{
					AppID: fmt.Sprintf("vkads_app_%d", apps[1].ID),
					Order: 0,
				},
				&sdkapi.VungleInitConfig{
					AppID: fmt.Sprintf("vungle_app_%d", apps[1].ID),
				},
				&sdkapi.YandexInitConfig{
					MetricaID: fmt.Sprintf("yandex_metrica_%d", apps[1].ID),
				},
			},
		},
		{
			name:           "setOrder = true",
			appID:          apps[1].ID,
			adapterKeys:    adapter.Keys,
			setAmazonSlots: true,
			setOrder:       true,
			want: []sdkapi.AdapterInitConfig{
				&sdkapi.ChartboostInitConfig{
					AppID:        fmt.Sprintf("chartboost_app_%d", apps[1].ID),
					AppSignature: "123",
					Order:        2,
				},
				&sdkapi.InmobiInitConfig{
					AccountID: "inmobi",
					AppKey:    fmt.Sprintf("inmobi_app_%d", apps[1].ID),
					Order:     3,
				},
				&sdkapi.IronSourceInitConfig{
					AppKey: fmt.Sprintf("ironsource_app_%d", apps[1].ID),
					Order:  2,
				},
				&sdkapi.MetaInitConfig{
					AppID:     fmt.Sprintf("meta_app_%d", apps[1].ID),
					AppSecret: fmt.Sprintf("meta_app_%d_secret", apps[1].ID),
					Order:     0,
				},
				&sdkapi.MintegralInitConfig{
					AppID:  fmt.Sprintf("mintegral_app_%d", apps[1].ID),
					AppKey: "mintegral",
					Order:  3,
				},
				&sdkapi.MobileFuseInitConfig{
					Order: 3,
				},
				&sdkapi.UnityAdsInitConfig{
					GameID: fmt.Sprintf("unityads_game_%d", apps[1].ID),
					Order:  2,
				},
				&sdkapi.VKAdsInitConfig{
					AppID: fmt.Sprintf("vkads_app_%d", apps[1].ID),
					Order: 2,
				},
				&sdkapi.VungleInitConfig{
					AppID: fmt.Sprintf("vungle_app_%d", apps[1].ID),
					Order: 2,
				},
				&sdkapi.YandexInitConfig{
					MetricaID: fmt.Sprintf("yandex_metrica_%d", apps[1].ID),
					Order:     2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetcher.FetchAdapterInitConfigs(context.Background(), tt.appID, tt.adapterKeys, tt.setAmazonSlots, tt.setOrder)
			if err != nil {
				t.Fatalf("FetchAdapterInitConfigs() error = %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FetchAdapterInitConfigs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAdapterInitConfigsFetcher_FetchAdapterInitConfigs_Amazon(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	rdb, _ := redismock.NewClientMock()

	demandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.AmazonKey)
	})

	account := dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
		account.DemandSource = demandSource
		account.Extra = dbtest.ValidDemandSourceAccountExtra(t, adapter.Key(demandSource.APIKey))
	})

	app := dbtest.CreateApp(t, tx)

	dbtest.CreateAppDemandProfile(t, tx, func(profile *db.AppDemandProfile) {
		profile.App = app
		profile.Account = account
		profile.Data = dbtest.ValidAppDemandProfileData(t, adapter.Key(demandSource.APIKey), app.ID)
	})

	dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = account
		item.AdType = db.BannerAdType
		item.Format = sql.NullString{
			String: string(ad.BannerFormat),
			Valid:  true,
		}
		item.Extra = map[string]any{"slot_uuid": "amazon_slot_banner", "format": "BANNER"}
	})

	dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = account
		item.AdType = db.BannerAdType
		item.Format = sql.NullString{
			String: string(ad.MRECFormat),
			Valid:  true,
		}
		item.Extra = map[string]any{"slot_uuid": "amazon_slot_mrec", "format": "MREC"}
	})

	dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = account
		item.AdType = db.InterstitialAdType
		item.Extra = map[string]any{"slot_uuid": "amazon_slot_interstitial", "format": "INTERSTITIAL"}
	})

	dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = account
		item.AdType = db.InterstitialAdType
		item.Extra = map[string]any{"slot_uuid": "amazon_slot_video", "format": "VIDEO"}
	})

	dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.App = app
		item.Account = account
		item.AdType = db.RewardedAdType
		item.Extra = map[string]any{"slot_uuid": "amazon_slot_rewarded", "format": "REWARDED"}
	})

	profilesCache := config.NewRedisCacheOf[[]db.AppDemandProfile](rdb, 10*time.Minute, "app_demand_profiles")
	amazonSlotsCache := config.NewRedisCacheOf[[]sdkapi.AmazonSlot](rdb, 10*time.Minute, "amazon_slots")
	fetcher := &AdapterInitConfigsFetcher{DB: tx, ProfilesCache: profilesCache, AmazonSlotsCache: amazonSlotsCache}

	tests := []struct {
		name           string
		appID          int64
		setAmazonSlots bool
		adapterKeys    []adapter.Key
		want           []sdkapi.AdapterInitConfig
	}{
		{
			name:           "set amazon slots",
			appID:          app.ID,
			adapterKeys:    adapter.Keys,
			setAmazonSlots: true,
			want: []sdkapi.AdapterInitConfig{
				&sdkapi.AmazonInitConfig{
					AppKey: fmt.Sprintf("amazon_app_%d", app.ID),
					Slots: []sdkapi.AmazonSlot{
						{
							SlotUUID: "amazon_slot_banner",
							Format:   "BANNER",
						},
						{
							SlotUUID: "amazon_slot_mrec",
							Format:   "MREC",
						},
						{
							SlotUUID: "amazon_slot_interstitial",
							Format:   "INTERSTITIAL",
						},
						{
							SlotUUID: "amazon_slot_video",
							Format:   "VIDEO",
						},
						{
							SlotUUID: "amazon_slot_rewarded",
							Format:   "REWARDED",
						},
					},
				},
			},
		},
		{
			name:           "do not set amazon slots",
			appID:          app.ID,
			adapterKeys:    adapter.Keys,
			setAmazonSlots: false,
			want: []sdkapi.AdapterInitConfig{
				&sdkapi.AmazonInitConfig{
					AppKey: fmt.Sprintf("amazon_app_%d", app.ID),
					Slots:  nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetcher.FetchAdapterInitConfigs(context.Background(), tt.appID, tt.adapterKeys, tt.setAmazonSlots, false)
			if err != nil {
				t.Fatalf("FetchAdapterInitConfigs() error = %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FetchAdapterInitConfigs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
