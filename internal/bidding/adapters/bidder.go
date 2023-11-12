package adapters

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/bidding/openrtb"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type BidderInterface interface {
	// CreateRequest makes the HTTP requests which should be made to fetch bids.
	CreateRequest(openrtb.BidRequest, *schema.BiddingRequest) (openrtb.BidRequest, error)

	// ExecuteRequest sends request to bidder endpoint.
	ExecuteRequest(context.Context, *http.Client, openrtb.BidRequest) *DemandResponse

	// ParseBids unpacks the server's response into Bids.
	ParseBids(*DemandResponse) (*DemandResponse, error)
}

type Bidder struct {
	Adapter BidderInterface
	Client  *http.Client
}

type Builder func(adapter.ProcessedConfigsMap, *http.Client) (*Bidder, error)

type DemandResponse struct {
	DemandID    adapter.Key
	RequestID   string
	RawRequest  string
	RawResponse string
	Status      int
	Bid         *BidDemandResponse
	Error       error
	TagID       string
	PlacementID string
	SlotUUID    string
	TimeoutURL  string
}

func (dr *DemandResponse) IsBid() bool {
	return dr.Bid != nil
}

func (dr *DemandResponse) ErrorMessage() string {
	errMsg := ""
	if dr.Error != nil {
		errMsg = dr.Error.Error()
	}

	return errMsg
}

func (dr *DemandResponse) Price() float64 {
	price := float64(0)
	if dr.IsBid() {
		price = dr.Bid.Price
	}
	return price
}

type BidDemandResponse struct {
	Payload    string
	Signaldata string
	ID         string
	ImpID      string
	AdID       string
	SeatID     string
	DemandID   adapter.Key
	Price      float64
	LURL       string
	NURL       string
	BURL       string
}
