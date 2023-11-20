package sdkapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type WinHandler struct {
	*BaseHandler[schema.WinRequest, *schema.WinRequest]
	EventLogger         *event.Logger
	NotificationHandler WinNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/win_mocks.go -pkg mocks . WinNotificationHandler
type WinNotificationHandler interface {
	HandleWin(context.Context, *schema.Imp, []*adapters.DemandResponse) error
}

func (h *WinHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	h.sendEvents(c, req)

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func (h *WinHandler) sendEvents(c echo.Context, req *request[schema.WinRequest, *schema.WinRequest]) {
	bid := req.raw.Bid

	auctionConfigurationUID, err := strconv.ParseInt(bid.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "win",
		AdType:                  string(req.raw.AdType),
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
	demandRequestEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
	h.EventLogger.Log(demandRequestEvent, func(err error) {
		logError(c, fmt.Errorf("log win event: %v", err))
	})
}
