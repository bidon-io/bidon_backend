package apihandlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/sdkapi"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type RewardHandler struct {
	*BaseHandler[schema.RewardRequest, *schema.RewardRequest]
	EventLogger *event.Logger
}

func (h *RewardHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	adEvent := prepareRewardEvent(req)
	h.EventLogger.Log(adEvent, func(err error) {
		sdkapi.LogError(c, fmt.Errorf("log reward event: %v", err))
	})

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func prepareRewardEvent(req *request[schema.RewardRequest, *schema.RewardRequest]) *event.AdEvent {
	bid := req.raw.Bid

	auctionConfigurationUID, err := strconv.ParseInt(bid.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "reward",
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

	return event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData)
}
