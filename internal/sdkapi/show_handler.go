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

type ShowHandler struct {
	*BaseHandler[schema.ShowRequest, *schema.ShowRequest]
	EventLogger         *event.Logger
	NotificationHandler ShowNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/show_mocks.go -pkg mocks . ShowNotificationHandler
type ShowNotificationHandler interface {
	HandleShow(ctx context.Context, imp *schema.Imp, responses []*adapters.DemandResponse) error
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
