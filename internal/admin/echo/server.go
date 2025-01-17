package adminecho

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/api"
	"github.com/bidon-io/bidon-backend/internal/admin/auth"
	"github.com/labstack/echo/v4"
	session "github.com/spazzymoto/echo-scs-session"
)

type Server struct {
	*admin.Service
	AuthService                *auth.Service
	AppHandler                 *appServiceHandler
	AppDemandProfileHandler    *appDemandProfileServiceHandler
	AucCfgHandler              *auctionConfigurationServiceHandler
	AucCfgV2Handler            *auctionConfigurationV2ServiceHandler
	CountryHandler             *countryServiceHandler
	DemandSourceHandler        *demandSourceServiceHandler
	DemandSourceAccountHandler *demandSourceAccountServiceHandler
	LineItemHandler            *lineItemServiceHandler
	LineItemImportHandler      *lineItemImportHandler
	SegmentHandler             *segmentServiceHandler
	UserHandler                *userHandler
	SettingsHandler            *settingsServiceHandler
}

var _ api.ServerInterface = (*Server)(nil)

func NewServer(service *admin.Service, authService *auth.Service) *Server {
	appHandler := &appServiceHandler{service.AppService}
	appDemandProfileHandler := &appDemandProfileServiceHandler{service.AppDemandProfileService}
	aucHandler := &auctionConfigurationServiceHandler{service.AuctionConfigurationService}
	aucV2Handler := &auctionConfigurationV2ServiceHandler{service.AuctionConfigurationV2Service}
	countryHandler := &countryServiceHandler{service.CountryService}
	demandSourceHandler := &demandSourceServiceHandler{service.DemandSourceService}
	demandSourceAccountHandler := &demandSourceAccountServiceHandler{service.DemandSourceAccountService}
	lineItemHandler := &lineItemServiceHandler{service.LineItemService}
	liImportHandler := &lineItemImportHandler{service.LineItemService}
	segmentHandler := &segmentServiceHandler{service.SegmentService}
	usrHandler := &userHandler{
		userServiceHandler: &userServiceHandler{service.UserService},
	}
	settingsHandler := &settingsServiceHandler{service.SettingsService}

	return &Server{
		Service:                    service,
		AuthService:                authService,
		AppHandler:                 appHandler,
		AppDemandProfileHandler:    appDemandProfileHandler,
		AucCfgHandler:              aucHandler,
		AucCfgV2Handler:            aucV2Handler,
		CountryHandler:             countryHandler,
		DemandSourceHandler:        demandSourceHandler,
		DemandSourceAccountHandler: demandSourceAccountHandler,
		LineItemHandler:            lineItemHandler,
		LineItemImportHandler:      liImportHandler,
		SegmentHandler:             segmentHandler,
		UserHandler:                usrHandler,
		SettingsHandler:            settingsHandler,
	}
}

// App handlers

func (s *Server) GetApps(c echo.Context) error {
	return s.AppHandler.list(c)
}

func (s *Server) CreateApp(c echo.Context) error {
	return s.AppHandler.create(c)
}

func (s *Server) GetApp(c echo.Context, _ api.IdParam) error {
	return s.AppHandler.get(c)
}

func (s *Server) UpdateApp(c echo.Context, _ api.IdParam) error {
	return s.AppHandler.update(c)
}

func (s *Server) DeleteApp(c echo.Context, _ api.IdParam) error {
	return s.AppHandler.delete(c)
}

// AppDemandProfile handlers

func (s *Server) GetAppDemandProfiles(c echo.Context) error {
	return s.AppDemandProfileHandler.list(c)
}

func (s *Server) GetAppDemandProfilesCollection(ctx echo.Context, _ api.GetAppDemandProfilesCollectionParams) error {
	return s.AppDemandProfileHandler.listCollection(ctx)
}

func (s *Server) CreateAppDemandProfile(c echo.Context) error {
	return s.AppDemandProfileHandler.create(c)
}

func (s *Server) GetAppDemandProfile(c echo.Context, _ api.IdParam) error {
	return s.AppDemandProfileHandler.get(c)
}

func (s *Server) UpdateAppDemandProfile(c echo.Context, _ api.IdParam) error {
	return s.AppDemandProfileHandler.update(c)
}

func (s *Server) DeleteAppDemandProfile(c echo.Context, _ api.IdParam) error {
	return s.AppDemandProfileHandler.delete(c)
}

// AuctionConfiguration handlers

