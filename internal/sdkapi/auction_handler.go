package sdkapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
)

type AuctionHandler struct {
	*BaseHandler[schema.AuctionRequest, *schema.AuctionRequest]
	AuctionBuilder   *auction.Builder
	AuctionBuilderV2 *auction.BuilderV2
	SegmentMatcher   *segment.Matcher
	EventLogger      *event.Logger
}

type AuctionResponse struct {
	*auction.Auction
	Token      string  `json:"token"`
	PriceFloor float64 `json:"pricefloor"`
	AuctionID  string  `json:"auction_id"`
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

	params := &auction.BuildParams{
		AppID:      req.app.ID,
		AdType:     req.raw.AdType,
		AdFormat:   req.raw.AdObject.Format(),
		DeviceType: req.raw.Device.Type,
		Adapters:   req.raw.Adapters.Keys(),
		Segment:    sgmnt,
		PriceFloor: &req.raw.AdObject.PriceFloor,
		AuctionKey: req.raw.AdObject.AuctionKey,
	}

	sdkVersion, err := req.raw.GetSDKVersionSemver()
	if err != nil {
		return ErrInvalidSDKVersion
	}

	var auc *auction.Auction
	if Version05GTEConstraint.Check(sdkVersion) {
		auc, err = h.AuctionBuilderV2.Build(c.Request().Context(), params)
	} else {
		auc, err = h.AuctionBuilder.Build(c.Request().Context(), params)
	}

	if err != nil {
		if errors.Is(err, auction.ErrNoAdsFound) {
			err = ErrNoAdsFound
		}

		return err
	}

	auctionConfigurationUID, err := strconv.Atoi(auc.ConfigUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	aucRequestEvent := prepareAuctionRequestEvent(req, auc, auctionConfigurationUID)
	h.EventLogger.Log(aucRequestEvent, func(err error) {
		logError(c, fmt.Errorf("log auction_request event: %v", err))
	})

	response := &AuctionResponse{
		Auction:    auc,
		Token:      "{}",
		PriceFloor: req.raw.AdObject.PriceFloor,
		AuctionID:  req.raw.AdObject.AuctionID,
	}

	return c.JSON(http.StatusOK, response)
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
