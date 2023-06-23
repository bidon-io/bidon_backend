package auction

import (
	"github.com/bidon-io/bidon-backend/internal/adapter"
)

type Auction struct {
	ConfigID                 int64         `json:"auction_configuration_id"`
	ExternalWinNotifications bool          `json:"external_win_notifications"`
	Rounds                   []RoundConfig `json:"rounds"`
	LineItems                []LineItem    `json:"line_items"`
}
type Config struct {
	ID                       int64
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
	ID         string  `json:"id"`
	PriceFloor float64 `json:"pricefloor"`
	AdUnitID   string  `json:"ad_unit_id"`
}
