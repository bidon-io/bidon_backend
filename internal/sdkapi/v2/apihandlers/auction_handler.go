package apihandlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auctionv2"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"strconv"
)

type AuctionHandler struct {
	*BaseHandler[schema.AuctionV2Request, *schema.AuctionV2Request]
	AuctionBuilder        *auctionv2.Builder
	SegmentMatcher        *segment.Matcher
	BiddingBuilder        BiddingBuilder
	AdUnitsMatcher        AdUnitsMatcher
	AdaptersConfigBuilder AdaptersConfigBuilder
	EventLogger           *event.Logger
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/auction_mocks.go -pkg mocks . BiddingBuilder AdaptersConfigBuilder AdUnitsMatcher

type BiddingBuilder interface {
	HoldAuction(ctx context.Context, params *bidding.BuildParams) (bidding.AuctionResult, error)
}

type AdaptersConfigBuilder interface {
	Build(ctx context.Context, appID int64, adapterKeys []adapter.Key, imp schema.Imp, adUnitsMap *map[adapter.Key][]auction.AdUnit) (adapter.ProcessedConfigsMap, error)
}

type AdUnitsMatcher interface {
	MatchCached(ctx context.Context, params *auction.BuildParams) ([]auction.AdUnit, error)
}

type AuctionResponse struct {
	ConfigID                 int64            `json:"auction_configuration_id"`
	ConfigUID                string           `json:"auction_configuration_uid"`
	ExternalWinNotifications bool             `json:"external_win_notifications"`
	AdUnits                  []auction.AdUnit `json:"ad_units"`
	Segment                  auction.Segment  `json:"segment"`
	Token                    string           `json:"token"`
	PriceFloor               float64          `json:"pricefloor"`
	AuctionID                string           `json:"auction_id"`
}

func (h *AuctionHandler) Handle(c echo.Context) error {
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
	req.raw.Segment.ID = sgmnt.StringID()
	req.raw.Segment.UID = sgmnt.UID

	params := &auctionv2.BuildParams{
		AppID:                req.app.ID,
		AdType:               req.raw.AdType,
		AdFormat:             req.raw.AdObject.Format(),
		DeviceType:           req.raw.Device.Type,
		Adapters:             req.raw.Adapters.Keys(),
		Segment:              sgmnt,
		PriceFloor:           req.raw.AdObject.PriceFloor,
		MergedAuctionRequest: &req.raw,
		GeoData:              req.geoData,
		AuctionKey:           req.raw.AdObject.AuctionKey,
	}

	auctionResult, err := h.AuctionBuilder.Build(c.Request().Context(), params)
	if err != nil {
		if errors.Is(err, auction.ErrNoAdsFound) {
			err = sdkapi.ErrNoAdsFound
		}

		return err
	}
	c.Logger().Printf("[AUCTION V2] auction: (%+v), err: (%s), took (%ms)", auctionResult, err, auctionResult.Stat.DurationTS)

	adUnitsMap := make(map[adapter.Key][]auction.AdUnit)
	for _, adUnit := range *auctionResult.AdUnits {
		key := adapter.Key(adUnit.DemandID)
		adUnitsMap[key] = append(adUnitsMap[key], adUnit)
	}

	h.logEvents(c, req, auctionResult, &adUnitsMap)

	response, err := h.buildResponse(req, auctionResult, &adUnitsMap)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AuctionHandler) buildResponse(
	req *request[schema.AuctionV2Request, *schema.AuctionV2Request],
	auctionResult *auctionv2.AuctionResult,
	adUnitsMap *map[adapter.Key][]auction.AdUnit,
) (*AuctionResponse, error) {
	adObject := req.raw.AdObject
	response := AuctionResponse{
		ConfigID:   auctionResult.AuctionConfiguration.ID,
		ConfigUID:  auctionResult.AuctionConfiguration.UID,
		Segment:    auction.Segment{ID: req.raw.Segment.ID, UID: req.raw.Segment.UID},
		Token:      "{}",
		AuctionID:  adObject.AuctionID,
		PriceFloor: adObject.PriceFloor,
	}

	// Store CPM AdUnits from AuctionConfiguration
	for _, adUnit := range *auctionResult.AdUnits {
		if adUnit.BidType == schema.CPMBidType {
			response.AdUnits = append(response.AdUnits, adUnit)
		}
	}

	// Store Bids AS RTB AdUnits from BiddingAuctionResult
	for _, bidResponse := range auctionResult.BiddingAuctionResult.Bids {
		if bidResponse.IsBid() && bidResponse.Price() >= adObject.PriceFloor {
			adUnit := convertBidToAdUnit(bidResponse, adUnitsMap)
			if adUnit != nil {
				response.AdUnits = append(response.AdUnits, *adUnit)
			}
		}
	}

	// Sort AdUnits by price
	sort.Slice(response.AdUnits, func(i, j int) bool {
		return response.AdUnits[i].GetPriceFloor() > response.AdUnits[j].GetPriceFloor()
	})

	return &response, nil
}

func (h *AuctionHandler) logEvents(
	c echo.Context,
	req *request[schema.AuctionV2Request, *schema.AuctionV2Request],
	auctionResult *auctionv2.AuctionResult,
	adUnitsMap *map[adapter.Key][]auction.AdUnit,
) {
	auctionRequest := &request[schema.AuctionRequest, *schema.AuctionRequest]{
		raw:           req.raw.ToAuctionRequest(),
		app:           req.app,
		auctionConfig: req.auctionConfig,
		geoData:       req.geoData,
	}
	auc := &auction.Auction{
		ConfigID:  auctionResult.AuctionConfiguration.ID,
		ConfigUID: auctionResult.AuctionConfiguration.UID,
	}
	auctionConfigurationUID, err := strconv.Atoi(auc.ConfigUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	var firstRoundID string
	if len(auctionResult.AuctionConfiguration.Rounds) > 0 {
		firstRoundID = auctionResult.AuctionConfiguration.Rounds[0].ID
	}
	biddingRequest := &request[schema.BiddingRequest, *schema.BiddingRequest]{
		raw:           req.raw.ToBiddingRequest(firstRoundID),
		app:           req.app,
		auctionConfig: req.auctionConfig,
		geoData:       req.geoData,
	}
	events := prepareBiddingEvents(biddingRequest, auctionResult.BiddingAuctionResult, adUnitsMap)
	aucRequestEvent := prepareAuctionRequestEvent(auctionRequest, auc, auctionConfigurationUID)

	events = append(events, aucRequestEvent)
	for _, ev := range events {
		h.EventLogger.Log(ev, func(err error) {
			sdkapi.LogError(c, fmt.Errorf("log %v event: %v", ev.EventType, err))
		})
	}
}

func convertBidToAdUnit(demandResponse adapters.DemandResponse, adUnitsMap *map[adapter.Key][]auction.AdUnit) *auction.AdUnit {
	storeAdUnit, err := selectAdUnit(demandResponse, adUnitsMap)
	if err != nil {
		return nil
	}
	if storeAdUnit == nil {
		return nil
	}

	priceFloor := demandResponse.Price()
	ext := buildDemandExt(demandResponse)
	for key, value := range storeAdUnit.Extra {
		ext[key] = value
	}

	return &auction.AdUnit{
		DemandID:   string(demandResponse.DemandID),
		UID:        storeAdUnit.UID,
		Label:      storeAdUnit.Label,
		BidType:    storeAdUnit.BidType,
		PriceFloor: &priceFloor,
		Extra:      ext,
	}
}

func prepareAuctionRequestEvent(req *request[schema.AuctionRequest, *schema.AuctionRequest], auc *auction.Auction, auctionConfigurationUID int) *event.AdEvent {
	adRequestParams := event.AdRequestParams{
		EventType:               "auction_request",
		AdType:                  string(req.raw.AdType),
		AdFormat:                string(req.raw.AdObject.Format()),
		AuctionID:               req.raw.AdObject.AuctionID,
		AuctionConfigurationID:  auc.ConfigID,
		AuctionConfigurationUID: int64(auctionConfigurationUID),
		Status:                  "",
		RoundID:                 "",
		RoundNumber:             0,
		ImpID:                   "",
		DemandID:                "",
		AdUnitUID:               0,
		AdUnitLabel:             "",
		ECPM:                    0,
		PriceFloor:              req.raw.AdObject.PriceFloor,
	}

	return event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData)
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
