package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/ad"
	"os"
	"testing"

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

	user := dbtest.CreateUser(t, tx, 1)
	app := &db.App{
		UserID:      user.ID,
		AppKey:      sql.NullString{String: "asdf", Valid: true},
		PackageName: sql.NullString{String: "com.example.app", Valid: true},
	}
	if err := tx.Create(app).Error; err != nil {
		t.Fatalf("Error creating app: %v", err)
	}

	fetcher := &AppFetcher{DB: tx}

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
		app, err := fetcher.Fetch(context.Background(), tC.appKey, tC.appBundle)

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

	keys := adapter.Keys

	demandSources := dbtest.CreateList[db.DemandSource](t, tx,
		dbtest.DemandSourceFactory{
			APIKey: func(i int) string { return string(keys[i]) },
		},
		len(keys),
	)

	accounts := dbtest.CreateList[db.DemandSourceAccount](t, tx,
		dbtest.DemandSourceAccountFactory{
			DemandSource: func(i int) db.DemandSource { return demandSources[i] },
			Extra: func(i int) []byte {
				demandSource := demandSources[i]
				return dbtest.ValidDemandSourceAccountExtra(t, adapter.Key(demandSource.APIKey))
			},
		},
		len(demandSources),
	)

	apps := dbtest.CreateList[db.App](t, tx,
		dbtest.AppFactory{},
		2,
	)

	selectApp := func(i int) db.App {
		if i < (len(accounts) / 2) {
			return apps[0]
		} else {
			return apps[1]
		}
	}
	_ = dbtest.CreateList[db.AppDemandProfile](t, tx,
		dbtest.AppDemandProfileFactory{
			App: func(i int) db.App {
				return selectApp(i)
			},
			Account: func(i int) db.DemandSourceAccount {
				return accounts[i]
			},
			Data: func(i int) []byte {
				demandSource := accounts[i].DemandSource
				app := selectApp(i)

				return dbtest.ValidAppDemandProfileData(t, adapter.Key(demandSource.APIKey), app.ID)
			},
		},
		len(accounts),
	)

	fetcher := &AdapterInitConfigsFetcher{DB: tx}

	tests := []struct {
		name        string
		appID       int64
		adapterKeys []adapter.Key
		want        []sdkapi.AdapterInitConfig
	}{
		{
			name:        "first app with all adapters",
			appID:       apps[0].ID,
			adapterKeys: adapter.Keys,
			want: []sdkapi.AdapterInitConfig{
				&sdkapi.AdmobInitConfig{
					AppID: fmt.Sprintf("admob_app_%d", apps[0].ID),
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
				&sdkapi.DTExchangeInitConfig{
					AppID: fmt.Sprintf("dtexchange_app_%d", apps[0].ID),
				},
				&sdkapi.MetaInitConfig{
					AppID:     fmt.Sprintf("meta_app_%d", apps[0].ID),
					AppSecret: fmt.Sprintf("meta_app_%d_secret", apps[0].ID),
				},
				&sdkapi.MintegralInitConfig{
					AppID:  fmt.Sprintf("mintegral_app_%d", apps[0].ID),
					AppKey: "mintegral",
				},
			},
		},
		{
			name:        "second app with all adapters",
			appID:       apps[1].ID,
			adapterKeys: adapter.Keys,
			want: []sdkapi.AdapterInitConfig{
				&sdkapi.MobileFuseInitConfig{},
				&sdkapi.UnityAdsInitConfig{
					GameID: fmt.Sprintf("unityads_game_%d", apps[1].ID),
				},
				&sdkapi.VungleInitConfig{
					AppID: fmt.Sprintf("vungle_app_%d", apps[1].ID),
				},
				&sdkapi.BigoAdsInitConfig{
					AppID: fmt.Sprintf("bigoads_app_%d", apps[1].ID),
				},
				&sdkapi.InmobiInitConfig{
					AccountID: "inmobi",
					AppKey:    fmt.Sprintf("inmobi_app_%d", apps[1].ID),
				},
				&sdkapi.AmazonInitConfig{
					AppKey: fmt.Sprintf("amazon_app_%d", apps[1].ID),
					Slots:  []sdkapi.AmazonSlot{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetcher.FetchAdapterInitConfigs(context.Background(), tt.appID, tt.adapterKeys)
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

	demandSource := dbtest.Create[db.DemandSource](t, tx,
		dbtest.DemandSourceFactory{
			APIKey: func(i int) string { return string(adapter.AmazonKey) },
		},
		0,
	)

	account := dbtest.Create[db.DemandSourceAccount](t, tx,
		dbtest.DemandSourceAccountFactory{
			DemandSource: func(i int) db.DemandSource { return demandSource },
			Extra: func(i int) []byte {
				demandSource := demandSource
				return dbtest.ValidDemandSourceAccountExtra(t, adapter.Key(demandSource.APIKey))
			},
		},
		0,
	)

	app := dbtest.Create[db.App](t, tx,
		dbtest.AppFactory{},
		1,
	)

	_ = dbtest.Create[db.AppDemandProfile](t, tx,
		dbtest.AppDemandProfileFactory{
			App: func(i int) db.App {
				return app
			},
			Account: func(i int) db.DemandSourceAccount {
				return account
			},
			Data: func(i int) []byte {
				return dbtest.ValidAppDemandProfileData(t, adapter.Key(demandSource.APIKey), app.ID)
			},
		},
		0,
	)

	dbtest.Create[db.LineItem](t, tx,
		dbtest.LineItemFactory{
			App:     func(i int) db.App { return app },
			Account: func(i int) db.DemandSourceAccount { return account },
			AdType:  func(i int) db.AdType { return db.BannerAdType },
			Extra: func(i int) map[string]any {
				return map[string]any{"slot_uuid": "amazon_slot_banner"}
			},
		}, 0)

	dbtest.Create[db.LineItem](t, tx,
		dbtest.LineItemFactory{
			App:     func(i int) db.App { return app },
			Account: func(i int) db.DemandSourceAccount { return account },
			AdType:  func(i int) db.AdType { return db.BannerAdType },
			Format: func(i int) sql.NullString {
				return sql.NullString{
					String: string(ad.MRECFormat),
					Valid:  true,
				}
			},
			Extra: func(i int) map[string]any {
				return map[string]any{"slot_uuid": "amazon_slot_mrec"}
			},
		}, 0)

	dbtest.Create[db.LineItem](t, tx,
		dbtest.LineItemFactory{
			App:     func(i int) db.App { return app },
			Account: func(i int) db.DemandSourceAccount { return account },
			AdType:  func(i int) db.AdType { return db.InterstitialAdType },
			Extra: func(i int) map[string]any {
				return map[string]any{"slot_uuid": "amazon_slot_interstitial"}
			},
		}, 0)

	dbtest.Create[db.LineItem](t, tx,
		dbtest.LineItemFactory{
			App:     func(i int) db.App { return app },
			Account: func(i int) db.DemandSourceAccount { return account },
			AdType:  func(i int) db.AdType { return db.InterstitialAdType },
			Extra: func(i int) map[string]any {
				return map[string]any{"slot_uuid": "amazon_slot_video", "is_video": true}
			},
		}, 0)

	fetcher := &AdapterInitConfigsFetcher{DB: tx}

	tests := []struct {
		name        string
		appID       int64
		adapterKeys []adapter.Key
		want        []sdkapi.AdapterInitConfig
	}{
		{
			name:        "amazon app with all line items",
			appID:       app.ID,
			adapterKeys: adapter.Keys,
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
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetcher.FetchAdapterInitConfigs(context.Background(), tt.appID, tt.adapterKeys)
			if err != nil {
				t.Fatalf("FetchAdapterInitConfigs() error = %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FetchAdapterInitConfigs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
