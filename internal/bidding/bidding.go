package bidding

import (
	"github.com/bidon-io/bidon-backend/internal/auction"
)

type Bidding struct {
	ConfigID                 int64                 `json:"auction_configuration_id"`
	ConfigUID                string                `json:"auction_configuration_uid"`
	ExternalWinNotifications bool                  `json:"external_win_notifications"`
	Rounds                   []auction.RoundConfig `json:"rounds"`
	Segment                  Segment               `json:"segment"`
}

type Segment struct {
	ID string `json:"id"`
}
