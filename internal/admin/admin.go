// Package admin implements an HTTP API handlers for managing entities.
package admin

import (
	"github.com/labstack/echo/v4"
)

type Service struct {
	AuctionConfigurations *AuctionConfigurationService
	Apps                  *AppService
	Segments              *SegmentService
}

func (s *Service) RegisterAPIRoutes(g *echo.Group) {
	g.GET("/auction_configurations", s.AuctionConfigurations.list)
	g.POST("/auction_configurations", s.AuctionConfigurations.create)
	g.GET("/auction_configurations/:id", s.AuctionConfigurations.get)
	g.PATCH("/auction_configurations/:id", s.AuctionConfigurations.update)
	g.DELETE("/auction_configurations/:id", s.AuctionConfigurations.delete)

	g.GET("/apps", s.Apps.list)
	g.POST("/apps", s.Apps.create)
	g.GET("/apps/:id", s.Apps.get)
	g.PATCH("/apps/:id", s.Apps.update)
	g.DELETE("/apps/:id", s.Apps.delete)

	g.GET("/segments", s.Segments.list)
	g.POST("/segments", s.Segments.create)
	g.GET("/segments/:id", s.Segments.get)
	g.PATCH("/segments/:id", s.Segments.update)
	g.DELETE("/segments/:id", s.Segments.delete)
}
