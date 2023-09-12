// Package adminecho implements Echo bindings for the admin package.
package adminecho

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/auth"
	v8n "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func UseAuthorization(g *echo.Group, authService *auth.Service) {
	g.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		// Do not check basic auth if JWT token is present.
		Skipper: authIsBearer,
		Validator: func(username, password string, c echo.Context) (bool, error) {
			if authService.IsSuperUser(username, password) {
				c.Set("authCtx", stubAuthContext{})

				return true, nil
			}

			return false, nil
		},
	}))
	g.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			// Skip if basic auth already set auth context.
			authCtx := c.Get("authCtx")
			return authCtx != nil
		},
		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*auth.JWTClaims)

			c.Set("authCtx", claims)
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTClaims)
		},
		KeyFunc: func(_ *jwt.Token) (any, error) {
			return authService.GetSecretKey(), nil
		},
	}))
}

func RegisterAuthService(g *echo.Group, service *auth.Service) {
	g.POST("/login", func(c echo.Context) error {
		var r auth.LogInRequest
		if err := c.Bind(&r); err != nil {
			return err
		}

		response, err := service.LogIn(c.Request().Context(), r)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			return err
		}

		return c.JSON(http.StatusOK, response)
	})
}

func RegisterAdminService(g *echo.Group, service *admin.Service) {
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

	g.GET("/permissions", getPermissionsHandler)
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

func getPermissionsHandler(c echo.Context) error {
	authContext := stubAuthContext{}
	permissions := admin.GetPermissions(authContext)

	return c.JSON(http.StatusOK, permissions)
}

type resourceServiceHandler[Resource, ResourceAttrs any] struct {
	service resourceService[Resource, ResourceAttrs]
}

type resourceService[Resource, ResourceAttrs any] interface {
	List(ctx context.Context, authCtx admin.AuthContext) ([]Resource, error)
	Find(ctx context.Context, authCtx admin.AuthContext, id int64) (*Resource, error)
	Create(ctx context.Context, attrs *ResourceAttrs) (*Resource, error)
	Update(ctx context.Context, id int64, attrs *ResourceAttrs) (*Resource, error)
	Delete(ctx context.Context, id int64) error
}

// stubAuthContext is a stub implementation of admin.AuthContext. Build auth context from JWT token.
type stubAuthContext struct{}

func (s stubAuthContext) UserID() int64 {
	return 0
}

func (s stubAuthContext) IsAdmin() bool {
	return true
}

func (s *resourceServiceHandler[Resource, ResourceAttrs]) list(c echo.Context) error {
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	resources, err := s.service.List(c.Request().Context(), authCtx)
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
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	resource, err := s.service.Find(c.Request().Context(), authCtx, int64(id))
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

func getAuthContext(c echo.Context) (admin.AuthContext, error) {
	authCtx, ok := c.Get("authCtx").(admin.AuthContext)
	if !ok {
		return nil, fmt.Errorf("failed to get auth context from request")
	}

	return authCtx, nil
}

func authIsBearer(c echo.Context) bool {
	header := c.Request().Header.Get(echo.HeaderAuthorization)

	const prefix = "Bearer "
	if len(header) < len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return false
	}

	return true
}
