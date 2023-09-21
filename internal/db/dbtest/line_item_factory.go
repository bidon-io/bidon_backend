package dbtest

import (
	"database/sql"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/shopspring/decimal"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type LineItemFactory struct {
	App       func(int) db.App
	Account   func(int) db.DemandSourceAccount
	HumanName func(int) string
	Code      func(int) string
	BidFloor  func() decimal.NullDecimal
	AdType    func(int) db.AdType
	Format    func(int) sql.NullString
	Extra     func(int) map[string]any
	Width     func(int) int32
	Height    func(int) int32
	PublicUID func(int) sql.NullInt64
}

func (f LineItemFactory) Build(i int) db.LineItem {
	li := db.LineItem{}

	var app db.App
	if f.App == nil {
		app = AppFactory{}.Build(i)
	} else {
		app = f.App(i)
	}
	li.AppID = app.ID
	li.App = app

	var account db.DemandSourceAccount
	if f.Account == nil {
		account = DemandSourceAccountFactory{}.Build(i)
	} else {
		account = f.Account(i)
	}
	li.AccountID = account.ID
	li.Account = account

	if f.HumanName == nil {
		li.HumanName = fmt.Sprintf("Test Line Item %d", i)
	} else {
		li.HumanName = f.HumanName(i)
	}

	var code string
	if f.Code == nil {
		code = fmt.Sprintf("code%d", i)
	} else {
		code = f.Code(i)
	}
	li.Code = &code

	if f.BidFloor == nil {
		li.BidFloor = decimal.NewNullDecimal(decimal.RequireFromString("0.1"))
	} else {
		li.BidFloor = f.BidFloor()
	}

	if f.AdType == nil {
		li.AdType = db.BannerAdType
	} else {
		li.AdType = f.AdType(i)
	}

	if f.Extra == nil {
		li.Extra = map[string]any{
			"foo": "bar",
		}
	} else {
		li.Extra = f.Extra(i)
	}

	if f.Format == nil {
		li.Format = sql.NullString{
			String: string(ad.BannerFormat),
			Valid:  true,
		}
	} else {
		li.Format = f.Format(i)
	}

	if f.Width == nil {
		li.Width = 320
	} else {
		li.Width = f.Width(i)
	}

	if f.Height == nil {
		li.Height = 50
	} else {
		li.Height = f.Height(i)
	}

	if f.PublicUID == nil {
		li.PublicUID = sql.NullInt64{
			Int64: int64(i),
			Valid: true,
		}
	} else {
		li.PublicUID = f.PublicUID(i)
	}

	return li
}

func ValidLineItemExtra(t *testing.T, key adapter.Key) map[string]any {
	t.Helper()

	switch key {
	case adapter.AdmobKey:
		return map[string]any{
			"ad_unit_id": "admob_line_item",
		}
	case adapter.ApplovinKey:
		return map[string]any{
			"zone_id": "applovin_line_item_zone_id",
		}
	case adapter.BidmachineKey:
		return map[string]any{}
	case adapter.BigoAdsKey:
		return map[string]any{
			"slot_id": "bigo_ads_line_item_slot_id",
		}
	case adapter.DTExchangeKey:
		return map[string]any{
			"placement_id": "dt_exchange_line_item_placement_id",
		}
	case adapter.MetaKey:
		return map[string]any{
			"placement_id": "meta_line_item_placement_id",
		}
	case adapter.MintegralKey:
		return map[string]any{
			"placement_id": "mintegral_line_item_placement_id",
			"ad_unit_id":   "mintegral_line_item_ad_unit_id",
		}
	case adapter.UnityAdsKey:
		return map[string]any{
			"placement_id": "unity_ads_line_item_placement_id",
		}
	case adapter.VungleKey:
		return map[string]any{
			"placement_id": "vungle_line_item_placement_id",
		}
	case adapter.MobileFuseKey:
		return map[string]any{
			"placement_id": "mobile_fuse_line_item_placement_id",
		}
	case adapter.InmobiKey:
		return map[string]any{
			"placement_id": "inmobi_line_item_placement_id",
		}
	default:
		t.Fatalf("Invalid adapter key or missing valid ACCOUNT config for adapter %q", key)
		return nil
	}
}
