package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/bidon-io/bidon-backend/cmd/bidon-admin/web"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := db.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	adminService := newAdminService(db)

	e := echo.New()
	e.Use(middleware.Logger())

	configureCORS(e)

	apiGroup := e.Group("/api")
	adminService.RegisterAPIRoutes(apiGroup)

	redocFileSystem, _ := fs.Sub(web.FS, "redoc")
	redocWebServer := http.FileServer(http.FS(redocFileSystem))
	e.GET("/redoc/*", echo.WrapHandler(http.StripPrefix("/redoc/", redocWebServer)))

	uiFileSystem, _ := fs.Sub(web.FS, "ui")
	uiWebServer := http.FileServer(http.FS(uiFileSystem))
	e.GET("/*", echo.WrapHandler(uiWebServer))

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
