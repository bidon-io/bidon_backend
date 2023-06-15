package main

import (
	"log"
	"os"

	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionstore "github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	sdkapistore "github.com/bidon-io/bidon-backend/internal/sdkapi/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := db.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	service := sdkapi.Service{
		AuctionBuilder: &auction.Builder{
			ConfigMatcher:    &auctionstore.ConfigMatcher{DB: db},
			LineItemsMatcher: &auctionstore.LineItemsMatcher{DB: db},
		},
		AppFetcher: &sdkapistore.AppFetcher{DB: db},
	}

	e := echo.New()
	e.HTTPErrorHandler = sdkapi.ErrorHandler

	e.Use(middleware.Logger())
	e.Use(sdkapi.CheckBidonHeader)

	e.POST("/auction/:ad_type", service.HandleAuction)
	e.POST("/:ad_type/auction", service.HandleAuction)

	e.Logger.Fatal(e.Start(":1323"))
}
