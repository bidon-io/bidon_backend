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

type LossHandler struct {
	*BaseHandler[schema.LossRequest, *schema.LossRequest]
	EventLogger         *event.Logger
	NotificationHandler LossNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/loss_mocks.go -pkg mocks . LossNotificationHandler
type LossNotificationHandler interface {
	HandleLoss(context.Context, *schema.Imp, []*adapters.DemandResponse) error
}

func (h *LossHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	h.sendEvents(c, req)

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func (h *LossHandler) sendEvents(c echo.Context, req *request[schema.LossRequest, *schema.LossRequest]) {
	bid := req.raw.Bid

	auctionConfigurationUID, err := strconv.ParseInt(bid.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "loss",
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
		ExternalWinnerDemandID:  req.raw.ExternalWinner.DemandID,
		ExternalWinnerEcpm:      req.raw.ExternalWinner.ECPM,
	}
	adEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
	h.EventLogger.Log(adEvent, func(err error) {
		logError(c, fmt.Errorf("log loss event: %v", err))
	})
}
