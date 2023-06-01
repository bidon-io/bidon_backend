package admin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// A resourceService implements CRUD handlers.
// Resource and ResourceAttrs must be JSON encodable
type resourceService[Resource, ResourceAttrs any] struct {
	Repo ResourceRepo[Resource, ResourceAttrs]
}

// A ResourceRepo defines necessary interactions with a store
type ResourceRepo[Resource, ResourceAttrs any] interface {
	List(ctx context.Context) ([]Resource, error)
	Find(ctx context.Context, id int64) (*Resource, error)
	Create(ctx context.Context, attrs *ResourceAttrs) (*Resource, error)
	Update(ctx context.Context, id int64, attrs *ResourceAttrs) (*Resource, error)
	Delete(ctx context.Context, id int64) error
}

func (s *resourceService[Resource, ResourceAttrs]) list(c echo.Context) error {
	resources, err := s.Repo.List(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, resources, "  ")
}

func (s *resourceService[Resource, ResourceAttrs]) create(c echo.Context) error {
	attrs := new(ResourceAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	resource, err := s.Repo.Create(c.Request().Context(), attrs)
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusCreated, resource, "  ")
}

func (s *resourceService[Resource, ResourceAttrs]) get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	resource, err := s.Repo.Find(c.Request().Context(), int64(id))
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, resource, "  ")
}

func (s *resourceService[Resource, ResourceAttrs]) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	attrs := new(ResourceAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	resource, err := s.Repo.Update(c.Request().Context(), int64(id), attrs)
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, resource, "  ")
}

func (s *resourceService[Resource, ResourceAttrs]) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	if err := s.Repo.Delete(c.Request().Context(), int64(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
