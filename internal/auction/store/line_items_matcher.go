package store

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"sort"
	"strconv"

	"gorm.io/gorm"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/device"
)

type LineItemsMatcher struct {
	DB    *db.DB
	Cache cache[[]auction.LineItem]
}

func (m *LineItemsMatcher) MatchCached(ctx context.Context, params *auction.BuildParams) ([]auction.LineItem, error) {
	key, err := m.cacheKey(*params)
	if err != nil {
		return nil, err
	}
	return m.Cache.Get(ctx, key, func(ctx context.Context) ([]auction.LineItem, error) {
		return m.Match(ctx, params)
	})
}

func (m *LineItemsMatcher) Match(ctx context.Context, params *auction.BuildParams) ([]auction.LineItem, error) {
	if params.AdType == ad.BannerType && !params.AdFormat.IsBannerFormat() {
		return []auction.LineItem{}, nil
	}

	query := m.DB.
		WithContext(ctx).
		Select("bid_floor", "line_items.extra", "line_items.public_uid").
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

	adFormats := m.selectAdFormats(params)
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
			setAdUnitIDIfEmpty(&lineItems[i], placementID)
		}

		spotID, ok := dbLineItem.Extra["spot_id"].(string)
		if ok {
			lineItems[i].PlacementID = spotID
			setAdUnitIDIfEmpty(&lineItems[i], spotID)
		}

		zoneID, ok := dbLineItem.Extra["zone_id"].(string)
		if ok {
			lineItems[i].ZonedID = zoneID
			setAdUnitIDIfEmpty(&lineItems[i], zoneID)
		}

		slotUUID, ok := dbLineItem.Extra["slot_uuid"].(string)
		if ok {
			lineItems[i].SlotUUID = slotUUID
			setAdUnitIDIfEmpty(&lineItems[i], slotUUID)
		}

		slotID, ok := dbLineItem.Extra["slot_id"].(string)
		if ok {
			lineItems[i].SlotID = slotID
			setAdUnitIDIfEmpty(&lineItems[i], slotID)
		}

		mediation, ok := dbLineItem.Extra["mediation"].(string)
		if ok {
			lineItems[i].Mediation = mediation
		}

	}

	return lineItems, nil
}

func setAdUnitIDIfEmpty(lineItem *auction.LineItem, value string) {
	if lineItem.AdUnitID == "" {
		lineItem.AdUnitID = value
	}
}

func (m *LineItemsMatcher) cacheKey(params auction.BuildParams) ([]byte, error) {
	// Sort adapter keys to get deterministic cache key
	sort.Slice(params.Adapters, func(i, j int) bool {
		return params.Adapters[i] < params.Adapters[j]
	})
	jsonData, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(jsonData)
	return hash[:], nil
}

func (m *LineItemsMatcher) selectAdFormats(params *auction.BuildParams) []ad.Format {
	adFormats := []ad.Format{params.AdFormat}
	switch params.AdFormat {
	case ad.AdaptiveFormat:
		switch params.DeviceType {
		case device.TabletType:
			adFormats = append(adFormats, ad.LeaderboardFormat)
		case device.PhoneType:
			adFormats = append(adFormats, ad.BannerFormat)
		}
	case ad.BannerFormat, ad.LeaderboardFormat:
		adFormats = append(adFormats, ad.AdaptiveFormat)
	}

	return adFormats
}
