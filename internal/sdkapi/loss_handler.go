package sdkapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
)

type LossHandler struct {
	*BaseHandler[schema.LossRequest, *schema.LossRequest]
	EventLogger         *event.Logger
	NotificationHandler LossNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/loss_mocks.go -pkg mocks . LossNotificationHandler
type LossNotificationHandler interface {
	HandleLoss(context.Context, *schema.Imp, []*adapters.DemandResponse) error
}

func (h *LossHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	lossEvent := event.NewLoss(&req.raw, req.geoData)
	h.EventLogger.Log(lossEvent, func(err error) {
		logError(c, fmt.Errorf("log loss event: %v", err))
	})

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
