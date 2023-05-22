// Package admin implements an HTTP API handlers for managing entities.
package admin

import (
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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
	configuration := new(auction.Configuration)
	if err := c.Bind(configuration); err != nil {
		return err
	}

	if err := s.AuctionConfigurationRepo.Create(c.Request().Context(), configuration); err != nil {
		return err
	}

	return c.JSONPretty(http.StatusCreated, configuration, "  ")
}

func (s *Handlers) getAuctionConfiguration(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	configuration, err := s.AuctionConfigurationRepo.Find(c.Request().Context(), uint(id))
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

	configuration := new(auction.Configuration)
	if err := c.Bind(configuration); err != nil {
		return err
	}

	configuration.ID = uint(id)
	if err := s.AuctionConfigurationRepo.Update(c.Request().Context(), configuration); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Handlers) deleteAuctionConfiguration(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	if err := s.AuctionConfigurationRepo.Delete(c.Request().Context(), uint(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
