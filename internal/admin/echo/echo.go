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
	session "github.com/spazzymoto/echo-scs-session"
)

func UseAuthorization(g *echo.Group, authService *auth.Service) {
	sm := authService.GetSessionManager()
	g.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Skipper: skipIfWebAppOrAuth("Bearer"),
		Validator: func(username, password string, c echo.Context) (bool, error) {
			if authService.IsSuperUser(username, password) {
				c.Set("authCtx", stubAuthContext{})

				return true, nil
			}

			return false, nil
		},
	}))
	g.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: skipIfWebAppOrAuth("Basic"),
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
	g.Use(session.LoadAndSaveWithConfig(session.SessionConfig{
		Skipper:        skipIfNotWebApp(),
		SessionManager: sm,
	}))
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipIfNotWebApp()(c) {
				return next(c)
			}

			authCtx := authService.NewSessionAuthContext(c.Request().Context())
			if authCtx != nil {
				c.Set("authCtx", authCtx)
			}

			return next(c)
		}
	})
}

func RegisterAuthService(g *echo.Group, service *auth.Service) {
	g.POST("/login", func(c echo.Context) error {
		var r auth.LogInRequest
		if err := c.Bind(&r); err != nil {
			return err
		}

		err := service.LogInWithSession(c.Request().Context(), r)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			return err
		}

		return c.JSON(http.StatusOK, map[string]any{"success": true})
	}, session.LoadAndSaveWithConfig(session.SessionConfig{
		SessionManager: service.GetSessionManager(),
	}))
	g.POST("/logout", func(c echo.Context) error {
		err := service.DestroySession(c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{"success": true})
	}, session.LoadAndSaveWithConfig(session.SessionConfig{
		SessionManager: service.GetSessionManager(),
	}))

	g.POST("/authorize", func(c echo.Context) error {
		var r auth.LogInRequest
		if err := c.Bind(&r); err != nil {
			return err
		}

		response, err := service.LogInWithAccessToken(c.Request().Context(), r)
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
	g.GET("/rest/resources", func(c echo.Context) error {
		authCtx, err := getAuthContext(c)
		if err != nil {
			return err
		}

		services := []interface {
			Meta(context.Context, admin.AuthContext) admin.ResourceMeta
		}{
			service.AppService,
			service.AppDemandProfileService,
			service.AuctionConfigurationService,
			service.CountryService,
			service.DemandSourceService,
			service.DemandSourceAccountService,
			service.LineItemService,
			service.SegmentService,
			service.UserService,
		}

		response := make(map[string]admin.ResourceMeta, len(services))

		for _, s := range services {
			meta := s.Meta(c.Request().Context(), authCtx)
			if !meta.Permissions.Read {
				continue
			}

			response[meta.Key] = meta
		}

		return c.JSON(http.StatusOK, response)
	})

	resourceRoutes := []resourceRoute{
		{
			group:   g.Group("/apps"),
			handler: &appServiceHandler{service.AppService},
		},
		{
			group:   g.Group("/app_demand_profiles"),
			handler: &appDemandProfileServiceHandler{service.AppDemandProfileService},
		},
		{
			group:   g.Group("/auction_configurations"),
			handler: &auctionConfigurationServiceHandler{service.AuctionConfigurationService},
		},
		{
			group:   g.Group("/countries"),
			handler: &countryServiceHandler{service.CountryService},
		},
		{
			group:   g.Group("/demand_sources"),
			handler: &demandSourceServiceHandler{service.DemandSourceService},
		},
		{
			group:   g.Group("/demand_source_accounts"),
			handler: &demandSourceAccountServiceHandler{service.DemandSourceAccountService},
		},
		{
			group:   g.Group("/line_items"),
			handler: &lineItemServiceHandler{service.LineItemService},
		},
		{
			group:   g.Group("/segments"),
			handler: &segmentServiceHandler{service.SegmentService},
		},
		{
			group: g.Group("/users"),
			handler: &userHandler{
				userServiceHandler: &userServiceHandler{service.UserService},
			},
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

	lineItemImportHandler := &lineItemImportHandler{service.LineItemService}
	g.POST("/line_items/import", lineItemImportHandler.handleImport)
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

type appServiceHandler = resourceServiceHandler[admin.AppResource, admin.App, admin.AppAttrs]
type appDemandProfileServiceHandler = resourceServiceHandler[admin.AppDemandProfileResource, admin.AppDemandProfile, admin.AppDemandProfileAttrs]
type auctionConfigurationServiceHandler = resourceServiceHandler[admin.AuctionConfigurationResource, admin.AuctionConfiguration, admin.AuctionConfigurationAttrs]
type countryServiceHandler = resourceServiceHandler[admin.CountryResource, admin.Country, admin.CountryAttrs]
type demandSourceServiceHandler = resourceServiceHandler[admin.DemandSourceResource, admin.DemandSource, admin.DemandSourceAttrs]
type demandSourceAccountServiceHandler = resourceServiceHandler[admin.DemandSourceAccountResource, admin.DemandSourceAccount, admin.DemandSourceAccountAttrs]
type lineItemServiceHandler = resourceServiceHandler[admin.LineItemResource, admin.LineItem, admin.LineItemAttrs]
type segmentServiceHandler = resourceServiceHandler[admin.SegmentResource, admin.Segment, admin.SegmentAttrs]
type userServiceHandler = resourceServiceHandler[admin.UserResource, admin.User, admin.UserAttrs]

type resourceServiceHandler[Resource, ResourceData, ResourceAttrs any] struct {
	service resourceService[Resource, ResourceData, ResourceAttrs]
}

type resourceService[Resource, ResourceData, ResourceAttrs any] interface {
	List(ctx context.Context, authCtx admin.AuthContext) ([]Resource, error)
	Find(ctx context.Context, authCtx admin.AuthContext, id int64) (*Resource, error)
	Create(ctx context.Context, authCtx admin.AuthContext, attrs *ResourceAttrs) (*ResourceData, error)
	Update(ctx context.Context, authCtx admin.AuthContext, id int64, attrs *ResourceAttrs) (*ResourceData, error)
	Delete(ctx context.Context, authCtx admin.AuthContext, id int64) error
}

// stubAuthContext is a stub implementation of admin.AuthContext. Build auth context from JWT token.
type stubAuthContext struct{}

func (s stubAuthContext) UserID() int64 {
	return 0
}

func (s stubAuthContext) IsAdmin() bool {
	return true
}

func (s *resourceServiceHandler[Resource, ResourceData, ResourceAttrs]) list(c echo.Context) error {
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

func (s *resourceServiceHandler[Resource, ResourceData, ResourceAttrs]) create(c echo.Context) error {
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	attrs := new(ResourceAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	resource, err := s.service.Create(c.Request().Context(), authCtx, attrs)
	if err != nil {
		var validationError v8n.Errors
		if errors.As(err, &validationError) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, validationError.Error())
		}

		return err
	}

	return c.JSON(http.StatusCreated, resource)
}

func (s *resourceServiceHandler[Resource, ResourceData, ResourceAttrs]) get(c echo.Context) error {
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

func (s *resourceServiceHandler[Resource, ResourceData, ResourceAttrs]) update(c echo.Context) error {
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	attrs := new(ResourceAttrs)
	if err := c.Bind(attrs); err != nil {
		return err
	}

	resource, err := s.service.Update(c.Request().Context(), authCtx, int64(id), attrs)
	if err != nil {
		var validationError v8n.Errors
		if errors.As(err, &validationError) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, validationError.Error())
		}

		return err
	}

	return c.JSON(http.StatusOK, resource)
}

func (s *resourceServiceHandler[Resource, ResourceData, ResourceAttrs]) delete(c echo.Context) error {
	authCtx, err := getAuthContext(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	if err := s.service.Delete(c.Request().Context(), authCtx, int64(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func getAuthContext(c echo.Context) (admin.AuthContext, error) {
	authCtx, ok := c.Get("authCtx").(admin.AuthContext)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized").SetInternal(
			fmt.Errorf("failed to get auth context from request"),
		)
	}

	return authCtx, nil
}

func skipIfWebAppOrAuth(prefixes ...string) middleware.Skipper {
	webAppSkipper := skipIfWebApp()
	authSkipper := skipIfAuthIs(prefixes...)
	return func(c echo.Context) bool {
		return webAppSkipper(c) || authSkipper(c)
	}
}

func skipIfWebApp() middleware.Skipper {
	return func(c echo.Context) bool {
		if c.Request().Header.Get("X-Bidon-App") == "web" {
			return true
		}

		return false
	}
}

func skipIfNotWebApp() middleware.Skipper {
	return func(c echo.Context) bool {
		if c.Request().Header.Get("X-Bidon-App") != "web" {
			return true
		}

		return false
	}
}

func skipIfAuthIs(prefixes ...string) middleware.Skipper {
	return func(c echo.Context) bool {
		header := c.Request().Header.Get(echo.HeaderAuthorization)

		for _, prefix := range prefixes {
			prefix = prefix + " "
			if len(header) >= len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
				return true
			}
		}

		return false
	}
}

type userHandler struct {
	*userServiceHandler
}

func (h *userHandler) get(c echo.Context) error {
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	var id int64
	idParam := c.Param("id")
	if idParam == "me" {
		id = authCtx.UserID()
	} else {
		convID, err := strconv.Atoi(idParam)
		if err != nil {
			return fmt.Errorf("invalid id: %v", err)
		}
		id = int64(convID)
	}

	resource, err := h.service.Find(c.Request().Context(), authCtx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resource)
}

type lineItemImportHandler struct {
	service *admin.LineItemService
}

func (h *lineItemImportHandler) handleImport(c echo.Context) error {
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	attrs := admin.LineItemImportCSVAttrs{}
	if err := c.Bind(&attrs); err != nil {
		return err
	}

	fileHeader, err := c.FormFile("csv")
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("open csv file: %v", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			c.Logger().Errorf("close csv file: %v", err)
		}
	}()

	err = h.service.ImportCSV(c.Request().Context(), authCtx, file, attrs)
	if err != nil {
		return fmt.Errorf("import csv: %v", err)
	}

	return c.NoContent(http.StatusNoContent)
}
