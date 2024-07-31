package dbtest

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
)

func demandSourceAccountDefaults(n uint32) func(*db.DemandSourceAccount) {
	return func(account *db.DemandSourceAccount) {
		if account.DemandSourceID == 0 && account.DemandSource.ID == 0 {
			account.DemandSource = BuildDemandSource(func(source *db.DemandSource) {
				*source = account.DemandSource
			})
		}
		if account.UserID == 0 && account.User.ID == 0 {
			account.User = BuildUser(func(user *db.User) {
				*user = account.User
			})
		}
		if account.Label == (sql.NullString{}) {
			account.Label = sql.NullString{
				String: fmt.Sprintf("Test Account %d", n),
				Valid:  true,
			}
		}
		if account.Type == "" {
			account.Type = fmt.Sprintf("DemandSourceAccount::%s", account.DemandSource.APIKey)
		}
		if account.Extra == nil {
			account.Extra = []byte(`{"foo": "bar"}`)
		}
		if account.IsBidding == (sql.NullBool{}) {
			account.IsBidding = sql.NullBool{
				Bool:  false,
				Valid: true,
			}
		}
		if account.IsDefault == (sql.NullBool{}) {
			account.IsDefault = sql.NullBool{
				Valid: true,
				Bool:  true,
			}
		}
	}
}

func BuildDemandSourceAccount(opts ...func(*db.DemandSourceAccount)) db.DemandSourceAccount {
	var account db.DemandSourceAccount

	n := counter.get("demand_source_account")

	opts = append(opts, demandSourceAccountDefaults(n))
	for _, opt := range opts {
		opt(&account)
	}

	return account
}

func CreateDemandSourceAccount(t *testing.T, tx *db.DB, opts ...func(*db.DemandSourceAccount)) db.DemandSourceAccount {
	t.Helper()

	account := BuildDemandSourceAccount(opts...)
	if err := tx.Create(&account).Error; err != nil {
		t.Fatalf("Failed to create demand source account: %v", err)
	}

	return account
}

func ValidDemandSourceAccountExtra(t *testing.T, key adapter.Key) []byte {
	t.Helper()

	switch key {
	case adapter.AdmobKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.ApplovinKey:
		return []byte(`{"sdk_key": "applovin", "foo": "bar"}`)
	case adapter.BidmachineKey:
		return []byte(`{"seller_id": "1", "endpoint": "x.appbaqend.com", "mediation_config": ["one", "two"], "foo": "bar"}`)
	case adapter.BigoAdsKey:
		return []byte(`{"publisher_id": "1", "endpoint": "https://example.com", "foo": "bar"}`)
	case adapter.DTExchangeKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.GAMKey:
		return []byte(`{"network_code": "111", "foo": "bar"}`)
	case adapter.MetaKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.MintegralKey:
		return []byte(`{"publisher_id": "1", "app_key": "mintegral", "foo": "bar"}`)
	case adapter.UnityAdsKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.VungleKey:
		return []byte(`{"account_id": "vungle", "foo": "bar"}`)
	case adapter.VKAdsKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.MobileFuseKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.InmobiKey:
		return []byte(`{"account_id": "inmobi", "foo": "bar"}`)
	case adapter.AmazonKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.YandexKey:
		return []byte(`{"oauth_token": "yandex"}`)
	case adapter.IronSourceKey:
		return []byte(`{"app_key": "ironsource"}`)
	default:
		t.Fatalf("Invalid adapter key or missing valid ACCOUNT config for adapter %q", key)
		return nil
	}
}
