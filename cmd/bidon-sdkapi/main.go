package main

import (
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/oschwald/maxminddb-golang"
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
	segmentstore "github.com/bidon-io/bidon-backend/internal/segment/store"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
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

	var maxMindDB *maxminddb.Reader

	if os.Getenv("USE_GEOCODING") == "true" {
		maxMindDB, err = maxminddb.Open(os.Getenv("MAXMIND_GEOIP_FILE_PATH"))
		if err != nil {
			log.Fatalf("maxminddb.Open(%v): %v", os.Getenv("MAXMIND_GEOIP_FILE_PATH"), err)
		}
	}

	baseHandler := sdkapi.BaseHandler{
		AppFetcher: &sdkapistore.AppFetcher{DB: db},
		Geocoder:   &geocoder.Geocoder{DB: db, MaxMindDB: maxMindDB},
	}
	segmentMatcher := segment.Matcher{
		Fetcher: &segmentstore.SegmentFetcher{DB: db},
	}
	auctionHandler := sdkapi.AuctionHandler{
		BaseHandler:    &baseHandler,
		SegmentMatcher: &segmentMatcher,
		AuctionBuilder: &auction.Builder{
			ConfigMatcher:    &auctionstore.ConfigMatcher{DB: db},
			LineItemsMatcher: &auctionstore.LineItemsMatcher{DB: db},
		},
	}
	configHandler := sdkapi.ConfigHandler{
		BaseHandler:    &baseHandler,
		SegmentMatcher: &segmentMatcher,
		AdaptersBuilder: &bidonconfig.AdaptersBuilder{
			AppDemandProfileFetcher: &configstore.AppDemandProfileFetcher{DB: db},
		},
	}

	e := config.Echo("bidon-sdkapi", logger)

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
