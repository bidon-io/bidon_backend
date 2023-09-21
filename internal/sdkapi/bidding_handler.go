package sdkapi

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
)

type BiddingHandler struct {
	*BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]
	BiddingBuilder        *bidding.Builder
	SegmentMatcher        *segment.Matcher
	AdaptersConfigBuilder AdaptersConfigBuilder
	EventLogger           *event.Logger
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/bidding_mocks.go -pkg mocks . AdaptersConfigBuilder

type AdaptersConfigBuilder interface {
	Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp) (adapter.ProcessedConfigsMap, error)
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
	imp := req.raw.Imp
	adapterConfigs, err := h.AdaptersConfigBuilder.Build(ctx, req.app.ID, req.raw.Adapters.Keys(), imp)
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
	auctionResult, err := h.BiddingBuilder.HoldAuction(ctx, params)
	c.Logger().Printf("[BIDDING] bids: (%+v), err: (%s), took (%s)", auctionResult, err, time.Since(start))
	if err != nil {
		return err
	}

	h.sendEvents(c, req, &auctionResult)

	response := BiddingResponse{
		Status: "NO_BID",
	}

	for _, result := range auctionResult.Bids {
		if result.IsBid() && result.Bid.Price >= imp.GetBidFloor() {
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
	// Sort bids by price descending
	sort.Slice(response.Bids, func(i, j int) bool {
		return response.Bids[i].Price > response.Bids[j].Price
	})

	return c.JSON(http.StatusOK, response)
}

func (h *BiddingHandler) sendEvents(c echo.Context, req *request[schema.BiddingRequest, *schema.BiddingRequest], auctionResult *bidding.AuctionResult) {
	imp := req.raw.Imp
	auctionConfigurationUID, err := strconv.Atoi(imp.AuctionConfigUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	for _, result := range auctionResult.Bids {
		adRequestParams := event.AdRequestParams{
			EventType:               "bid_request",
			AdType:                  string(req.raw.AdType),
			AuctionID:               imp.AuctionID,
			AuctionConfigurationID:  imp.AuctionConfigID,
			AuctionConfigurationUID: int64(auctionConfigurationUID),
			Status:                  fmt.Sprint(result.Status),
			RoundID:                 imp.RoundID,
			RoundNumber:             auctionResult.RoundNumber,
			ImpID:                   imp.ID,
			DemandID:                string(result.DemandID),
			AdUnitID:                0,
			LineItemUID:             0,
			AdUnitCode:              "",
			Ecpm:                    0,
			PriceFloor:              imp.GetBidFloor(),
			Bidding:                 true,
			RawRequest:              result.RawRequest,
			RawResponse:             result.RawResponse,
		}
		bidRequestEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
		h.EventLogger.Log(bidRequestEvent, func(err error) {
			logError(c, fmt.Errorf("log bid_request event: %v", err))
		})
		if result.IsBid() {
			adRequestParams = event.AdRequestParams{
				EventType:               "bid",
				AdType:                  string(req.raw.AdType),
				AuctionID:               imp.AuctionID,
				AuctionConfigurationID:  imp.AuctionConfigID,
				AuctionConfigurationUID: int64(auctionConfigurationUID),
				Status:                  "SUCCESS",
				RoundID:                 imp.RoundID,
				RoundNumber:             auctionResult.RoundNumber,
				ImpID:                   imp.ID,
				DemandID:                string(result.DemandID),
				AdUnitID:                0,
				LineItemUID:             0,
				AdUnitCode:              result.TagID,
				Ecpm:                    result.Bid.Price,
				PriceFloor:              imp.GetBidFloor(),
				Bidding:                 true,
			}
			bidEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
			h.EventLogger.Log(bidEvent, func(err error) {
				logError(c, fmt.Errorf("log bid event: %v", err))
			})
		}
	}
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
