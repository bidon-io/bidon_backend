package sdkapi

import (
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/config"
	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	*BaseHandler
	AdaptersBuilder *config.AdaptersBuilder
}

type ConfigResponse struct {
	Init       ConfigResponseInit `json:"init"`
	Placements []any              `json:"placements"`
	Token      string             `json:"token"`
	SegmentID  string             `json:"segment_id"`
}

type ConfigResponseInit struct {
	TMax     int             `json:"tmax"`
	Adapters config.Adapters `json:"adapters"`
}

func (h *ConfigHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	adapters, err := h.AdaptersBuilder.Build(c.Request().Context(), req.app.ID, req.adapterKeys())
	if err != nil {
		return err
	}
	if len(adapters) == 0 {
		return ErrNoAdaptersFound
	}

	resp := &ConfigResponse{
		Init: ConfigResponseInit{
			TMax:     5000,
			Adapters: adapters,
		},
		Placements: []any{},
		Token:      "{}",
		SegmentID:  "",
	}

	return c.JSON(http.StatusOK, resp)
}
