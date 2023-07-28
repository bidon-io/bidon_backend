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

type WinHandler struct {
	*BaseHandler[schema.WinRequest, *schema.WinRequest]
	EventLogger         *event.Logger
	NotificationHandler WinNotificationHandler
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/win_mocks.go -pkg mocks . WinNotificationHandler
type WinNotificationHandler interface {
	HandleWin(context.Context, *schema.Imp, []*adapters.DemandResponse) error
}

func (h *WinHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	winEvent := event.NewWin(&req.raw, req.geoData)
	h.EventLogger.Log(winEvent, func(err error) {
		logError(c, fmt.Errorf("log win event: %v", err))
	})

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
