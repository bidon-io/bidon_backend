package sdkapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type StatsV2Handler struct {
	*BaseHandler[schema.StatsV2Request, *schema.StatsV2Request]
	EventLogger         *event.Logger
	NotificationHandler StatsV2NotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/stats_v2_mocks.go -pkg mocks . StatsV2NotificationHandler

type StatsV2NotificationHandler interface {
	HandleStats(context.Context, schema.StatsV2, *auction.Config, string, string)
}

func (h *StatsV2Handler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	events := prepareStatsV2Events(req)
	for _, ev := range events {
		h.EventLogger.Log(ev, func(err error) {
			logError(c, fmt.Errorf("log %v event: %v", ev.EventType, err))
		})
	}

	h.NotificationHandler.HandleStats(c.Request().Context(), req.raw.Stats, req.auctionConfig, req.raw.App.Bundle, string(req.raw.AdType))

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func prepareStatsV2Events(req *request[schema.StatsV2Request, *schema.StatsV2Request]) []*event.AdEvent {
	// 1 event whole auction
	// 1 event for each Auction Ad Unit Result
	stats := req.raw.Stats

	auctionConfigurationUID, err := strconv.ParseInt(stats.AuctionConfigurationUID, 10, 64)
	if err != nil {
		auctionConfigurationUID = 0
	}

	// 1 event whole auction + 1 event for each ad unit
	events := make([]*event.AdEvent, 0, 1+len(stats.AdUnits))

	adRequestParams := event.AdRequestParams{
		EventType:               "stats_request",
		AdType:                  string(req.raw.AdType),
		AdFormat:                string(req.raw.Stats.Result.Format()),
		AuctionID:               stats.AuctionID,
		AuctionConfigurationID:  stats.AuctionConfigurationID,
		AuctionConfigurationUID: auctionConfigurationUID,
		Status:                  stats.Result.Status,
		ImpID:                   "",
		DemandID:                stats.Result.GetWinnerDemandID(),
		AdUnitUID:               int64(stats.Result.GetWinnerAdUnitUID()),
		AdUnitLabel:             stats.Result.WinnerAdUnitLabel,
		ECPM:                    stats.Result.GetWinnerPrice(),
		PriceFloor:              stats.AuctionPricefloor,
		TimingMap:               event.TimingMap{"auction": {stats.Result.AuctionStartTS, stats.Result.AuctionFinishTS}},
	}
	events = append(events, event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData))

	for _, adUnit := range stats.AdUnits {
		if adUnit.BidType == schema.CPMBidType {
			adRequestParams = event.AdRequestParams{
				EventType:               "demand_request",
				AdType:                  string(req.raw.AdType),
				AdFormat:                string(req.raw.Stats.Result.Format()),
				AuctionID:               stats.AuctionID,
				AuctionConfigurationID:  stats.AuctionConfigurationID,
				AuctionConfigurationUID: auctionConfigurationUID,
				Status:                  adUnit.Status,
				ImpID:                   "",
				DemandID:                adUnit.DemandID,
				AdUnitUID:               int64(adUnit.GetAdUnitUID()),
				AdUnitLabel:             adUnit.AdUnitLabel,
				ECPM:                    adUnit.GetPrice(),
				PriceFloor:              stats.AuctionPricefloor,
				Bidding:                 false,
				TimingMap:               event.TimingMap{"fill": {adUnit.FillStartTS, adUnit.FillFinishTS}},
			}
			events = append(events, event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData))
		} else {
			adRequestParams = event.AdRequestParams{
				EventType:               "client_bid",
				AdType:                  string(req.raw.AdType),
				AdFormat:                string(req.raw.Stats.Result.Format()),
				AuctionID:               stats.AuctionID,
				AuctionConfigurationID:  stats.AuctionConfigurationID,
				AuctionConfigurationUID: auctionConfigurationUID,
				Status:                  adUnit.Status,
				ImpID:                   "",
				DemandID:                adUnit.DemandID,
				AdUnitUID:               int64(adUnit.GetAdUnitUID()),
				AdUnitLabel:             adUnit.AdUnitLabel,
				ECPM:                    adUnit.GetPrice(),
				PriceFloor:              stats.AuctionPricefloor,
				Bidding:                 true,
				TimingMap: event.TimingMap{
					"fill":  {adUnit.FillStartTS, adUnit.FillFinishTS},
					"token": {adUnit.TokenStartTS, adUnit.TokenFinishTS},
				},
			}
			events = append(events, event.NewAdEvent(&req.raw.BaseRequest, adRequestParams, req.geoData))
		}
	}

	return events
}
