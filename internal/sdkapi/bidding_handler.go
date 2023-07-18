package sdkapi

import (
	"context"
	"net/http"
	"time"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters_builder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
)

type BiddingHandler struct {
	*BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]
	BiddingBuilder        *bidding.Builder
	SegmentMatcher        *segment.Matcher
	AdaptersConfigBuilder *adapters_builder.AdaptersConfigBuilder
}

type BiddingResponse struct {
	Bid    *Bid   `json:"bid,omitempty"`
	Status string `json:"status"`
}

type Bid struct {
	ID       string                 `json:"id"`
	ImpID    string                 `json:"impid"`
	Price    float64                `json:"price"`
	Payload  string                 `json:"payload"`
	DemandID adapter.Key            `json:"demand_id"`
	Ext      map[string]interface{} `json:"ext,omitempty"` // TODO: remove interface{} with concrete type
}

func (h *BiddingHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	start := time.Now()

	ctx := c.Request().Context()

	timeout := time.Duration(req.raw.TMax) * time.Millisecond
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, start.Add(timeout))
		defer cancel()
	}

	segmentParams := &segment.Params{
		Country: req.countryCode(),
		Ext:     req.raw.Segment.Ext,
		AppID:   req.app.ID,
	}

	sgmnt := h.SegmentMatcher.Match(ctx, segmentParams)
	adapterConfigs, err := h.AdaptersConfigBuilder.Build(ctx, req.app.ID, req.raw.Adapters.Keys())
	if err != nil {
		return err
	}

	params := &bidding.BuildParams{
		AppID:          req.app.ID,
		BiddingRequest: req.raw,
		SegmentID:      sgmnt.ID,
		GeoData:        req.geoData,
		AdapterConfigs: adapterConfigs,
	}
	result, err := h.BiddingBuilder.HoldAuction(ctx, params)
	if err != nil && err != bidding.ErrNoBids {
		return err
	}

	response := &BiddingResponse{
		Status: "NO_BID",
	}

	if result.IsBid() {
		response.Bid = &Bid{
			ID:       result.Bid.ID,
			ImpID:    result.Bid.ImpID,
			Price:    result.Bid.Price,
			Payload:  result.Bid.Payload,
			DemandID: result.Bid.DemandID,
		}
		response.Status = "SUCCESS"

	}

	return c.JSON(http.StatusOK, response)
}
