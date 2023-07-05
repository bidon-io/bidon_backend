package sdkapi

import (
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/config"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	*BaseHandler[schema.ConfigRequest, *schema.ConfigRequest]
	AdaptersBuilder *config.AdaptersBuilder
	SegmentMatcher  *segment.Matcher
}

type ConfigResponse struct {
	Init       ConfigResponseInit `json:"init"`
	Placements []any              `json:"placements"`
	Token      string             `json:"token"`
	Segment    Segment            `json:"segment"`
}

type Segment struct {
	ID string `json:"id"`
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

	segmentParams := &segment.Params{
		Country: req.countryCode(),
		Ext:     req.raw.Segment.Ext,
	}

	sgmnt := h.SegmentMatcher.Match(c.Request().Context(), segmentParams)

	var segmentID string
	if sgmnt.ID != 0 {
		segmentID = strconv.Itoa(int(sgmnt.ID))
	} else {
		segmentID = ""
	}

	adapters, err := h.AdaptersBuilder.Build(c.Request().Context(), req.app.ID, req.raw.Adapters.Keys())
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
		Segment:    Segment{ID: segmentID},
	}

	return c.JSON(http.StatusOK, resp)
}
