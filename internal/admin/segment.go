// Package admin implements an HTTP API handlers for managing entities.
package admin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Segment struct {
	ID int64 `json:"id"`
	SegmentAttrs
}

type SegmentAttrs struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Filters     []SegmentFilter `json:"filters"`
	Enabled     *bool           `json:"enabled"`
	AppID       int64           `json:"app_id"`
}

type SegmentFilter struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type SegmentRepo interface {
	List(ctx context.Context) ([]Segment, error)
	Find(ctx context.Context, id int64) (*Segment, error)
	Create(ctx context.Context, attrs *SegmentAttrs) (*Segment, error)
	Update(ctx context.Context, id int64, attrs *SegmentAttrs) (*Segment, error)
	Delete(ctx context.Context, id int64) error
}

func (s *Handlers) getSegments(c echo.Context) error {
	segments, err := s.SegmentRepo.List(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, segments, "  ")
}

func (s *Handlers) createSegment(c echo.Context) error {
	attrs := new(SegmentAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	segment, err := s.SegmentRepo.Create(c.Request().Context(), attrs)
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusCreated, segment, "  ")
}

func (s *Handlers) getSegment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	segment, err := s.SegmentRepo.Find(c.Request().Context(), int64(id))
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, segment, "  ")
}

func (s *Handlers) updateSegment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	attrs := new(SegmentAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	segment, err := s.SegmentRepo.Update(c.Request().Context(), int64(id), attrs)
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, segment, "  ")
}

func (s *Handlers) deleteSegment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	if err := s.SegmentRepo.Delete(c.Request().Context(), int64(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
