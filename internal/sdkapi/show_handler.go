package sdkapi

import (
	"fmt"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type ShowHandler struct {
	*BaseHandler[schema.ShowRequest, *schema.ShowRequest]
	EventLogger *event.Logger
}

func (h *ShowHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	showEvent := event.NewShow(&req.raw, req.geoData)
	h.EventLogger.Log(showEvent, func(err error) {
		logError(c, fmt.Errorf("log show event: %v", err))
	})

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
