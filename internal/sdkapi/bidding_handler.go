package sdkapi

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/bidon-io/bidon-backend/internal/adapter/store"
	"github.com/bidon-io/bidon-backend/internal/auction"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type BiddingHandler struct {
	*BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]
	BiddingBuilder        *bidding.Builder
	AdUnitsMapBuilder     AdUnitsMapBuilder
	AdaptersConfigBuilder AdaptersConfigBuilder
	EventLogger           *event.Logger
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/bidding_mocks.go -pkg mocks . AdaptersConfigBuilder AdUnitsMapBuilder

type AdaptersConfigBuilder interface {
	Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp, adUnitsMap *store.AdUnitsMap) (adapter.ProcessedConfigsMap, error)
}

type AdUnitsMapBuilder interface {
	Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp) (store.AdUnitsMap, error)
}

type BiddingResponse struct {
	Bids   []Bid  `json:"bids,omitempty"`
	Status string `json:"status"`
}

type Bid struct {
	ID      string                 `json:"id"`
	ImpID   string                 `json:"impid"`
	Price   float64                `json:"price"`
	Demands map[adapter.Key]Demand `json:"demands,omitempty"` // Deprecated: uses AdUnit instead of Demands since SDK 0.5
	AdUnit  AdUnit                 `json:"ad_unit,omitempty"`
	Ext     map[string]interface{} `json:"ext,omitempty"` // TODO: remove interface{} with concrete type
}

type Demand struct {
	Payload     string `json:"payload"`
	Signaldata  string `json:"signaldata"`
	UnitID      string `json:"unit_id,omitempty"`
	SlotID      string `json:"slot_id,omitempty"`
	SlotUUID    string `json:"slot_uuid,omitempty"`
	PlacementID string `json:"placement_id,omitempty"`
}

