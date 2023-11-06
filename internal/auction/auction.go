package auction

import (
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type Auction struct {
	ConfigID                 int64         `json:"auction_configuration_id"`
	ConfigUID                string        `json:"auction_configuration_uid"`
	ExternalWinNotifications bool          `json:"external_win_notifications"`
	Rounds                   []RoundConfig `json:"rounds"`
	LineItems                []LineItem    `json:"line_items"` // Deprecated: use AdUnits instead
	AdUnits                  []AdUnit      `json:"ad_units"`
	Segment                  Segment       `json:"segment"`
}
type Config struct {
	ID                       int64
	UID                      string
	ExternalWinNotifications bool
	Rounds                   []RoundConfig
}

type RoundConfig struct {
	ID      string        `json:"id"`
	Demands []adapter.Key `json:"demands"`
	Bidding []adapter.Key `json:"bidding"`
	Timeout int           `json:"timeout"`
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
}

type AdUnit struct {
	DemandID   string         `json:"demand_id"`
	UID        string         `json:"uid"`
	Label      string         `json:"label"`
	PriceFloor *float64       `json:"pricefloor,omitempty"`
	BidType    schema.BidType `json:"bid_type"`
	Extra      map[string]any `json:"ext"`
}

type Segment struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
}
