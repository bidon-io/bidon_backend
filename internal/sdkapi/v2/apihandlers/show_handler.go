package apihandlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/util"
)

type ShowHandler struct {
	*BaseHandler[schema.ShowRequest, *schema.ShowRequest]
	EventLogger         *event.Logger
	NotificationHandler ShowNotificationHandler
	AdUnitLookup        AdUnitLookup
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/show_mocks.go -pkg mocks . ShowNotificationHandler AdUnitLookup
type ShowNotificationHandler interface {
	HandleShow(context.Context, *schema.Bid, string, string)
}

type AdUnitLookup interface {
	GetByUIDCached(context.Context, string) (*db.LineItem, error)
}

func (h *ShowHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	demandRequestEvent := prepareShowEvent(c.Request().Context(), req, h.AdUnitLookup)
	h.EventLogger.Log(demandRequestEvent, func(err error) {
		sdkapi.LogError(c, fmt.Errorf("log show event: %v", err))
	})

	h.NotificationHandler.HandleShow(c.Request().Context(), req.raw.Bid, req.raw.App.Bundle, string(req.raw.AdType))

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func prepareShowEvent(ctx context.Context, req *request[schema.ShowRequest, *schema.ShowRequest], adUnitLookup AdUnitLookup) *event.AdEvent {
	bid := req.raw.Bid

	auctionConfigurationUID, err := strconv.ParseInt(bid.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	var adUnitInternalID int64
	var adUnitCredentials map[string]string

	adUnit, err := adUnitLookup.GetByUIDCached(ctx, bid.AdUnitUID)
	if err == nil && adUnit != nil {
		adUnitInternalID = adUnit.ID
		adUnitCredentials = util.ConvertToStringMap(adUnit.Extra)
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
		AdUnitInternalID:        adUnitInternalID,
		AdUnitLabel:             bid.AdUnitLabel,
		AdUnitCredentials:       adUnitCredentials,
		ECPM:                    bid.GetPrice(),
		PriceFloor:              bid.AuctionPriceFloor,
		Bidding:                 bid.IsBidding(),
	}

	return event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData)
}
