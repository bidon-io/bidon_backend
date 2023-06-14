package sdkapi

import (
	"context"
	"net/http"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
)

type Service struct {
	AuctionBuilder *auction.Builder
	AppFetcher     AppFetcher
}

// App represents an app for the purposes of the SDK API
type App struct {
	ID int64
}

type AppFetcher interface {
	Fetch(ctx context.Context, appKey, appBundle string) (*App, error)
}

type AuctionResponse struct {
	auction.Auction
	Token      string  `json:"token"`
	PriceFloor float64 `json:"pricefloor"`
	AuctionID  string  `json:"auction_id"`
}

func (s *Service) HandleAuction(ctx echo.Context) error {
	var request schema.Request
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	requestCtx := ctx.Request().Context()

	app, err := s.AppFetcher.Fetch(requestCtx, request.App.Key, request.App.Bundle)
	if err != nil {
		return err
	}

	params := &auction.BuildParams{
		AppID:      app.ID,
		AdType:     request.AdType,
		AdFormat:   request.AdObject.AdFormat(),
		DeviceType: request.Device.Type,
		Adapters:   maps.Keys(request.Adapters),
	}
	auction, err := s.AuctionBuilder.Build(ctx.Request().Context(), params)
	if err != nil {
		return err
	}

	response := &AuctionResponse{
		Auction:    *auction,
		Token:      "{}",
		PriceFloor: request.AdObject.PriceFloor,
		AuctionID:  request.AdObject.AuctionID,
	}

	return ctx.JSON(http.StatusOK, response)
}
