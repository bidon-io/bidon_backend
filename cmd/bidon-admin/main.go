package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bidon-io/bidon-backend/cmd/bidon-admin/web"
	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.ConfigureOTel()

	logger, err := config.NewLogger()
	if err != nil {
		log.Fatalf("config.NewLogger(): %v", err)
	}
	defer logger.Sync()

	sentryConf := config.Sentry()
	err = sentry.Init(sentryConf.ClientOptions)
	if err != nil {
		log.Fatalf("sentry.Init(%+v): %v", sentryConf.ClientOptions, err)
	}
	defer sentry.Flush(sentryConf.FlushTimeout)

	dbURL := os.Getenv("DATABASE_URL")
	db, err := db.Open(dbURL)
	if err != nil {
		log.Fatalf("db.Open(%v): %v", dbURL, err)
	}

	e := config.Echo("bidon-admin", logger)

	configureCORS(e)

	adminService := newAdminService(db)

	apiGroup := e.Group("/api")
	adminService.RegisterAPIRoutes(apiGroup)

	redocFileSystem, _ := fs.Sub(web.FS, "redoc")
	redocWebServer := http.FileServer(http.FS(redocFileSystem))
	e.GET("/redoc/*", echo.WrapHandler(http.StripPrefix("/redoc/", redocWebServer)))

	uiFileSystem, _ := fs.Sub(web.FS, "ui")
	uiWebServer := http.FileServer(http.FS(uiFileSystem))
	e.GET("/*", func(c echo.Context) error {
		_, err := uiFileSystem.Open(strings.TrimPrefix(c.Request().URL.Path, "/"))
		if err != nil {
			c.Request().URL.Path = "/"
		}
		echo.WrapHandler(uiWebServer)(c)

		return nil
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	addr := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(addr))
}

func newAdminService(db *db.DB) *admin.Service {
	return &admin.Service{
		AuctionConfigurations: &admin.AuctionConfigurationService{
			Repo: store.NewAuctionConfigurationRepo(db),
		},
		Apps: &admin.AppService{
			Repo: store.NewAppRepo(db),
		},
		AppDemandProfiles: &admin.AppDemandProfileService{
			Repo: store.NewAppDemandProfileRepo(db),
		},
		Segments: &admin.SegmentService{
			Repo: store.NewSegmentRepo(db),
		},
		DemandSourceAccounts: &admin.DemandSourceAccountService{
			Repo: store.NewDemandSourceAccountRepo(db),
		},
		LineItems: &admin.LineItemService{
			Repo: store.NewLineItemRepo(db),
		},
		DemandSources: &admin.DemandSourceService{
			Repo: store.NewDemandSourceRepo(db),
		},
		Countries: &admin.CountryService{
			Repo: store.NewCountryRepo(db),
		},
		Users: &admin.UserService{
			Repo: store.NewUserRepo(db),
		},
	}
}

func configureCORS(e *echo.Echo) {
	if os.Getenv("ENVIRONMENT") == "development" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		}))
	}
}
