package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bidon-io/bidon-backend/cmd/bidon-admin/web"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := openDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	handlers := &admin.Handlers{
		AuctionConfigurationRepo: &store.AuctionConfigurationRepo{
			DB: db,
		},
	}

	e := echo.New()
	e.Use(middleware.Logger())

	apiGroup := e.Group("/api")
	handlers.RegisterRoutes(apiGroup)

	webServer := http.FileServer(http.FS(web.FS))
	e.GET("/*", echo.WrapHandler(webServer))

	e.Logger.Fatal(e.Start(":1323"))
}

func openDB(databaseUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl))
	if err != nil {
		return nil, err
	}

	return db, nil
}
