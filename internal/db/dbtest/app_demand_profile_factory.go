package dbtest

import (
	"fmt"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type AppDemandProfileFactory struct {
	App          func(int) db.App
	Account      func(int) db.DemandSourceAccount
	DemandSource func(int) db.DemandSource
	Data         func(int) []byte
}

func (f AppDemandProfileFactory) Build(i int) db.AppDemandProfile {
	profile := db.AppDemandProfile{}

	var app db.App
	if f.App == nil {
		app = AppFactory{}.Build(i)
	} else {
		app = f.App(i)
	}
	profile.AppID = app.ID
	profile.App = app

	var demandSource db.DemandSource
	if f.DemandSource == nil {
		demandSource = DemandSourceFactory{}.Build(i)
	} else {
		demandSource = f.DemandSource(i)
	}
	profile.DemandSourceID = demandSource.ID
	profile.DemandSource = demandSource

	var account db.DemandSourceAccount
	if f.Account == nil {
		account = DemandSourceAccountFactory{
			DemandSource: func(i int) db.DemandSource {
				return demandSource
			},
			User: func(i int) db.User {
				return app.User
			},
		}.Build(i)
	} else {
		account = f.Account(i)
	}
	profile.AccountID = account.ID
	profile.AccountType = account.Type
	profile.Account = account

	if f.Data == nil {
		profile.Data = []byte(`{"foo": "bar"}`)
	} else {
		profile.Data = f.Data(i)
	}

	return profile
}

func ValidAppDemandProfileData(t *testing.T, key adapter.Key, appID int64) []byte {
	t.Helper()

	switch key {
	case adapter.AdmobKey:
		return []byte(fmt.Sprintf(`{"app_id": "admob_app_%d", "foo": "bar"}`, appID))
	case adapter.ApplovinKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.BidmachineKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.BigoAdsKey:
		return []byte(fmt.Sprintf(`{"app_id": "bigoads_app_%d", "foo": "bar"}`, appID))
	case adapter.DTExchangeKey:
		return []byte(fmt.Sprintf(`{"app_id": "dtexchange_app_%d", "foo": "bar"}`, appID))
	case adapter.MetaKey:
		return []byte(fmt.Sprintf(`{"app_id": "meta_app_%d", "app_secret": "meta_app_%d_secret", "foo": "bar"}`, appID, appID))
	case adapter.MintegralKey:
		return []byte(fmt.Sprintf(`{"app_id": "mintegral_app_%d", "foo": "bar"}`, appID))
	case adapter.UnityAdsKey:
		return []byte(fmt.Sprintf(`{"game_id": "unityads_game_%d", "foo": "bar"}`, appID))
	case adapter.VungleKey:
		return []byte(fmt.Sprintf(`{"app_id": "vungle_app_%d", "foo": "bar"}`, appID))
	case adapter.MobileFuseKey:
		return []byte(`{"foo": "bar"}`)
	default:
		t.Fatalf("Invalid adapter key or missing valid APP config for adapter %q", key)
		return nil
	}
}
