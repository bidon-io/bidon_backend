package store

import (
	"context"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/device"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"gorm.io/gorm"
)

type AdUnitsMatcher struct {
	DB *db.DB
}

func (m *AdUnitsMatcher) Match(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
	if params.AdType == ad.BannerType && !params.AdFormat.IsBannerFormat() {
		return []auction.AdUnit{}, nil
	}

	query := m.DB.
		WithContext(ctx).
		Select("bid_floor", "line_items.human_name", "line_items.bidding", "line_items.extra", "line_items.public_uid").
		Where(map[string]any{
			"app_id":  params.AppID,
			"ad_type": db.AdTypeFromDomain(params.AdType),
		}).
		InnerJoins("Account", m.DB.Select("id")).
		InnerJoins("Account.DemandSource", m.DB.Select("api_key").Where(map[string]any{"api_key": params.Adapters}))

	if params.PriceFloor != nil {
		query = query.Where("(bid_floor >= ? OR line_items.bidding)", params.PriceFloor)
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

func (m *AdUnitsMatcher) find(query *gorm.DB) ([]auction.AdUnit, error) {
	var dbLineItems []db.LineItem
	if err := query.Find(&dbLineItems).Error; err != nil {
		return nil, err
	}

	adUnits := make([]auction.AdUnit, len(dbLineItems))
	for i := range dbLineItems {
		dbLineItem := &dbLineItems[i]

		adUnits[i].DemandID = dbLineItem.Account.DemandSource.APIKey
		adUnits[i].UID = strconv.FormatInt(dbLineItem.PublicUID.Int64, 10)
		adUnits[i].Label = dbLineItem.HumanName
		isRTB := dbLineItem.IsBidding.Valid && dbLineItem.IsBidding.Bool
		if isRTB {
			adUnits[i].BidType = schema.RTBBidType
		} else {
			adUnits[i].BidType = schema.CPMBidType
		}
		if dbLineItem.BidFloor.Valid && !isRTB {
			pf := dbLineItem.BidFloor.Decimal.InexactFloat64()
			adUnits[i].PriceFloor = &pf
		}
		adUnits[i].Extra = dbLineItem.Extra
	}

	return adUnits, nil
}
