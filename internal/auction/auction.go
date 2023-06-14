package auction

type Auction struct {
	ConfigID  int64         `json:"auction_configuration_id"`
	Rounds    []RoundConfig `json:"rounds"`
	LineItems []LineItem    `json:"line_items"`
}
type Config struct {
	ID     int64
	Rounds []RoundConfig
}

type RoundConfig struct {
	ID      string   `json:"id"`
	Demands []string `json:"demands"`
	Timeout int      `json:"timeout"`
}

type LineItem struct {
	ID         string  `json:"id"`
	PriceFloor float64 `json:"pricefloor"`
	AdUnitID   string  `json:"ad_unit_id"`
}
