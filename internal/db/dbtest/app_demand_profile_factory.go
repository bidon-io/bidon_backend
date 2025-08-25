package dbtest

import (
	"fmt"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
)

func appDemandProfileDefaults(n uint32) func(*db.AppDemandProfile) {
	return func(profile *db.AppDemandProfile) {
		if profile.AppID == 0 && profile.App.ID == 0 {
			profile.App = BuildApp(func(app *db.App) {
				*app = profile.App
			})
		}

		if profile.AccountID == 0 && profile.Account.ID == 0 {
			profile.Account = BuildDemandSourceAccount(func(account *db.DemandSourceAccount) {
				*account = profile.Account
			})
		}
		profile.AccountType = profile.Account.Type

		if profile.Account.DemandSourceID != 0 {
			profile.DemandSourceID = profile.Account.DemandSourceID
		} else if profile.Account.DemandSource.ID != 0 {
			profile.DemandSourceID = profile.Account.DemandSource.ID
		} else if profile.DemandSourceID == 0 && profile.DemandSource.ID == 0 {
			profile.DemandSource = BuildDemandSource(func(source *db.DemandSource) {
				*source = profile.DemandSource
			})
		}

		if profile.Data == nil {
			profile.Data = []byte(fmt.Sprintf(`{"profile_num": %d, "foo": "bar"}`, n))
		}
	}
}

func BuildAppDemandProfile(opts ...func(*db.AppDemandProfile)) db.AppDemandProfile {
	var profile db.AppDemandProfile

	n := counter.get("app_demand_profile")

	opts = append(opts, appDemandProfileDefaults(n))
	for _, opt := range opts {
		opt(&profile)
	}

	return profile
}

func CreateAppDemandProfile(t *testing.T, tx *db.DB, opts ...func(*db.AppDemandProfile)) db.AppDemandProfile {
	t.Helper()

	profile := BuildAppDemandProfile(opts...)
	if err := tx.Create(&profile).Error; err != nil {
		t.Fatalf("Failed to create app demand profile: %v", err)
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
	case adapter.ChartboostKey:
		return []byte(fmt.Sprintf(`{"app_id": "chartboost_app_%d", "app_signature": "123"}`, appID))
	case adapter.DTExchangeKey:
		return []byte(fmt.Sprintf(`{"app_id": "dtexchange_app_%d", "foo": "bar"}`, appID))
	case adapter.GAMKey:
		return []byte(fmt.Sprintf(`{"app_id": "gam_app_%d", "foo": "bar"}`, appID))
	case adapter.MetaKey:
		return []byte(fmt.Sprintf(`{"app_id": "meta_app_%d", "app_secret": "meta_app_%d_secret", "foo": "bar"}`, appID, appID))
	case adapter.MintegralKey:
		return []byte(fmt.Sprintf(`{"app_id": "mintegral_app_%d", "foo": "bar"}`, appID))
	case adapter.MolocoKey:
		return []byte(fmt.Sprintf(`{"app_key": "moloco_app_%d"}`, appID))
	case adapter.UnityAdsKey:
		return []byte(fmt.Sprintf(`{"game_id": "unityads_game_%d", "foo": "bar"}`, appID))
	case adapter.VKAdsKey:
		return []byte(fmt.Sprintf(`{"app_id": "vkads_app_%d", "foo": "bar"}`, appID))
	case adapter.VungleKey:
		return []byte(fmt.Sprintf(`{"app_id": "vungle_app_%d", "foo": "bar"}`, appID))
	case adapter.MobileFuseKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.InmobiKey:
		return []byte(fmt.Sprintf(`{"app_key": "inmobi_app_%d", "foo": "bar"}`, appID))
	case adapter.AmazonKey:
		return []byte(fmt.Sprintf(`{"app_key": "amazon_app_%d"}`, appID))
	case adapter.YandexKey:
		return []byte(fmt.Sprintf(`{"metrica_id": "yandex_metrica_%d"}`, appID))
	case adapter.IronSourceKey:
		return []byte(fmt.Sprintf(`{"app_key": "ironsource_app_%d"}`, appID))
	default:
		t.Fatalf("Invalid adapter key or missing valid APP config for adapter %q", key)
		return nil
	}
}
