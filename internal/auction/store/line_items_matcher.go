package store

import (
	"context"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/device"
	"gorm.io/gorm"
)

type LineItemsMatcher struct {
	DB *db.DB
}

func (m *LineItemsMatcher) Match(ctx context.Context, params *auction.BuildParams) ([]auction.LineItem, error) {
	if params.AdType == ad.BannerType && !params.AdFormat.IsBannerFormat() {
		return []auction.LineItem{}, nil
	}

	query := m.DB.
		WithContext(ctx).
		Select("bid_floor", "code", "line_items.extra", "line_items.public_uid").
		Where(map[string]any{
			"app_id":  params.AppID,
			"ad_type": db.AdTypeFromDomain(params.AdType),
		}).
		InnerJoins("Account", m.DB.Select("id")).
		InnerJoins("Account.DemandSource", m.DB.Select("api_key").Where(map[string]any{"api_key": params.Adapters}))

	if params.PriceFloor != nil {
		query = query.Where("bid_floor >= ?", params.PriceFloor)
	}

	if params.AdType != ad.BannerType {
		return m.find(query)
	}

	adFormats := []ad.Format{params.AdFormat}
	if params.AdFormat == ad.AdaptiveFormat {
		switch params.DeviceType {
		case device.TabletType:
			adFormats = append(adFormats, ad.LeaderboardFormat)
		case device.PhoneType:
			adFormats = append(adFormats, ad.BannerFormat)
		}
	}

	query = query.Where(map[string]any{"format": adFormats})

	return m.find(query)
}

func (m *LineItemsMatcher) find(query *gorm.DB) ([]auction.LineItem, error) {
	var dbLineItems []db.LineItem
	if err := query.Find(&dbLineItems).Error; err != nil {
		return nil, err
	}

	lineItems := make([]auction.LineItem, len(dbLineItems))
	for i := range dbLineItems {
		dbLineItem := &dbLineItems[i]

		lineItems[i].ID = dbLineItem.Account.DemandSource.APIKey
		lineItems[i].UID = strconv.FormatInt(dbLineItem.PublicUID.Int64, 10)
		lineItems[i].PriceFloor = dbLineItem.BidFloor.Decimal.InexactFloat64()
		lineItems[i].AdUnitID = *dbLineItem.Code

		adUnitID, ok := dbLineItem.Extra["ad_unit_id"].(string)
		if ok {
			lineItems[i].AdUnitID = adUnitID
		}

		unitID, ok := dbLineItem.Extra["unit_id"].(string)
		if ok {
			lineItems[i].AdUnitID = unitID
		}

		placementID, ok := dbLineItem.Extra["placement_id"].(string)
		if ok {
			lineItems[i].PlacementID = placementID
		}

		spotID, ok := dbLineItem.Extra["spot_id"].(string)
		if ok {
			lineItems[i].PlacementID = spotID
		}

		zoneID, ok := dbLineItem.Extra["zone_id"].(string)
		if ok {
			lineItems[i].ZonedID = zoneID
		}

		slotUUID, ok := dbLineItem.Extra["slot_uuid"].(string)
		if ok {
			lineItems[i].SlotUUID = slotUUID
		}

		slotID, ok := dbLineItem.Extra["slot_id"].(string)
		if ok {
			lineItems[i].SlotID = slotID
		}
	}

	return lineItems, nil
}
