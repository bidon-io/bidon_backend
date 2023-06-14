package store

import (
	"context"

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
		Select("bid_floor", "code").
		Where(map[string]any{
			"app_id":  params.AppID,
			"ad_type": db.AdTypeFromDomain(params.AdType),
		}).
		InnerJoins("Account", m.DB.Select("id")).
		InnerJoins("Account.DemandSource", m.DB.Select("api_key").Where(map[string]any{"api_key": params.Adapters}))

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
		lineItems[i].PriceFloor = dbLineItem.BidFloor.Decimal.InexactFloat64()
		lineItems[i].AdUnitID = *dbLineItem.Code
	}

	return lineItems, nil
}
