package sdkapi

import (
	"errors"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auctionv2"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"strconv"
)

type AuctionHandlerV2 struct {
	*BaseHandler[schema.AuctionV2Request, *schema.AuctionV2Request]
	AuctionBuilder        *auctionv2.Builder
	SegmentMatcher        *segment.Matcher
	BiddingBuilder        BiddingBuilder
	AdUnitsMatcher        AdUnitsMatcher
	AdaptersConfigBuilder AdaptersConfigBuilder
	EventLogger           *event.Logger
}

type AuctionV2Response struct {
	ConfigID                 int64            `json:"auction_configuration_id"`
	ConfigUID                string           `json:"auction_configuration_uid"`
	ExternalWinNotifications bool             `json:"external_win_notifications"`
	AdUnits                  []auction.AdUnit `json:"ad_units"`
	Segment                  auction.Segment  `json:"segment"`
	Token                    string           `json:"token"`
	PriceFloor               float64          `json:"pricefloor"`
	AuctionID                string           `json:"auction_id"`
}

func (h *AuctionHandlerV2) Handle(c echo.Context) error {
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
		AdFormat:             req.raw.AdObjectV2.Format(),
		DeviceType:           req.raw.Device.Type,
		Adapters:             req.raw.Adapters.Keys(),
		Segment:              sgmnt,
		PriceFloor:           req.raw.AdObjectV2.PriceFloor,
		MergedAuctionRequest: &req.raw,
		GeoData:              req.geoData,
	}

	auctionResult, err := h.AuctionBuilder.Build(c.Request().Context(), params)
	if err != nil {
		if errors.Is(err, auction.ErrNoAdsFound) {
			err = ErrNoAdsFound
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

func (h *AuctionHandlerV2) buildResponse(
	req *request[schema.AuctionV2Request, *schema.AuctionV2Request],
	auctionResult *auctionv2.AuctionResult,
	adUnitsMap *map[adapter.Key][]auction.AdUnit,
) (*AuctionV2Response, error) {
	adObject := req.raw.AdObjectV2
	response := AuctionV2Response{
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

func (h *AuctionHandlerV2) logEvents(
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
			logError(c, fmt.Errorf("log %v event: %v", ev.EventType, err))
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
