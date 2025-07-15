package apihandlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type AuctionHandler struct {
	*BaseHandler[schema.AuctionRequest, *schema.AuctionRequest]
	AuctionService AuctionService
}

type AuctionService interface {
	Run(ctx context.Context, params *auction.ExecutionParams) (*auction.Response, error)
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

	if h.shouldReturnErrNoAdsFound(&req.raw) {
		return sdkapi.ErrNoAdsFound
	}

	params := &auction.ExecutionParams{
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

// shouldReturnEmptyResponse checks if the request matches conditions for returning ErrNoAdsFound:
// - OS is NOT android
// - Mediator is max (BCAMAX case)
// - SDK version is 0.7.x or 0.8.1
// - ad_type = rewarded_video
func (h *AuctionHandler) shouldReturnErrNoAdsFound(req *schema.AuctionRequest) bool {
	if req.Device.OS == "android" {
		return false
	}

	if req.GetMediator() != "max" {
		return false
	}

	if req.AdType != ad.RewardedType {
		return false
	}

	sdkVersion, err := req.GetSDKVersionSemver()
	if err != nil {
		return false
	}

	return sdkapi.Version07xConstraint.Check(sdkVersion) || sdkapi.Version081Constraint.Check(sdkVersion)
}
