package auction

import (
	"errors"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type Auction struct {
	ConfigID                 int64    `json:"auction_configuration_id"`
	ConfigUID                string   `json:"auction_configuration_uid"`
	ExternalWinNotifications bool     `json:"external_win_notifications"`
	AdUnits                  []AdUnit `json:"ad_units"`
	Segment                  Segment  `json:"segment"`
}
type Config struct {
	ID                       int64
	UID                      string
	ExternalWinNotifications bool
	Bidding                  []adapter.Key `json:"bidding"`
	Demands                  []adapter.Key `json:"demands"`
	AdUnitIDs                []int64       `json:"ad_unit_ids"`
	Timeout                  int           `json:"timeout"`
	PriceFloor               float64       `json:"pricefloor"`
}

type LineItem struct {
	ID          string  `json:"id"`
	UID         string  `json:"uid"`
	PriceFloor  float64 `json:"pricefloor"`
	AdUnitID    string  `json:"ad_unit_id"`
	PlacementID string  `json:"placement_id"`
	ZonedID     string  `json:"zoned_id"`
	SlotUUID    string  `json:"slot_uuid"`
	SlotID      string  `json:"slot_id"`
	Mediation   string  `json:"mediation"`
}

type AdUnit struct {
	DemandID   string         `json:"demand_id"`
	UID        string         `json:"uid"`
	Label      string         `json:"label"`
	PriceFloor *float64       `json:"pricefloor,omitempty"`
	BidType    schema.BidType `json:"bid_type"`
	Timeout    int32          `json:"timeout"`
	Extra      map[string]any `json:"ext"`
}

func (a *AdUnit) GetPriceFloor() float64 {
	if a.PriceFloor == nil {
		return 0
	}
	return *a.PriceFloor
}

func (a *AdUnit) IsCPM() bool {
	return a.BidType == schema.CPMBidType
}

func (a *AdUnit) IsRTB() bool {
	return a.BidType == schema.RTBBidType
}

var ErrNoAdUnitsFound = errors.New("no ad units found")

type AdUnitsMap map[adapter.Key][]AdUnit

func (m *AdUnitsMap) First(key adapter.Key, bidType schema.BidType) (*AdUnit, error) {
	if adUnits, ok := (*m)[key]; ok {
		for _, adUnit := range adUnits {
			if adUnit.BidType == bidType {
				return &adUnit, nil
			}
		}
	}
	return nil, ErrNoAdUnitsFound
}

func (m *AdUnitsMap) All(key adapter.Key, bidType schema.BidType) ([]AdUnit, error) {
	if adUnits, ok := (*m)[key]; ok {
		var result []AdUnit
		for _, adUnit := range adUnits {
			if adUnit.BidType == bidType {
				result = append(result, adUnit)
			}
		}
		return result, nil
	}
	return nil, ErrNoAdUnitsFound
}

func BuildAdUnitsMap(adUnits *[]AdUnit) *AdUnitsMap {
	m := make(AdUnitsMap)
	if adUnits == nil {
		return &m
	}

	for _, adUnit := range *adUnits {
		key := adapter.Key(adUnit.DemandID)
		m[key] = append(m[key], adUnit)
	}

	return &m
}

type Segment struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
}
