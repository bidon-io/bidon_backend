//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml --config=config.yaml ../openapi/openapi.yaml

package api

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler interface {
	Handle(c echo.Context) error
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks/mocks.go -pkg mocks . Handler

type Server struct {
	AuctionHandler Handler
	ConfigHandler  Handler
}

func (s *Server) GetAuction(c echo.Context, _ GetAuctionParamsAdType) error {
	return s.AuctionHandler.Handle(c)
}

func (s *Server) GetConfig(c echo.Context) error {
	return s.ConfigHandler.Handle(c)
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