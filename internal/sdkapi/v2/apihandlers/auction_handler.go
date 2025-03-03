package apihandlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/auctionv2"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type AuctionHandler struct {
	*BaseHandler[schema.AuctionV2Request, *schema.AuctionV2Request]
	AuctionService AuctionService
}

type AuctionService interface {
	Run(ctx context.Context, params *auctionv2.ExecutionParams) (*auctionv2.Response, error)
}

type AuctionResponse struct {
	ConfigID                 int64            `json:"auction_configuration_id"`
	ConfigUID                string           `json:"auction_configuration_uid"`
	ExternalWinNotifications bool             `json:"external_win_notifications"`
	AdUnits                  []auction.AdUnit `json:"ad_units"`
	NoBids                   []auction.AdUnit `json:"no_bids"`
	Segment                  auction.Segment  `json:"segment"`
	Token                    string           `json:"token"`
	AuctionPriceFloor        float64          `json:"auction_pricefloor"`
	AuctionTimeout           int              `json:"auction_timeout"`
	AuctionID                string           `json:"auction_id"`
}

func (h *AuctionHandler) Handle(c echo.Context) error {
	req, err := h.resolveRequest(c)
	if err != nil {
		return err
	}

	params := &auctionv2.ExecutionParams{
		Req:     &req.raw,
		AppID:   req.app.ID,
		Country: req.countryCode(),
		GeoData: req.geoData,
		Log: func(str string) {
			c.Logger().Printf(str)
		},
		LogErr: func(err error) {
			sdkapi.LogError(c, err)
		},
	}
	result, err := h.AuctionService.Run(c.Request().Context(), params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
