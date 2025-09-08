//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml --config=config.yaml ../openapi/openapi.yaml

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	Handle(c echo.Context) error
}

//go:generate go run -mod=mod github.com/matryer/moq@v0.5.3 -out mocks/mocks.go -pkg mocks . Handler

type Server struct {
	AuctionHandler Handler
	ClickHandler   Handler
	ConfigHandler  Handler
	LossHandler    Handler
	StatsHandler   Handler
	ShowHandler    Handler
	RewardHandler  Handler
	WinHandler     Handler
}

func (s *Server) GetAuction(c echo.Context, _ GetAuctionParamsAdType, _ GetAuctionParams) error {
	return s.AuctionHandler.Handle(c)
}

func (s *Server) GetConfig(c echo.Context, _ GetConfigParams) error {
	return s.ConfigHandler.Handle(c)
}

func (s *Server) PostClick(c echo.Context, _ PostClickParamsAdType, _ PostClickParams) error {
	return s.ClickHandler.Handle(c)
}

func (s *Server) PostLoss(c echo.Context, _ PostLossParamsAdType, _ PostLossParams) error {
	return s.LossHandler.Handle(c)
}

func (s *Server) PostStats(c echo.Context, _ PostStatsParamsAdType, _ PostStatsParams) error {
	return s.StatsHandler.Handle(c)
}

func (s *Server) PostShow(c echo.Context, _ PostShowParamsAdType, _ PostShowParams) error {
	return s.ShowHandler.Handle(c)
}

func (s *Server) PostReward(c echo.Context, _ PostRewardParamsAdType, _ PostRewardParams) error {
	return s.RewardHandler.Handle(c)
}

func (s *Server) PostWin(c echo.Context, _ PostWinParamsAdType, _ PostWinParams) error {
	return s.WinHandler.Handle(c)
}

func (s *Server) GetOpenAPISpec(c echo.Context) error {
	spec, err := GetSwagger()
	if err != nil {
		return err
	}

	swaggerJSON, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to generate OpenAPI spec")
	}

	return c.JSONBlob(http.StatusOK, swaggerJSON)
}

// Ensure that we implement the server interface
var _ ServerInterface = (*Server)(nil)
