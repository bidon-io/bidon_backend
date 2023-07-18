package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/oschwald/maxminddb-golang"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionstore "github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters_builder"
	biddingstore "github.com/bidon-io/bidon-backend/internal/bidding/store"
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

	var loggerEngine event.LoggerEngine
	if os.Getenv("USE_KAFKA") == "true" {
		conf, err := config.Kafka()
		if err != nil {
			log.Fatalf("config.Kafka(): %v", err)
		}

		client, err := kgo.NewClient(conf.ClientOpts...)
		if err != nil {
			log.Fatalf("kgo.NewClient(): %v", err)
		}
		defer func() {
			ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer ctxCancel()

			err := client.Flush(ctx)
			if err != nil {
				log.Printf("client.Flush(): %v", err)
			}
		}()

		loggerEngine = &engine.Kafka{Client: client, Topics: conf.Topics}
	} else {
		loggerEngine = &engine.Log{}
	}

	appFetcher := &sdkapistore.AppFetcher{DB: db}
	geocoder := &geocoder.Geocoder{DB: db, MaxMindDB: maxMindDB}
	segmentMatcher := segment.Matcher{
		Fetcher: &segmentstore.SegmentFetcher{DB: db},
	}

	biddingHttpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxConnsPerHost:     50,
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50, // TODO: Move to config
		},
	}
	auctionHandler := sdkapi.AuctionHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.AuctionRequest, *schema.AuctionRequest]{
			AppFetcher: appFetcher,
			Geocoder:   geocoder,
		},
		SegmentMatcher: &segmentMatcher,
		AuctionBuilder: &auction.Builder{
			ConfigMatcher:    &auctionstore.ConfigMatcher{DB: db},
			LineItemsMatcher: &auctionstore.LineItemsMatcher{DB: db},
		},
	}
	configHandler := sdkapi.ConfigHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.ConfigRequest, *schema.ConfigRequest]{
			AppFetcher: appFetcher,
			Geocoder:   geocoder,
		},
		SegmentMatcher: &segmentMatcher,
		AdaptersBuilder: &bidonconfig.AdaptersBuilder{
			AppDemandProfileFetcher: &configstore.AppDemandProfileFetcher{DB: db},
		},
		EventLogger: &event.Logger{Engine: loggerEngine},
	}
	biddingHandler := sdkapi.BiddingHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]{
			AppFetcher: appFetcher,
			Geocoder:   geocoder,
		},
		SegmentMatcher: &segmentMatcher,
		BiddingBuilder: &bidding.Builder{
			ConfigMatcher:   &auctionstore.ConfigMatcher{DB: db},
			AdaptersBuilder: adapters_builder.BuildBiddingAdapters(biddingHttpClient),
		},
		AdaptersConfigBuilder: &adapters_builder.AdaptersConfigBuilder{
			AppDemandProfileFetcher: &biddingstore.AppDemandProfileFetcher{DB: db},
		},
	}

	e := config.Echo("bidon-sdkapi", logger)

	e.Use(sdkapi.CheckBidonHeader)

	e.POST("/config", configHandler.Handle)
	e.POST("/auction/:ad_type", auctionHandler.Handle)
	e.POST("/:ad_type/auction", auctionHandler.Handle)
	e.POST("/bidding/:ad_type", biddingHandler.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	addr := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(addr))
}
