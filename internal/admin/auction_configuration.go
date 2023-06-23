package admin

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
)

type AuctionConfiguration struct {
	ID int64 `json:"id"`
	AuctionConfigurationAttrs
}

// AuctionConfigurationAttrs is attributes of Configuration. Used to create and update configurations
type AuctionConfigurationAttrs struct {
	Name                     string                `json:"name"`
	AppID                    int64                 `json:"app_id"`
	AdType                   ad.Type               `json:"ad_type"`
	Rounds                   []auction.RoundConfig `json:"rounds"`
	Pricefloor               float64               `json:"pricefloor"`
	SegmentID                *int64                `json:"segment_id"`
	ExternalWinNotifications *bool                 `json:"external_win_notifications"`
}

type AuctionConfigurationService = resourceService[AuctionConfiguration, AuctionConfigurationAttrs]
