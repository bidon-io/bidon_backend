package store

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"sort"
	"strconv"

	"gorm.io/gorm"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/device"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

const AdUnitTimeout = 6_000

type AdUnitsMatcher struct {
	DB    *db.DB
	Cache cache[[]auction.AdUnit]
}

func (m *AdUnitsMatcher) MatchCached(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error) {
	key, err := m.cacheKey(*params)
	if err != nil {
		return nil, err
	}
	return m.Cache.Get(ctx, key, func(ctx context.Context) ([]auction.AdUnit, error) {
		return m.Match(ctx, params)
	})
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

	if len(params.AdUnitIDs) > 0 {
		query = query.Where("line_items.id IN ?", params.AdUnitIDs)
	}

	if params.AdType != ad.BannerType {
		return m.find(query)
	}

	adFormats := m.selectAdFormats(params)
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
		adUnits[i].Timeout = m.timeout(adUnits[i].DemandID)
	}

	return adUnits, nil
}

func (m *AdUnitsMatcher) cacheKey(params auction.BuildParams) ([]byte, error) {
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

func (m *AdUnitsMatcher) selectAdFormats(params *auction.BuildParams) []ad.Format {
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

// timeout 6 seconds for all adapters except admob, which is 10 seconds
func (m *AdUnitsMatcher) timeout(demandID string) int32 {
	if demandID == string(adapter.AdmobKey) {
		return 10_000
	}

	return AdUnitTimeout
}
