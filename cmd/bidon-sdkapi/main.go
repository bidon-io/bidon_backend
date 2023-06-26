package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionstore "github.com/bidon-io/bidon-backend/internal/auction/store"
	bidonconfig "github.com/bidon-io/bidon-backend/internal/config"
	configstore "github.com/bidon-io/bidon-backend/internal/config/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	sdkapistore "github.com/bidon-io/bidon-backend/internal/sdkapi/store"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	sentryConf := config.Sentry()
	err := sentry.Init(sentryConf.ClientOptions)
	if err != nil {
		log.Fatalf("sentry.Init(%+v): %v", sentryConf.ClientOptions, err)
	}
	defer sentry.Flush(sentryConf.FlushTimeout)

	dbURL := os.Getenv("DATABASE_URL")
	db, err := db.Open(dbURL)
	if err != nil {
		log.Fatalf("db.Open(%v): %v", dbURL, err)
	}

	baseHandler := sdkapi.BaseHandler{
		AppFetcher: &sdkapistore.AppFetcher{DB: db},
	}
	auctionHandler := sdkapi.AuctionHandler{
		BaseHandler: &baseHandler,
		AuctionBuilder: &auction.Builder{
			ConfigMatcher:    &auctionstore.ConfigMatcher{DB: db},
			LineItemsMatcher: &auctionstore.LineItemsMatcher{DB: db},
		},
	}
	configHandler := sdkapi.ConfigHandler{
		BaseHandler: &baseHandler,
		AdaptersBuilder: &bidonconfig.AdaptersBuilder{
			AppDemandProfileFetcher: &configstore.AppDemandProfileFetcher{DB: db},
		},
	}

	e := config.Echo()

	e.Use(sdkapi.CheckBidonHeader)

	e.POST("/config", configHandler.Handle)
	e.POST("/auction/:ad_type", auctionHandler.Handle)
	e.POST("/:ad_type/auction", auctionHandler.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	addr := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(addr))
}
