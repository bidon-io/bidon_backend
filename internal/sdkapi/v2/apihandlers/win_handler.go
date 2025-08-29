package apihandlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type WinHandler struct {
	*BaseHandler[schema.WinRequest, *schema.WinRequest]
	EventLogger         *event.Logger
	NotificationHandler WinNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/win_mocks.go -pkg mocks . WinNotificationHandler
type WinNotificationHandler interface {
	HandleWin(ctx context.Context, bid *schema.Bid, config *auction.Config, bundle, adType string) error
}

func (h *WinHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	demandRequestEvent := prepareWinEvent(req)
	h.EventLogger.Log(demandRequestEvent, func(err error) {
		sdkapi.LogError(c, fmt.Errorf("log win event: %v", err))
	})

	// Call notification handler
	bid := req.raw.Bid
	err = h.NotificationHandler.HandleWin(c.Request().Context(), bid, req.auctionConfig, req.raw.App.Bundle, string(req.raw.AdType))
	if err != nil {
		sdkapi.LogError(c, fmt.Errorf("handle win notification: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func prepareWinEvent(req *request[schema.WinRequest, *schema.WinRequest]) *event.AdEvent {
	bid := req.raw.Bid

	auctionConfigurationUID, err := strconv.ParseInt(bid.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "win",
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
