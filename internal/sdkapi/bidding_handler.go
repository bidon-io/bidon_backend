package sdkapi

import (
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type BiddingHandler struct {
	*BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]
	BiddingBuilder *bidding.Builder
	SegmentMatcher *segment.Matcher
}

type BiddingResponse struct {
	Bid    *Bid   `json:"bid,omitempty"`
	Status string `json:"status"`
}

type Bid struct {
	ID       uuid.UUID              `json:"id"`
	ImpID    uuid.UUID              `json:"impid"`
	Price    float64                `json:"price"`
	Payload  string                 `json:"payload"`
	DemandID adapter.Key            `json:"demand_id"`
	Ext      map[string]interface{} `json:"ext"`
}

func (h *BiddingHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	segmentParams := &segment.Params{
		Country: req.countryCode(),
		Ext:     req.raw.Segment.Ext,
		AppID:   req.app.ID,
	}

	sgmnt := h.SegmentMatcher.Match(c.Request().Context(), segmentParams)

	params := &bidding.BuildParams{
		AppID:      req.app.ID,
		AdType:     req.raw.AdType,
		AdFormat:   req.raw.Imp.Format(),
		DeviceType: req.raw.Device.Type,
		Adapters:   req.raw.Adapters.Keys(),
		SegmentID:  sgmnt.ID,
	}
	bid, err := h.BiddingBuilder.Build(c.Request().Context(), params)
	if err != nil {
		return err
	}

	response := &BiddingResponse{
		Status: "NO_BID",
	}

	if bid.IsBid() {
		id, _ := uuid.NewV4()    // Temporary stub
		impid, _ := uuid.NewV4() // Temporary stub
		response.Bid = &Bid{
			ID:       id,
			ImpID:    impid,
			Price:    bid.Price,
			Payload:  "Some payload",
			DemandID: adapter.BidmachineKey,
		}
		response.Status = "SUCCESS"

	}

	return c.JSON(http.StatusOK, response)
}
