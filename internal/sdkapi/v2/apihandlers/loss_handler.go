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

type LossHandler struct {
	*BaseHandler[schema.LossRequest, *schema.LossRequest]
	EventLogger         *event.Logger
	NotificationHandler LossNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@v0.5.3 -out mocks/loss_mocks.go -pkg mocks . LossNotificationHandler
type LossNotificationHandler interface {
	HandleLoss(ctx context.Context, bid *schema.Bid, externalWinner *schema.ExternalWinner, config *auction.Config, bundle, adType string) error
}

func (h *LossHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	adEvent := prepareLossEvent(req)
	h.EventLogger.Log(adEvent, func(err error) {
		sdkapi.LogError(c, fmt.Errorf("log loss event: %v", err))
	})

	// Call notification handler
	bid := req.raw.Bid
	err = h.NotificationHandler.HandleLoss(c.Request().Context(), bid, &req.raw.ExternalWinner, req.auctionConfig, req.raw.App.Bundle, string(req.raw.AdType))
	if err != nil {
		sdkapi.LogError(c, fmt.Errorf("handle loss notification: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func prepareLossEvent(req *request[schema.LossRequest, *schema.LossRequest]) *event.AdEvent {
	bid := req.raw.Bid

	auctionConfigurationUID, err := strconv.ParseInt(bid.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "loss",
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
		ExternalWinnerDemandID:  req.raw.ExternalWinner.DemandID,
		ExternalWinnerEcpm:      req.raw.ExternalWinner.GetPrice(),
	}

	return event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData)
}
