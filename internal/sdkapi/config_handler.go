package sdkapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/config"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	*BaseHandler[schema.ConfigRequest, *schema.ConfigRequest]
	AdaptersBuilder *config.AdaptersBuilder
	SegmentMatcher  *segment.Matcher
	EventLogger     *event.Logger
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
	ctx := c.Request().Context()

	event, err := event.Prepare(event.ConfigTopic, &req.raw, req.geoData)
	if err != nil {
		logError(c, fmt.Errorf("prepare config event: %v", err))
	}
	h.EventLogger.Log(ctx, event, func(err error) {
		logError(c, fmt.Errorf("log config event: %v", err))
	})

	segmentParams := &segment.Params{
		Country: req.countryCode(),
		Ext:     req.raw.Segment.Ext,
	}

	sgmnt := h.SegmentMatcher.Match(ctx, segmentParams)

	var segmentID string
	if sgmnt.ID != 0 {
		segmentID = strconv.Itoa(int(sgmnt.ID))
	}

	adapters, err := h.AdaptersBuilder.Build(ctx, req.app.ID, req.raw.Adapters.Keys())
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
