package sdkapi

import (
	"context"
	"net/http"
	"time"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
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
	Bids   []Bid  `json:"bids,omitempty"`
	Status string `json:"status"`
}

type Bid struct {
	ID      string                 `json:"id"`
	ImpID   string                 `json:"impid"`
	Price   float64                `json:"price"`
	Demands map[adapter.Key]Demand `json:"demands"`
	Ext     map[string]interface{} `json:"ext,omitempty"` // TODO: remove interface{} with concrete type
}

type Demand struct {
	Payload     string `json:"payload"`
	UnitID      string `json:"unit_id,omitempty"`
	SlotID      string `json:"slot_id,omitempty"`
	PlacementID string `json:"placement_id,omitempty"`
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
	adapterConfigs, err := h.AdaptersConfigBuilder.Build(ctx, req.app.ID, req.raw.Adapters.Keys(), req.raw.Imp)
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
	demandResponses, err := h.BiddingBuilder.HoldAuction(ctx, params)
	if err != nil && err != bidding.ErrNoBids {
		return err
	}

	response := BiddingResponse{
		Status: "NO_BID",
	}

	for _, result := range demandResponses {
		if result.IsBid() {
			response.Bids = append(response.Bids, Bid{
				ID:    result.Bid.ID,
				ImpID: result.Bid.ImpID,
				Price: result.Bid.Price,
				Demands: map[adapter.Key]Demand{
					result.DemandID: buildDemandInfo(result),
				},
			})
			response.Status = "SUCCESS"
		}
	}

	return c.JSON(http.StatusOK, response)
}

func buildDemandInfo(demandResponse adapters.DemandResponse) Demand {
	switch demandResponse.DemandID {
	case adapter.MintegralKey:
		return Demand{
			Payload:     demandResponse.Bid.Payload,
			UnitID:      demandResponse.TagID,
			PlacementID: demandResponse.PlacementID,
		}
	case adapter.BidmachineKey:
		return Demand{
			Payload: demandResponse.Bid.Payload,
		}
	case adapter.BigoAdsKey:
		return Demand{
			Payload: demandResponse.Bid.Payload,
			SlotID:  demandResponse.TagID,
		}
	default:
		return Demand{
			Payload:     demandResponse.Bid.Payload,
			PlacementID: demandResponse.TagID,
		}
	}
}
