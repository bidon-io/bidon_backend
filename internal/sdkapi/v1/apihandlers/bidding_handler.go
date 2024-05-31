package apihandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type BiddingHandler struct {
	*BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]
	BiddingBuilder        BiddingBuilder
	AdUnitsMatcher        AdUnitsMatcher
	AdaptersConfigBuilder AdaptersConfigBuilder
	EventLogger           *event.Logger
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/bidding_mocks.go -pkg mocks . BiddingBuilder AdaptersConfigBuilder AdUnitsMatcher

type BiddingBuilder interface {
	HoldAuction(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error)
}

type AdaptersConfigBuilder interface {
	Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp, adUnitsMap *map[adapter.Key][]auction.AdUnit) (adapter.ProcessedConfigsMap, error)
}

type AdUnitsMatcher interface {
	MatchCached(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error)
}

type BiddingResponse struct {
	Bids   []Bid  `json:"bids,omitempty"`
	Status string `json:"status"`
}

type Bid struct {
	ID      string
	ImpID   string
	Price   float64
	Demands *map[adapter.Key]Demand // Deprecated: use AdUnit instead of Demands since SDK 0.5
	AdUnit  *AdUnit
	Ext     map[string]interface{}
}

func (b *Bid) MarshalJSON() ([]byte, error) {
	if b.Demands != nil {
		return json.Marshal(map[string]interface{}{
			"id":      b.ID,
			"impid":   b.ImpID,
			"price":   b.Price,
			"demands": b.Demands,
		})
	} else {
		return json.Marshal(map[string]interface{}{
			"id":      b.ID,
			"imp_id":  b.ImpID,
			"price":   b.Price,
			"ad_unit": b.AdUnit,
			"ext":     b.Ext,
		})
	}
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

	sdkVersion, err := req.raw.GetSDKVersionSemver()
	if err != nil {
		return sdkapi.ErrInvalidSDKVersion
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
	adUnits, err := h.AdUnitsMatcher.MatchCached(ctx, &auction.BuildParams{
		Adapters:   req.raw.Adapters.Keys(),
		AppID:      req.app.ID,
		AdType:     req.raw.Imp.Type(),
		AdFormat:   req.raw.Imp.Format(),
		DeviceType: req.raw.Device.Type,
	})
	if err != nil {
		return err
	}

	adUnitsMap := make(map[adapter.Key][]auction.AdUnit)
	for _, adUnit := range adUnits {
		key := adapter.Key(adUnit.DemandID)
		adUnitsMap[key] = append(adUnitsMap[key], adUnit)
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
		StartTS:        start.UnixMilli(),
	}
	auctionResult, err := h.BiddingBuilder.HoldAuction(ctx, params)
	c.Logger().Printf("[BIDDING] bids: (%+v), err: (%s), took (%s)", auctionResult, err, time.Since(start))
	if err != nil {
		return err
	}

	events := prepareBiddingEvents(req, &auctionResult, &adUnitsMap)
	for _, ev := range events {
		h.EventLogger.Log(ev, func(err error) {
			sdkapi.LogError(c, fmt.Errorf("log %v event: %v", ev.EventType, err))
		})
	}

	response, err := h.buildResponse(auctionResult, imp, adUnitsMap, sdkVersion)
	if err != nil {
		c.Logger().Errorf("Error building response: ", err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *BiddingHandler) buildResponse(auctionResult bidding.AuctionResult, imp schema.Imp, adUnitsMap map[adapter.Key][]auction.AdUnit, sdkVersion *semver.Version) (BiddingResponse, error) {
	response := BiddingResponse{
		Status: "NO_BID",
	}

	for _, result := range auctionResult.Bids {
		if result.IsBid() && result.Price() >= imp.GetBidFloor() {
			var bid *Bid

			if sdkapi.Version05GTEConstraint.Check(sdkVersion) {
				bid = buildBid(result, &adUnitsMap)
			} else {
				bid = h.buildBidDeprecated(result)
			}

			if bid != nil {
				response.Bids = append(response.Bids, *bid)
			}
		}
	}
	if len(response.Bids) > 0 {
		response.Status = "SUCCESS"
	}

	// Sort bids by price descending
	sort.Slice(response.Bids, func(i, j int) bool {
		return response.Bids[i].Price > response.Bids[j].Price
	})
	return response, nil
}

func buildBid(demandResponse adapters.DemandResponse, adUnitsMap *map[adapter.Key][]auction.AdUnit) *Bid {
	storeAdUnit, err := selectAdUnit(demandResponse, adUnitsMap)
	if err != nil {
		return nil
	}

	var adUnit AdUnit

	if storeAdUnit != nil {
		adUnit = AdUnit{
			DemandID: demandResponse.DemandID,
			UID:      storeAdUnit.UID,
			Label:    storeAdUnit.Label,
			BidType:  storeAdUnit.BidType,
			Extra:    storeAdUnit.Extra,
		}
	}

	return &Bid{
		ID:     demandResponse.Bid.ID,
		ImpID:  demandResponse.Bid.ImpID,
		Price:  demandResponse.Bid.Price,
		AdUnit: &adUnit,
		Ext:    buildDemandExt(demandResponse),
	}
}

// Deprecated: uses AdUnit instead of Demands since SDK 0.5
func (h *BiddingHandler) buildBidDeprecated(demandResponse adapters.DemandResponse) *Bid {
	return &Bid{
		ID:    demandResponse.Bid.ID,
		ImpID: demandResponse.Bid.ImpID,
		Price: demandResponse.Bid.Price,
		Demands: &map[adapter.Key]Demand{
			demandResponse.DemandID: buildDemandInfo(demandResponse),
		},
	}
}

func prepareBiddingEvents(
	req *request[schema.BiddingRequest, *schema.BiddingRequest],
	auctionResult *bidding.AuctionResult,
	adUnitsMap *map[adapter.Key][]auction.AdUnit,
) []*event.AdEvent {
	imp := req.raw.Imp
	auctionConfigurationUID, err := strconv.Atoi(imp.AuctionConfigurationUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	events := make([]*event.AdEvent, 0, len(auctionResult.Bids))
	for _, result := range auctionResult.Bids {
		adUnit, _ := selectAdUnit(result, adUnitsMap)
		adUnitUID := int64(0)
		adUnitLabel := ""
		if adUnit != nil {
			uid, _ := strconv.ParseInt(adUnit.UID, 10, 64)
			adUnitUID = uid
			adUnitLabel = adUnit.Label
		}

		adRequestParams := event.AdRequestParams{
			EventType:               "bid_request",
			AdType:                  string(req.raw.AdType),
			AdFormat:                string(req.raw.Imp.Format()),
			AuctionID:               imp.AuctionID,
			AuctionConfigurationID:  imp.AuctionConfigurationID,
			AuctionConfigurationUID: int64(auctionConfigurationUID),
			Status:                  fmt.Sprint(result.Status),
			RoundID:                 imp.RoundID,
			RoundNumber:             auctionResult.RoundNumber,
			ImpID:                   "",
			DemandID:                string(result.DemandID),
			AdUnitUID:               adUnitUID,
			AdUnitLabel:             adUnitLabel,
			ECPM:                    result.Price(),
			PriceFloor:              imp.GetBidFloor(),
			Bidding:                 true,
			RawRequest:              result.RawRequest,
			RawResponse:             result.RawResponse,
			Error:                   result.ErrorMessage(),
			TimingMap: event.TimingMap{
				"bid": {result.StartTS, result.EndTS},
			},
		}
		events = append(events, event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData))
		if result.IsBid() {
			adRequestParams = event.AdRequestParams{
				EventType:               "bid",
				AdType:                  string(req.raw.AdType),
				AdFormat:                string(req.raw.Imp.Format()),
				AuctionID:               imp.AuctionID,
				AuctionConfigurationID:  imp.AuctionConfigurationID,
				AuctionConfigurationUID: int64(auctionConfigurationUID),
				Status:                  "SUCCESS",
				RoundID:                 imp.RoundID,
				RoundNumber:             auctionResult.RoundNumber,
				ImpID:                   "",
				DemandID:                string(result.DemandID),
				AdUnitUID:               adUnitUID,
				AdUnitLabel:             adUnitLabel,
				ECPM:                    result.Bid.Price,
				PriceFloor:              imp.GetBidFloor(),
				Bidding:                 true,
				TimingMap: event.TimingMap{
					"bid": {result.StartTS, result.EndTS},
				},
			}
			events = append(events, event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData))
		}
	}

	return events
}

func selectAdUnit(demandResponse adapters.DemandResponse, adUnitsMap *map[adapter.Key][]auction.AdUnit) (*auction.AdUnit, error) {
	adUnits, ok := (*adUnitsMap)[demandResponse.DemandID]
	if !ok {
		return nil, fmt.Errorf("ad units not found for demand %s", demandResponse.DemandID)
	}

	if demandResponse.DemandID == adapter.AmazonKey {
		for _, adUnit := range adUnits {
			if demandResponse.SlotUUID == adUnit.Extra["slot_uuid"] {
				return &adUnit, nil
			}
		}
	} else {
		adUnit := adUnits[0]
		return &adUnit, nil
	}

	return nil, fmt.Errorf("ad unit not found for demand %s", demandResponse.DemandID)
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

func buildDemandExt(demandResponse adapters.DemandResponse) map[string]any {
	switch demandResponse.DemandID {
	case adapter.AmazonKey:
		return map[string]any{}
	case adapter.MobileFuseKey:
		return map[string]any{
			"signaldata": demandResponse.Bid.Signaldata,
		}
	default:
		return map[string]any{
			"payload": demandResponse.Bid.Payload,
		}
	}
}