func (s *Server) GetAuctionConfigurations(c echo.Context) error {
	return s.AucCfgHandler.list(c)
}

func (s *Server) GetAuctionConfigurationsCollection(ctx echo.Context, _ api.GetAuctionConfigurationsCollectionParams) error {
	return s.AucCfgHandler.listCollection(ctx)
}

func (s *Server) CreateAuctionConfiguration(c echo.Context) error {
	return s.AucCfgHandler.create(c)
}

func (s *Server) GetAuctionConfiguration(c echo.Context, _ api.IdParam) error {
	return s.AucCfgHandler.get(c)
}

func (s *Server) UpdateAuctionConfiguration(c echo.Context, _ api.IdParam) error {
	return s.AucCfgHandler.update(c)
}

func (s *Server) DeleteAuctionConfiguration(c echo.Context, _ api.IdParam) error {
	return s.AucCfgHandler.delete(c)
}

// AuctionConfigurationV2 handlers

func (s *Server) GetAuctionConfigurationsV2(c echo.Context) error {
	return s.AucCfgV2Handler.list(c)
}

func (s *Server) GetAuctionConfigurationsCollectionV2(ctx echo.Context, _ api.GetAuctionConfigurationsCollectionV2Params) error {
	return s.AucCfgV2Handler.listCollection(ctx)
}

func (s *Server) CreateAuctionConfigurationV2(c echo.Context) error {
	return s.AucCfgV2Handler.create(c)
}

func (s *Server) GetAuctionConfigurationV2(c echo.Context, _ api.IdParam) error {
	return s.AucCfgV2Handler.get(c)
}

func (s *Server) UpdateAuctionConfigurationV2(c echo.Context, _ api.IdParam) error {
	return s.AucCfgV2Handler.update(c)
}

func (s *Server) DeleteAuctionConfigurationV2(c echo.Context, _ api.IdParam) error {
	return s.AucCfgV2Handler.delete(c)
}

// Country handlers

func (s *Server) GetCountries(c echo.Context) error {
	return s.CountryHandler.list(c)
}

func (s *Server) CreateCountry(c echo.Context) error {
	return s.CountryHandler.create(c)
}

func (s *Server) GetCountry(c echo.Context, _ api.IdParam) error {
	return s.CountryHandler.get(c)
}

func (s *Server) UpdateCountry(c echo.Context, _ api.IdParam) error {
	return s.CountryHandler.update(c)
}

func (s *Server) DeleteCountry(c echo.Context, _ api.IdParam) error {
	return s.CountryHandler.delete(c)
}

// DemandSource handlers

func (s *Server) GetDemandSources(c echo.Context) error {
	return s.DemandSourceHandler.list(c)
}

func (s *Server) CreateDemandSource(c echo.Context) error {
	return s.DemandSourceHandler.create(c)
}

func (s *Server) GetDemandSource(c echo.Context, _ api.IdParam) error {
	return s.DemandSourceHandler.get(c)
}

func (s *Server) UpdateDemandSource(c echo.Context, _ api.IdParam) error {
	return s.DemandSourceHandler.update(c)
}

func (s *Server) DeleteDemandSource(c echo.Context, _ api.IdParam) error {
	return s.DemandSourceHandler.delete(c)
}

// Demand Source Account handlers

func (s *Server) GetDemandSourceAccounts(c echo.Context) error {
	return s.DemandSourceAccountHandler.list(c)
}

func (s *Server) CreateDemandSourceAccount(c echo.Context) error {
	return s.DemandSourceAccountHandler.create(c)
}

func (s *Server) GetDemandSourceAccount(c echo.Context, _ api.IdParam) error {
	return s.DemandSourceAccountHandler.get(c)
}

func (s *Server) UpdateDemandSourceAccount(c echo.Context, _ api.IdParam) error {
	return s.DemandSourceAccountHandler.update(c)
}

func (s *Server) DeleteDemandSourceAccount(c echo.Context, _ api.IdParam) error {
	return s.DemandSourceAccountHandler.delete(c)
}

// Segment handlers

func (s *Server) GetSegments(c echo.Context) error {
	return s.SegmentHandler.list(c)
}

func (s *Server) CreateSegment(c echo.Context) error {
	return s.SegmentHandler.create(c)
}

func (s *Server) GetSegment(c echo.Context, _ api.IdParam) error {
	return s.SegmentHandler.get(c)
}

