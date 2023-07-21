package sdkapi

import (
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
	"net/http"
)

type StatsHandler struct {
	*BaseHandler[schema.StatsRequest, *schema.StatsRequest]
	EventLogger *event.Logger
}

func (h *StatsHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	configEvent := event.NewStats(&req.raw, req.geoData)
	h.EventLogger.Log(configEvent, func(err error) {
		logError(c, fmt.Errorf("log config event: %v", err))
	})

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
