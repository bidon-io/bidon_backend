// Package admin implements an HTTP API handlers for managing entities.
package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	AuctionConfigurationRepo auction.ConfigurationRepo
}

func (s *Handlers) RegisterRoutes(e *echo.Echo) {
	e.GET("/auction_configurations", s.getAuctionConfigurations)
	e.POST("/auction_configurations", s.createAuctionConfiguration)
	e.GET("/auction_configurations/:id", s.getAuctionConfiguration)
	e.PUT("/auction_configurations/:id", s.updateAuctionConfiguration)
	e.DELETE("/auction_configurations/:id", s.deleteAuctionConfiguration)
}

func (s *Handlers) getAuctionConfigurations(c echo.Context) error {
	configurations, err := s.AuctionConfigurationRepo.List(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, configurations, "  ")
}

func (s *Handlers) createAuctionConfiguration(c echo.Context) error {
	attrs := new(auction.ConfigurationAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	config, err := s.AuctionConfigurationRepo.Create(c.Request().Context(), attrs)
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusCreated, config, "  ")
}

func (s *Handlers) getAuctionConfiguration(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	configuration, err := s.AuctionConfigurationRepo.Find(c.Request().Context(), int64(id))
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, configuration, "  ")
}

func (s *Handlers) updateAuctionConfiguration(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	attrs := new(auction.ConfigurationAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	config, err := s.AuctionConfigurationRepo.Update(c.Request().Context(), int64(id), attrs)
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, config, "  ")
}

func (s *Handlers) deleteAuctionConfiguration(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	if err := s.AuctionConfigurationRepo.Delete(c.Request().Context(), int64(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
