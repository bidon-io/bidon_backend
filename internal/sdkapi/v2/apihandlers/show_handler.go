package apihandlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type ShowHandler struct {
	*BaseHandler[schema.ShowRequest, *schema.ShowRequest]
	EventLogger         *event.Logger
	NotificationHandler ShowNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/show_mocks.go -pkg mocks . ShowNotificationHandler
type ShowNotificationHandler interface {
	HandleShow(context.Context, *schema.Bid, string, string)
}

func (h *ShowHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	demandRequestEvent := prepareShowEvent(req)
	h.EventLogger.Log(demandRequestEvent, func(err error) {
		sdkapi.LogError(c, fmt.Errorf("log show event: %v", err))
	})

	h.NotificationHandler.HandleShow(c.Request().Context(), req.raw.Bid, req.raw.App.Bundle, string(req.raw.AdType))

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func prepareShowEvent(req *request[schema.ShowRequest, *schema.ShowRequest]) *event.AdEvent {
	bid := req.raw.Bid

	auctionConfigurationUID, err := strconv.ParseInt(bid.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "show",
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