type AdUnit struct {
	DemandID adapter.Key    `json:"demand_id"`
	UID      string         `json:"uid"`
	Label    string         `json:"label"`
	BidType  schema.BidType `json:"bid_type"`
	Extra    map[string]any `json:"ext"`
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

	imp := req.raw.Imp
	adUnitsMap, err := h.AdUnitsMapBuilder.Build(ctx, req.app.ID, req.raw.Adapters.Keys(), imp)
	if err != nil {
		return err
	}

	adapterConfigs, err := h.AdaptersConfigBuilder.Build(ctx, req.app.ID, req.raw.Adapters.Keys(), imp, &adUnitsMap)
	if err != nil {
		return err
	}

	auctionConfig := req.auctionConfig
	if auctionConfig == nil {
		return fmt.Errorf("cannot find config: %v", imp.AuctionConfigurationID)
	}

	params := &bidding.BuildParams{
		AppID:          req.app.ID,
		BiddingRequest: req.raw,
		GeoData:        req.geoData,
		AdapterConfigs: adapterConfigs,
		AuctionConfig:  *auctionConfig,
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

	sdkVersion, err := req.raw.GetSDKVersionSemver()
	if err != nil {
		return ErrInvalidSDKVersion
	}

	if Version05GTEConstraint.Check(sdkVersion) {
		h.buildBids(auctionResult, imp, &adUnitsMap, &response)
	} else {
		h.buildBidsDeprecated(auctionResult, imp, &response)
	}
	// Sort bids by price descending
	sort.Slice(response.Bids, func(i, j int) bool {
		return response.Bids[i].Price > response.Bids[j].Price
	})

	return c.JSON(http.StatusOK, response)
}

func (h *BiddingHandler) buildBids(auctionResult bidding.AuctionResult, imp schema.Imp, adUnitsMap *store.AdUnitsMap, response *BiddingResponse) {
	for _, result := range auctionResult.Bids {
		if result.IsBid() && result.Bid.Price >= imp.GetBidFloor() {

			adUnit, err := selectAdUnit(result, adUnitsMap)
			if err != nil {
				continue
			}
			response.Bids = append(response.Bids, Bid{
				ID:     result.Bid.ID,
				ImpID:  result.Bid.ImpID,
				Price:  result.Bid.Price,
				AdUnit: *adUnit,
				Ext: map[string]any{
					"payload": result.Bid.Payload,
				},
			})
			response.Status = "SUCCESS"
		}
	}
}

// Deprecated: uses AdUnit instead of Demands since SDK 0.5
func (h *BiddingHandler) buildBidsDeprecated(auctionResult bidding.AuctionResult, imp schema.Imp, response *BiddingResponse) {
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
}

func (h *BiddingHandler) sendEvents(c echo.Context, req *request[schema.BiddingRequest, *schema.BiddingRequest], auctionResult *bidding.AuctionResult) {
	imp := req.raw.Imp
	auctionConfigurationUID, err := strconv.Atoi(imp.AuctionConfigurationUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	for _, result := range auctionResult.Bids {
		adRequestParams := event.AdRequestParams{
			EventType:               "bid_request",
			AdType:                  string(req.raw.AdType),
			AuctionID:               imp.AuctionID,
			AuctionConfigurationID:  imp.AuctionConfigurationID,
			AuctionConfigurationUID: int64(auctionConfigurationUID),
			Status:                  fmt.Sprint(result.Status),
			RoundID:                 imp.RoundID,
			RoundNumber:             auctionResult.RoundNumber,
			ImpID:                   imp.ID,
			DemandID:                string(result.DemandID),
			AdUnitID:                0,
			LineItemUID:             0,
			LineItemLabel:           "",
			AdUnitCode:              "",
			Ecpm:                    0,
			PriceFloor:              imp.GetBidFloor(),
			Bidding:                 true,
			RawRequest:              result.RawRequest,
			RawResponse:             result.RawResponse,
			Error:                   result.ErrorMessage(),
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
				AuctionConfigurationID:  imp.AuctionConfigurationID,
				AuctionConfigurationUID: int64(auctionConfigurationUID),
				Status:                  "SUCCESS",
				RoundID:                 imp.RoundID,
				RoundNumber:             auctionResult.RoundNumber,
				ImpID:                   imp.ID,
				DemandID:                string(result.DemandID),
				AdUnitID:                0,
				LineItemUID:             0,
				LineItemLabel:           "",
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

func selectAdUnit(demandResponse adapters.DemandResponse, adUnitsMap *store.AdUnitsMap) (*AdUnit, error) {
	adUnits, ok := (*adUnitsMap)[demandResponse.DemandID]
	if !ok {
		return nil, fmt.Errorf("ad units not found for demand %s", demandResponse.DemandID)
	}
	if demandResponse.DemandID == adapter.AmazonKey {
		return selectAmazonAdUnit(demandResponse.SlotUUID, adUnits)
	}

	adUnit := adUnits[0]
	return &AdUnit{
		DemandID: demandResponse.DemandID,
		UID:      adUnit.UID,
		Label:    adUnit.Label,
		BidType:  adUnit.BidType,
		Extra:    adUnit.Extra,
	}, nil
}

func selectAmazonAdUnit(slotUUID string, adUnits []auction.AdUnit) (*AdUnit, error) {
	for _, adUnit := range adUnits {
		if slotUUID == adUnit.Extra["slot_uuid"] {
			return &AdUnit{
				UID:   adUnit.UID,
				Label: adUnit.Label,
				Extra: adUnit.Extra,
			}, nil
		}
	}

	return nil, fmt.Errorf("ad unit not found for slot_uuid %s", slotUUID)
}

// Deprecated: uses AdUnit instead of Demands since SDK 0.5
func buildDemandInfo(demandResponse adapters.DemandResponse) Demand {
	switch demandResponse.DemandID {
	case adapter.AmazonKey:
		return Demand{
			SlotUUID: demandResponse.SlotUUID,
		}
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
	case adapter.MobileFuseKey:
		return Demand{
			Payload:     demandResponse.Bid.Payload,
			Signaldata:  demandResponse.Bid.Signaldata,
			PlacementID: demandResponse.TagID,
		}
	default:
		return Demand{
			Payload:     demandResponse.Bid.Payload,
			PlacementID: demandResponse.TagID,
		}
	}
}
