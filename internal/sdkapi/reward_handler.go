package sdkapi

import (
	"fmt"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type RewardHandler struct {
	*BaseHandler[schema.RewardRequest, *schema.RewardRequest]
	EventLogger *event.Logger
}

func (h *RewardHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	rewardEvent := event.NewReward(&req.raw, req.geoData)
	h.EventLogger.Log(rewardEvent, func(err error) {
		logError(c, fmt.Errorf("log reward event: %v", err))
	})

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
