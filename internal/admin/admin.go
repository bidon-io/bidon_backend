// Package admin implements an HTTP API handlers for managing entities.
package admin

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Service struct {
	AuctionConfigurations *AuctionConfigurationService
	Apps                  *AppService
	Segments              *SegmentService
}

func (s *Service) RegisterAPIRoutes(g *echo.Group) {
	resources := []resourceRoutes{
		{"auction_configurations", s.AuctionConfigurations},
		{"apps", s.Apps},
		{"segments", s.Segments},
	}

	for i := range resources {
		resource := &resources[i]
		resource.register(g)
	}
}

type resourceRoutes struct {
	name    string
	handler resourceHandler
}

func (rr *resourceRoutes) register(g *echo.Group) {
	collectionPath := fmt.Sprintf("/%s", rr.name)
	itemPath := fmt.Sprintf("%s/:id", collectionPath)

	g.GET(collectionPath, rr.handler.list)
	g.POST(collectionPath, rr.handler.create)
	g.GET(itemPath, rr.handler.get)
	g.PATCH(itemPath, rr.handler.update)
	g.DELETE(itemPath, rr.handler.delete)
}

type resourceHandler interface {
	list(echo.Context) error
	create(echo.Context) error
	get(echo.Context) error
	update(echo.Context) error
	delete(echo.Context) error
}
