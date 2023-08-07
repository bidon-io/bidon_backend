package sdkapi

import (
	"context"
	"fmt"
	"net/http"

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
