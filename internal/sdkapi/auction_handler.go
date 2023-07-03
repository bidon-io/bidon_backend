package sdkapi

import (
	"errors"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/labstack/echo/v4"
)

type AuctionHandler struct {
	*BaseHandler
	AuctionBuilder *auction.Builder
	SegmentMatcher *segment.Matcher
}

type AuctionResponse struct {
	*auction.Auction
	Token      string  `json:"token"`
	PriceFloor float64 `json:"pricefloor"`
	AuctionID  string  `json:"auction_id"`
}

func (h *AuctionHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	segmentParams := &segment.Params{
		Country: req.countryCode(),
		Ext:     req.raw.Segment.Ext,
		AppID:   req.app.ID,
	}

	sgmnt := h.SegmentMatcher.Match(c.Request().Context(), segmentParams)

	params := &auction.BuildParams{
		AppID:      req.app.ID,
		AdType:     req.raw.AdType,
		AdFormat:   req.adFormat(),
		DeviceType: req.raw.Device.Type,
		Adapters:   req.adapterKeys(),
		SegmentID:  sgmnt.ID,
	}
	auc, err := h.AuctionBuilder.Build(c.Request().Context(), params)
	if err != nil {
		if errors.Is(err, auction.ErrNoAdsFound) {
			err = ErrNoAdsFound
		}

		return err
	}

	response := &AuctionResponse{
		Auction:    auc,
		Token:      "{}",
		PriceFloor: req.raw.AdObject.PriceFloor,
		AuctionID:  req.raw.AdObject.AuctionID,
	}

	return c.JSON(http.StatusOK, response)
}