func (s *Server) UpdateSegment(c echo.Context, _ api.IdParam) error {
	return s.SegmentHandler.update(c)
}

func (s *Server) DeleteSegment(c echo.Context, _ api.IdParam) error {
	return s.SegmentHandler.delete(c)
}

// User handlers

type userHandler struct {
	*userServiceHandler
}

func (h *userHandler) get(c echo.Context) error {
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	var id int64
	if strings.HasSuffix(c.Path(), "/me") {
		id = authCtx.UserID()
	} else {
		idParam := c.Param("id")
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

func (s *Server) GetUsers(c echo.Context) error {
	return s.UserHandler.list(c)
}

func (s *Server) CreateUser(c echo.Context) error {
	return s.UserHandler.create(c)
}

func (s *Server) GetUser(c echo.Context, _ api.IdParam) error {
	return s.UserHandler.get(c)
}

func (s *Server) GetCurrentUser(c echo.Context) error {
	return s.UserHandler.get(c)
}

func (s *Server) UpdateUser(c echo.Context, _ api.IdParam) error {
	return s.UserHandler.update(c)
}

func (s *Server) DeleteUser(c echo.Context, _ api.IdParam) error {
	return s.UserHandler.delete(c)
}

// LineItem handlers

func (s *Server) GetLineItems(c echo.Context, _ api.GetLineItemsParams) error {
	return s.LineItemHandler.list(c)
}

func (s *Server) GetLineItemsCollection(c echo.Context, _ api.GetLineItemsCollectionParams) error {
	return s.LineItemHandler.listCollection(c)
}

func (s *Server) CreateLineItem(c echo.Context) error {
	return s.LineItemHandler.create(c)
}

func (s *Server) GetLineItem(c echo.Context, _ api.IdParam) error {
	return s.LineItemHandler.get(c)
}

func (s *Server) UpdateLineItem(c echo.Context, _ api.IdParam) error {
	return s.LineItemHandler.update(c)
}

func (s *Server) DeleteLineItem(c echo.Context, _ api.IdParam) error {
	return s.LineItemHandler.delete(c)
}

// Settings handler

func (s *Server) UpdatePassword(c echo.Context) error {
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	return s.SettingsHandler.updatePassword(c, authCtx)
}

func (s *Server) GetResources(c echo.Context) error {
	authCtx, err := getAuthContext(c)
	if err != nil {
		return err
	}

	services := []interface {
		Meta(context.Context, admin.AuthContext) admin.ResourceMeta
	}{
		s.AppService,
		s.AppDemandProfileService,
		s.AuctionConfigurationService,
		s.AuctionConfigurationV2Service,
		s.CountryService,
		s.DemandSourceService,
		s.DemandSourceAccountService,
		s.LineItemService,
		s.SegmentService,
		s.UserService,
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
}

// Import LineItems handlers

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

func (s *Server) ImportLineItems(ctx echo.Context) error {
	return s.LineItemImportHandler.handleImport(ctx)
}

// Auth handlers

func (s *Server) AuthorizeUser(c echo.Context) error {
	var r auth.LogInRequest
	if err := c.Bind(&r); err != nil {
		return err
	}

	response, err := s.AuthService.LogInWithAccessToken(c.Request().Context(), r)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (s *Server) LogIn(ctx echo.Context) error {
	middleware := session.LoadAndSaveWithConfig(session.SessionConfig{
		SessionManager: s.AuthService.GetSessionManager(),
	})

	handler := func(c echo.Context) error {
		var r auth.LogInRequest
		if err := c.Bind(&r); err != nil {
			return err
		}

		err := s.AuthService.LogInWithSession(c.Request().Context(), r)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			return err
		}

		return c.JSON(http.StatusOK, map[string]any{"success": true})
	}

	return middleware(handler)(ctx)
}

func (s *Server) LogOut(ctx echo.Context) error {
	middleware := session.LoadAndSaveWithConfig(session.SessionConfig{
		SessionManager: s.AuthService.GetSessionManager(),
	})

	handler := func(c echo.Context) error {
		err := s.AuthService.DestroySession(c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{"success": true})
	}

	return middleware(handler)(ctx)
}

// Utility handlers

func (s *Server) GetOpenAPISpec(c echo.Context) error {
	spec, err := api.GetSwagger()
	if err != nil {
		return err
	}

	swaggerJSON, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to generate OpenAPI spec")
	}

	return c.JSONBlob(http.StatusOK, swaggerJSON)
}
