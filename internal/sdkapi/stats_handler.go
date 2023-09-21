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
	ConfigMatcher       ConfigMatcher
	EventLogger         *event.Logger
	NotificationHandler StatsNotificationHandler
}

type StatsNotificationHandler interface {
	HandleStats(context.Context, schema.Stats, auction.Config) error
}

type ConfigMatcher interface {
	MatchById(ctx context.Context, appID, id int64) *auction.Config
}

func (h *StatsHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	configEvent := event.NewStats(&req.raw, req.geoData)
	h.EventLogger.Log(configEvent, func(err error) {
		logError(c, fmt.Errorf("log stats event: %v", err))
	})

	h.sendEvents(c, req)

	ctx := c.Request().Context()
	config := h.ConfigMatcher.MatchById(ctx, req.app.ID, int64(req.raw.Stats.AuctionConfigurationID))
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

	auctionConfigurationUID, err := strconv.Atoi(stats.AuctionConfigurationUID)
	if err != nil {
		auctionConfigurationUID = 0
	}

	adRequestParams := event.AdRequestParams{
		EventType:               "stats_request",
		AdType:                  string(req.raw.AdType),
		AuctionID:               stats.AuctionID,
		AuctionConfigurationID:  int64(stats.AuctionConfigurationID),
		AuctionConfigurationUID: int64(auctionConfigurationUID),
		Status:                  stats.Result.Status,
		RoundID:                 stats.Result.RoundID,
		RoundNumber:             0,
		ImpID:                   "",
		DemandID:                stats.Result.WinnerID,
		AdUnitID:                0,
		LineItemUID:             0,
		AdUnitCode:              "",
		Ecpm:                    stats.Result.ECPM,
		PriceFloor:              0,
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
			AuctionConfigurationID:  int64(stats.AuctionConfigurationID),
			AuctionConfigurationUID: int64(auctionConfigurationUID),
			RoundID:                 round.ID,
			RoundNumber:             roundNumber,
			ImpID:                   "",
			DemandID:                round.WinnerID,
			AdUnitID:                0,
			LineItemUID:             0,
			AdUnitCode:              "",
			Ecpm:                    round.WinnerECPM,
			PriceFloor:              round.PriceFloor,
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
			lineItemUID, err := strconv.Atoi(demand.LineItemUID)
			if err != nil {
				lineItemUID = 0
			}

			adRequestParams = event.AdRequestParams{
				EventType:               "demand_request",
				AdType:                  string(req.raw.AdType),
				AuctionID:               stats.AuctionID,
				AuctionConfigurationID:  int64(stats.AuctionConfigurationID),
				AuctionConfigurationUID: int64(auctionConfigurationUID),
				Status:                  demand.Status,
				RoundID:                 round.ID,
				RoundNumber:             roundNumber,
				ImpID:                   "",
				DemandID:                demand.ID,
				AdUnitID:                0,
				LineItemUID:             int64(lineItemUID),
				AdUnitCode:              demand.AdUnitID,
				Ecpm:                    demand.ECPM,
				PriceFloor:              round.PriceFloor,
				Bidding:                 false,
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
				AuctionConfigurationID:  int64(stats.AuctionConfigurationID),
				AuctionConfigurationUID: int64(auctionConfigurationUID),
				Status:                  bid.Status,
				RoundID:                 round.ID,
				RoundNumber:             roundNumber,
				ImpID:                   "",
				DemandID:                bid.ID,
				AdUnitID:                0,
				LineItemUID:             0,
				AdUnitCode:              "",
				Ecpm:                    bid.ECPM,
				PriceFloor:              round.PriceFloor,
				Bidding:                 true,
			}
			demandRequestEvent := event.NewRequest(&req.raw.BaseRequest, adRequestParams, req.geoData)
			h.EventLogger.Log(demandRequestEvent, func(err error) {
				logError(c, fmt.Errorf("log bid event: %v", err))
			})
		}
	}
}
