package bidding

import (
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
)

type Bidding struct {
	ConfigID                 int64                 `json:"auction_configuration_id"`
	ExternalWinNotifications bool                  `json:"external_win_notifications"`
	Rounds                   []auction.RoundConfig `json:"rounds"`
	Segment                  Segment               `json:"segment"`
}

type DemandResponse struct {
	DemandID    adapter.Key
	RawRequest  string
	RawResponse string
	Status      int
	Price       float64
	Bid         *BidDemandResponse
}

func (m *DemandResponse) IsBid() bool {
	return m.Bid != nil
}

type BidDemandResponse struct {
	Payload  string
	ID       string
	ImpID    string
	AdID     string
	SeatID   string
	DemandID string
	Price    float64
	LURL     string
	NURL     string
	BURL     string
}

type Segment struct {
	ID string `json:"id"`
}
