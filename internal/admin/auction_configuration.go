package admin

import "github.com/bidon-io/bidon-backend/internal/ad"

type AuctionConfiguration struct {
	ID int64 `json:"id"`
	AuctionConfigurationAttrs
}

// AuctionConfigurationAttrs is attributes of Configuration. Used to create and update configurations
type AuctionConfigurationAttrs struct {
	Name       string                      `json:"name"`
	AppID      int64                       `json:"app_id"`
	AdType     ad.Type                     `json:"ad_type"`
	Rounds     []AuctionRoundConfiguration `json:"rounds"`
	Pricefloor float64                     `json:"pricefloor"`
}

type AuctionRoundConfiguration struct {
	ID      string   `json:"id"`
	Demands []string `json:"demands"`
	Timeout int      `json:"timeout"`
}

type AuctionConfigurationService = resourceService[AuctionConfiguration, AuctionConfigurationAttrs]
