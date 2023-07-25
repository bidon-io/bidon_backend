package sdkapi

import (
	"fmt"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type ClickHandler struct {
	*BaseHandler[schema.ClickRequest, *schema.ClickRequest]
	EventLogger *event.Logger
}

func (h *ClickHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	clickEvent := event.NewClick(&req.raw, req.geoData)
	h.EventLogger.Log(clickEvent, func(err error) {
		logError(c, fmt.Errorf("log click event: %v", err))
	})

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
