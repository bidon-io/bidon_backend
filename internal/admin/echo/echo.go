// Package adminecho implements Echo bindings for the admin package.
package adminecho

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bidon-io/bidon-backend/internal/admin"
	v8n "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

func RegisterService(g *echo.Group, service *admin.Service) {
	resourceRoutes := []resourceRoute{
		{
			group:   g.Group("/apps"),
			handler: &appHandler{service.AppService},
		},
		{
			group:   g.Group("/app_demand_profiles"),
			handler: &appDemandProfileHandler{service.AppDemandProfileService},
		},
		{
			group:   g.Group("/auction_configurations"),
			handler: &auctionConfigurationHandler{service.AuctionConfigurationService},
		},
		{
			group:   g.Group("/countries"),
			handler: &countryHandler{service.CountryService},
		},
		{
			group:   g.Group("/demand_sources"),
			handler: &demandSourceHandler{service.DemandSourceService},
		},
		{
			group:   g.Group("/demand_source_accounts"),
			handler: &demandSourceAccountHandler{service.DemandSourceAccountService},
		},
		{
			group:   g.Group("/line_items"),
			handler: &lineItemHandler{service.LineItemService},
		},
		{
			group:   g.Group("/segments"),
			handler: &segmentHandler{service.SegmentService},
		},
		{
			group:   g.Group("/users"),
			handler: &userHandler{service.UserService},
		},
	}
	for _, r := range resourceRoutes {
		r.group.GET("", r.handler.list)
		r.group.POST("", r.handler.create)
		r.group.GET("/:id", r.handler.get)
		r.group.PUT("/:id", r.handler.update)
		r.group.PATCH("/:id", r.handler.update)
		r.group.DELETE("/:id", r.handler.delete)
	}
}

type resourceRoute struct {
	group   *echo.Group
	handler resourceHandler
}

type resourceHandler interface {
	list(c echo.Context) error
	create(c echo.Context) error
	get(c echo.Context) error
	update(c echo.Context) error
	delete(c echo.Context) error
}

type appHandler = resourceServiceHandler[admin.App, admin.AppAttrs]
type appDemandProfileHandler = resourceServiceHandler[admin.AppDemandProfile, admin.AppDemandProfileAttrs]
type auctionConfigurationHandler = resourceServiceHandler[admin.AuctionConfiguration, admin.AuctionConfigurationAttrs]
type countryHandler = resourceServiceHandler[admin.Country, admin.CountryAttrs]
type demandSourceHandler = resourceServiceHandler[admin.DemandSource, admin.DemandSourceAttrs]
type demandSourceAccountHandler = resourceServiceHandler[admin.DemandSourceAccount, admin.DemandSourceAccountAttrs]
type lineItemHandler = resourceServiceHandler[admin.LineItem, admin.LineItemAttrs]
type segmentHandler = resourceServiceHandler[admin.Segment, admin.SegmentAttrs]
type userHandler = resourceServiceHandler[admin.User, admin.UserAttrs]

type resourceServiceHandler[Resource, ResourceAttrs any] struct {
	service resourceService[Resource, ResourceAttrs]
}

type resourceService[Resource, ResourceAttrs any] interface {
	List(ctx context.Context) ([]Resource, error)
	Find(ctx context.Context, id int64) (*Resource, error)
	Create(ctx context.Context, attrs *ResourceAttrs) (*Resource, error)
	Update(ctx context.Context, id int64, attrs *ResourceAttrs) (*Resource, error)
	Delete(ctx context.Context, id int64) error
}

func (s *resourceServiceHandler[Resource, ResourceAttrs]) list(c echo.Context) error {
	resources, err := s.service.List(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resources)
}

func (s *resourceServiceHandler[Resource, ResourceAttrs]) create(c echo.Context) error {
	attrs := new(ResourceAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	resource, err := s.service.Create(c.Request().Context(), attrs)
	if err != nil {
		var validationError v8n.Errors
		if errors.As(err, &validationError) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, validationError.Error())
		}

		return err
	}

	return c.JSON(http.StatusCreated, resource)
}

func (s *resourceServiceHandler[Resource, ResourceAttrs]) get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	resource, err := s.service.Find(c.Request().Context(), int64(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resource)
}

func (s *resourceServiceHandler[Resource, ResourceAttrs]) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	attrs := new(ResourceAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	resource, err := s.service.Update(c.Request().Context(), int64(id), attrs)
	if err != nil {
		var validationError v8n.Errors
		if errors.As(err, &validationError) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, validationError.Error())
		}

		return err
	}

	return c.JSON(http.StatusOK, resource)
}

func (s *resourceServiceHandler[Resource, ResourceAttrs]) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	if err := s.service.Delete(c.Request().Context(), int64(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
