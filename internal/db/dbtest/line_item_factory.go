package dbtest

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/shopspring/decimal"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/db"
)

func lineItemDefaults(n uint32) func(*db.LineItem) {
	return func(item *db.LineItem) {
		if item.AppID == 0 && item.App.ID == 0 {
			item.App = BuildApp(func(app *db.App) {
				*app = item.App
			})
		}
		if item.AccountID == 0 && item.Account.ID == 0 {
			item.Account = BuildDemandSourceAccount(func(account *db.DemandSourceAccount) {
				*account = item.Account
			})
		}
		if item.HumanName == "" {
			item.HumanName = fmt.Sprintf("Test Line Item %d", n)
		}
		if item.Code == nil {
			code := fmt.Sprintf("code%d", n)
			item.Code = &code
		}
		if item.BidFloor == (decimal.NullDecimal{}) {
			item.BidFloor = decimal.NewNullDecimal(decimal.RequireFromString("0.1"))
		}
		if item.AdType == 0 {
			item.AdType = db.BannerAdType
		}
		if item.Extra == nil {
			item.Extra = map[string]any{
				"foo": "bar",
			}
		}
		if item.Format == (sql.NullString{}) {
			item.Format = sql.NullString{
				String: string(ad.BannerFormat),
				Valid:  true,
			}
		}
		if item.Width == 0 {
			item.Width = 320
		}
		if item.Height == 0 {
			item.Height = 50
		}
		if item.PublicUID == (sql.NullInt64{}) {
			item.PublicUID = sql.NullInt64{
				Int64: int64(n),
				Valid: true,
			}
		}
	}
}

func BuildLineItem(opts ...func(*db.LineItem)) db.LineItem {
	var item db.LineItem

	n := counter.get("line_item")

	opts = append(opts, lineItemDefaults(n))
	for _, opt := range opts {
		opt(&item)
	}

	return item
}

func CreateLineItem(t *testing.T, tx *db.DB, opts ...func(*db.LineItem)) db.LineItem {
	t.Helper()

	item := BuildLineItem(opts...)
	if err := tx.Create(&item).Error; err != nil {
		t.Fatalf("Failed to create line item: %v", err)
	}

	return item
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
			"spot_id": "dt_exchange_line_item_spot_id",
		}
	case adapter.MetaKey:
		return map[string]any{
			"placement_id": "meta_line_item_placement_id",
		}
	case adapter.MintegralKey:
		return map[string]any{
			"placement_id": "mintegral_line_item_placement_id",
			"unit_id":      "mintegral_line_item_unit_id",
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
	case adapter.AmazonKey:
		return map[string]any{
			"placement_id": "amazon_line_item_placement_id",
		}
	default:
		t.Fatalf("Invalid adapter key or missing valid ACCOUNT config for adapter %q", key)
		return nil
	}
}
