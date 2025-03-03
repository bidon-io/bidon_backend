// Package adminecho implements Echo bindings for the admin package.
package adminecho

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	v8n "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	session "github.com/spazzymoto/echo-scs-session"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/auth"
	"github.com/bidon-io/bidon-backend/internal/admin/resource"
)

func UseAuthorization(g *echo.Group, authService *auth.Service) {
	sm := authService.GetSessionManager()
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			skipper := skipIfAny(skipIfWebAppOrAuth("Bearer", "Basic"), skipIfAuthRoutes())
			if skipper(c) {
				return next(c)
			}

			apiKey := c.Request().Header.Get("X-Bidon-Api-Key")
			if apiKey == "" {
				return next(c)
			}

			authCtx, err := authService.ResolveAPIKey(c.Request().Context(), apiKey)
			if err != nil {
				return err
			}

			c.Set("authCtx", authCtx)

			return next(c)
		}
	})
	g.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Skipper: skipIfAny(skipIfWebAppOrAuth("Bearer"), skipIfAuthRoutes(), skipIfApiKey()),
		Validator: func(username, password string, c echo.Context) (bool, error) {
			if authService.IsSuperUser(username, password) {
				c.Set("authCtx", stubAuthContext{})

				return true, nil
			}

			return false, nil
		},
	}))
	g.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: skipIfAny(skipIfWebAppOrAuth("Basic"), skipIfAuthRoutes(), skipIfApiKey()),
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
		Skipper:        skipIfAny(skipIfNotWebApp(), skipIfAuthRoutes()),
		SessionManager: sm,
	}))
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			skipper := skipIfAny(skipIfNotWebApp(), skipIfAuthRoutes())
			if skipper(c) {
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

type appServiceHandler = resourceServiceHandler[admin.AppResource, admin.App, admin.AppAttrs]
type appDemandProfileServiceHandler = resourceServiceHandler[admin.AppDemandProfileResource, admin.AppDemandProfile, admin.AppDemandProfileAttrs]
type auctionConfigurationServiceHandler = resourceServiceHandler[admin.AuctionConfigurationResource, admin.AuctionConfiguration, admin.AuctionConfigurationAttrs]
type auctionConfigurationV2ServiceHandler = resourceServiceHandler[admin.AuctionConfigurationV2Resource, admin.AuctionConfigurationV2, admin.AuctionConfigurationV2Attrs]
type countryServiceHandler = resourceServiceHandler[admin.CountryResource, admin.Country, admin.CountryAttrs]
type demandSourceServiceHandler = resourceServiceHandler[admin.DemandSourceResource, admin.DemandSource, admin.DemandSourceAttrs]
type demandSourceAccountServiceHandler = resourceServiceHandler[admin.DemandSourceAccountResource, admin.DemandSourceAccount, admin.DemandSourceAccountAttrs]
type lineItemServiceHandler = resourceServiceHandler[admin.LineItemResource, admin.LineItem, admin.LineItemAttrs]
type segmentServiceHandler = resourceServiceHandler[admin.SegmentResource, admin.Segment, admin.SegmentAttrs]
type userServiceHandler = resourceServiceHandler[admin.UserResource, admin.User, admin.UserAttrs]
type settingsServiceHandler struct {
	service *admin.SettingsService
}

type resourceServiceHandler[Resource, ResourceData, ResourceAttrs any] struct {
	service resourceService[Resource, ResourceData, ResourceAttrs]
}

type resourceService[Resource, ResourceData, ResourceAttrs any] interface {
	List(ctx context.Context, authCtx admin.AuthContext, qParams map[string][]string) (*resource.Collection[Resource], error)
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

	collection, err := s.service.List(c.Request().Context(), authCtx, c.QueryParams())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collection.Items)
}

func (s *resourceServiceHandler[Resource, ResourceData, ResourceAttrs]) listCollection(c echo.Context) error {
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	collection, err := s.service.List(c.Request().Context(), authCtx, c.QueryParams())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collection)
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
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	if err := s.service.Delete(c.Request().Context(), authCtx, int64(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *settingsServiceHandler) updatePassword(c echo.Context, authCtx admin.AuthContext) error {
	return h.service.UpdatePassword(c, authCtx)
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

func skipIfApiKey() middleware.Skipper {
	return func(c echo.Context) bool {
		return c.Request().Header.Get("X-Bidon-Api-Key") != ""
	}
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
		return c.Request().Header.Get("X-Bidon-App") == "web"
	}
}

func skipIfNotWebApp() middleware.Skipper {
	return func(c echo.Context) bool {
		return c.Request().Header.Get("X-Bidon-App") != "web"
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

func skipIfAuthRoutes() middleware.Skipper {
	return func(c echo.Context) bool {
		return strings.HasPrefix(c.Path(), "/auth")
	}
}

// Combine skippers with OR logic
func skipIfAny(skippers ...middleware.Skipper) middleware.Skipper {
	return func(c echo.Context) bool {
		// Any skipper returning true will make the combined skipper skip
		for _, skipper := range skippers {
			if skipper(c) {
				return true
			}
		}
		return false
	}
}
