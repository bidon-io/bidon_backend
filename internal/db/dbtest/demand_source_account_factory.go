package dbtest

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type DemandSourceAccountFactory struct {
	DemandSource func(int) db.DemandSource
	User         func(int) db.User
	Type         func(int) string
	Extra        func(int) []byte
	IsBidding    func(int) bool
	IsDefault    func(int) bool
}

func (f DemandSourceAccountFactory) Build(i int) db.DemandSourceAccount {
	a := db.DemandSourceAccount{}

	var ds db.DemandSource
	if f.DemandSource == nil {
		ds = DemandSourceFactory{}.Build(i)
	} else {
		ds = f.DemandSource(i)
	}
	a.DemandSourceID = ds.ID
	a.DemandSource = ds

	var user db.User
	if f.User == nil {
		user = UserFactory{}.Build(i)
	} else {
		user = f.User(i)
	}
	a.UserID = user.ID
	a.User = user

	if f.Type == nil {
		a.Type = fmt.Sprintf("DemandSourceAccount::%s", ds.APIKey)
	} else {
		a.Type = f.Type(i)
	}

	if f.Extra == nil {
		a.Extra = []byte(`{"foo": "bar"}`)
	} else {
		a.Extra = f.Extra(i)
	}

	if f.IsBidding == nil {
		a.IsBidding = new(bool)
	} else {
		a.IsBidding = new(bool)
		*a.IsBidding = f.IsBidding(i)
	}

	if f.IsDefault == nil {
		a.IsDefault = sql.NullBool{
			Valid: true,
			Bool:  true,
		}
	} else {
		a.IsDefault = sql.NullBool{
			Valid: true,
			Bool:  f.IsDefault(i),
		}
	}

	return a
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
	case adapter.MetaKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.MintegralKey:
		return []byte(`{"publisher_id": "1", "app_key": "mintegral", "foo": "bar"}`)
	case adapter.UnityAdsKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.VungleKey:
		return []byte(`{"account_id": "vungle", "foo": "bar"}`)
	case adapter.MobileFuseKey:
		return []byte(`{"foo": "bar"}`)
	case adapter.InmobiKey:
		return []byte(`{"account_id": "inmobi", "foo": "bar"}`)
	case adapter.AmazonKey:
		return []byte(`{"foo": "bar"}`)
	default:
		t.Fatalf("Invalid adapter key or missing valid ACCOUNT config for adapter %q", key)
		return nil
	}
}
