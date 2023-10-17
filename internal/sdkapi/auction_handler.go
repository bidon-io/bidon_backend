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
	AuctionBuilder *auction.Builder
	SegmentMatcher *segment.Matcher
	EventLogger    *event.Logger
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
	}
	auc, err := h.AuctionBuilder.Build(c.Request().Context(), params)
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
	adRequestParams := event.AdRequestParams{
		EventType:               "auction_request",
		AdType:                  string(req.raw.AdType),
		AuctionID:               req.raw.AdObject.AuctionID,
		AuctionConfigurationID:  auc.ConfigID,
		AuctionConfigurationUID: int64(auctionConfigurationUID),
		Status:                  "",
		RoundID:                 "",
		RoundNumber:             0,
		ImpID:                   "",
		DemandID:                "",
		AdUnitID:                0,
		LineItemUID:             0,
		AdUnitCode:              "",
		Ecpm:                    0,
		PriceFloor:              req.raw.AdObject.PriceFloor,
	}
	aucRequestEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
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
