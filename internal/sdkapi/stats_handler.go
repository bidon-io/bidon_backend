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

type StatsHandler struct {
	*BaseHandler[schema.StatsRequest, *schema.StatsRequest]
	EventLogger         *event.Logger
	NotificationHandler StatsNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/stats_mocks.go -pkg mocks . StatsNotificationHandler

type StatsNotificationHandler interface {
	HandleStats(context.Context, schema.Stats, auction.Config) error
}

func (h *StatsHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	h.sendEvents(c, req)

	config := req.auctionConfig
	if config == nil {
		logError(c, fmt.Errorf("cannot find config: %v", req.raw.Stats.AuctionConfigurationID))
	} else {
		_ = ""
		// h.NotificationHandler.HandleStats(ctx, req.raw.Stats, *config)
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func (h *StatsHandler) sendEvents(c echo.Context, req *request[schema.StatsRequest, *schema.StatsRequest]) {
	// 1 event whole auction
	// 1 event for each round
	// 1 event for each demand in round
	stats := req.raw.Stats

	if req.auctionConfig == nil {
		err := fmt.Errorf(
			"no config found for app_id:%d, id: %d, uid: %s",
			req.app.ID, stats.AuctionConfigurationID, stats.AuctionConfigurationUID,
		)
		logError(c, err)
	}

	auctionConfigurationUID, err := strconv.Atoi(stats.AuctionConfigurationUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	// find round by ID
	statsPriceFloor := 0.0
	statsRoundNumber := 0
	for idx, round := range stats.Rounds {
		if round.ID == stats.Result.RoundID {
			statsPriceFloor = round.PriceFloor
			statsRoundNumber = idx
			break
		}
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "stats_request",
		AdType:                  string(req.raw.AdType),
		AuctionID:               stats.AuctionID,
		AuctionConfigurationID:  stats.AuctionConfigurationID,
		AuctionConfigurationUID: int64(auctionConfigurationUID),
		Status:                  stats.Result.Status,
		RoundID:                 stats.Result.RoundID,
		RoundNumber:             statsRoundNumber,
		ImpID:                   "",
		DemandID:                stats.Result.GetWinnerDemandID(),
		AdUnitUID:               int64(stats.Result.GetWinnerAdUnitUID()),
		AdUnitLabel:             stats.Result.WinnerAdUnitLabel,
		Ecpm:                    stats.Result.GetWinnerPrice(),
		PriceFloor:              statsPriceFloor,
		TimingMap:               event.TimingMap{"auction": {stats.Result.AuctionStartTS, stats.Result.AuctionFinishTS}},
	}
	statsRequestEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
	h.EventLogger.Log(statsRequestEvent, func(err error) {
		logError(c, fmt.Errorf("log stats_request event: %v", err))
	})

	for roundNumber, round := range stats.Rounds {
		adRequestParams = event.AdRequestParams{
			EventType:               "round_request",
			AdType:                  string(req.raw.AdType),
			AuctionID:               stats.AuctionID,
			AuctionConfigurationID:  stats.AuctionConfigurationID,
			AuctionConfigurationUID: int64(auctionConfigurationUID),
			RoundID:                 round.ID,
			RoundNumber:             roundNumber,
			ImpID:                   "",
			DemandID:                round.GetWinnerDemandID(),
			AdUnitUID:               int64(round.GetWinnerAdUnitUID()),
			AdUnitLabel:             round.WinnerAdUnitLabel,
			Ecpm:                    round.GetWinnerPrice(),
			PriceFloor:              round.PriceFloor,
			TimingMap:               event.TimingMap{"auction": {stats.Result.AuctionStartTS, stats.Result.AuctionFinishTS}},
		}
		if round.WinnerID != "" {
			adRequestParams.Status = "SUCCESS"
		} else {
			adRequestParams.Status = "FAIL"
		}
		roundRequestEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
		h.EventLogger.Log(roundRequestEvent, func(err error) {
			logError(c, fmt.Errorf("log round_request event: %v", err))
		})

		for _, demand := range round.Demands {
			adRequestParams = event.AdRequestParams{
				EventType:               "demand_request",
				AdType:                  string(req.raw.AdType),
				AuctionID:               stats.AuctionID,
				AuctionConfigurationID:  stats.AuctionConfigurationID,
				AuctionConfigurationUID: int64(auctionConfigurationUID),
				Status:                  demand.Status,
				RoundID:                 round.ID,
				RoundNumber:             roundNumber,
				ImpID:                   "",
				DemandID:                demand.ID,
				AdUnitUID:               int64(demand.GetAdUnitUID()),
				AdUnitLabel:             demand.AdUnitLabel,
				Ecpm:                    demand.GetPrice(),
				PriceFloor:              round.PriceFloor,
				Bidding:                 false,
				TimingMap:               event.TimingMap{"fill": {demand.FillStartTS, demand.FillFinishTS}},
			}
			demandRequestEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
			h.EventLogger.Log(demandRequestEvent, func(err error) {
				logError(c, fmt.Errorf("log demand_request event: %v", err))
			})
		}

		for _, bid := range round.Bidding.Bids {
			adRequestParams = event.AdRequestParams{
				EventType:               "client_bid",
				AdType:                  string(req.raw.AdType),
				AuctionID:               stats.AuctionID,
				AuctionConfigurationID:  stats.AuctionConfigurationID,
				AuctionConfigurationUID: int64(auctionConfigurationUID),
				Status:                  bid.Status,
				RoundID:                 round.ID,
				RoundNumber:             roundNumber,
				ImpID:                   "",
				DemandID:                bid.ID,
				AdUnitUID:               int64(bid.GetAdUnitUID()),
				AdUnitLabel:             bid.AdUnitLabel,
				Ecpm:                    bid.GetPrice(),
				PriceFloor:              round.PriceFloor,
				Bidding:                 true,
				TimingMap: event.TimingMap{
					"fill":  {bid.FillStartTS, bid.FillFinishTS},
					"token": {bid.TokenStartTS, bid.TokenFinishTS},
					"bid":   {round.Bidding.BidStartTS, round.Bidding.BidFinishTS},
				},
			}
			demandRequestEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
			h.EventLogger.Log(demandRequestEvent, func(err error) {
				logError(c, fmt.Errorf("log bid event: %v", err))
			})
		}
	}
}
