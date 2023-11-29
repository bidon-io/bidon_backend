package sdkapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type ClickHandler struct {
	*BaseHandler[schema.ClickRequest, *schema.ClickRequest]
	EventLogger *event.Logger
}

func (h *ClickHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	adEvent := prepareClickEvent(req)
	h.EventLogger.Log(adEvent, func(err error) {
		logError(c, fmt.Errorf("log click event: %v", err))
	})

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func prepareClickEvent(req *request[schema.ClickRequest, *schema.ClickRequest]) *event.RequestEvent {
	bid := req.raw.Bid

	auctionConfigurationUID, err := strconv.ParseInt(bid.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "click",
		AdType:                  string(req.raw.AdType),
		AdFormat:                string(bid.Format()),
		AuctionID:               bid.AuctionID,
		AuctionConfigurationID:  bid.AuctionConfigurationID,
		AuctionConfigurationUID: auctionConfigurationUID,
		Status:                  "",
		RoundID:                 bid.RoundID,
		RoundNumber:             bid.RoundIndex,
		ImpID:                   bid.ImpID,
		DemandID:                bid.DemandID,
		AdUnitUID:               int64(bid.GetAdUnitUID()),
		AdUnitLabel:             bid.AdUnitLabel,
		ECPM:                    bid.GetPrice(),
		PriceFloor:              bid.AuctionPriceFloor,
		Bidding:                 bid.IsBidding(),
	}

	return event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
}
