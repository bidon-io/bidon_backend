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

type ShowHandler struct {
	*BaseHandler[schema.ShowRequest, *schema.ShowRequest]
	EventLogger         *event.Logger
	NotificationHandler ShowNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/show_mocks.go -pkg mocks . ShowNotificationHandler
type ShowNotificationHandler interface {
	HandleShow(ctx context.Context, imp *schema.Imp, responses []*adapters.DemandResponse) error
}

func (h *ShowHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	h.sendEvents(c, req)

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func (h *ShowHandler) sendEvents(c echo.Context, req *request[schema.ShowRequest, *schema.ShowRequest]) {
	bid := req.raw.Bid

	auctionConfigurationUID, err := strconv.ParseInt(bid.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "show",
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
		logError(c, fmt.Errorf("log show event: %v", err))
	})
}
