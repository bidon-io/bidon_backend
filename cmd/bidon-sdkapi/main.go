package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bool64/cache"
	"github.com/labstack/echo-contrib/echoprometheus"

	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/notification"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/oschwald/maxminddb-golang"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/bidon-io/bidon-backend/config"
	adapterstore "github.com/bidon-io/bidon-backend/internal/adapter/store"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionstore "github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters_builder"
	dbpkg "github.com/bidon-io/bidon-backend/internal/db"
	notificationstore "github.com/bidon-io/bidon-backend/internal/notification/store"
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
	db, err := dbpkg.Open(dbURL)
	if err != nil {
		log.Fatalf("db.Open(%v): %v", dbURL, err)
	}

	redisURL := os.Getenv("REDIS_URL")
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("REDIS_URL parsing failed, using default options: %v", err)
		opts = &redis.Options{}
	}
	rdb := redis.NewClient(opts)

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
	eventLogger := &event.Logger{Engine: loggerEngine}

	configFetcher := &auctionstore.ConfigFetcher{
		DB:    db,
		Cache: config.NewMemoryCacheOf[*auction.Config](10 * time.Minute),
	}
	appFetcher := &sdkapistore.AppFetcher{
		DB:    db,
		Cache: config.NewMemoryCacheOf[sdkapi.App](10 * time.Minute),
	}
	geoCoder := &geocoder.Geocoder{
		DB:        db,
		MaxMindDB: maxMindDB,
		Cache:     config.NewMemoryCacheOf[*dbpkg.Country](cache.UnlimitedTTL), // We don't update countries
	}
	segmentMatcher := segment.Matcher{
		Fetcher: &segmentstore.SegmentFetcher{
			DB:    db,
			Cache: config.NewMemoryCacheOf[[]segment.Segment](10 * time.Minute),
		},
	}

	biddingHttpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxConnsPerHost:     200,
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 50, // TODO: Move to config
		},
	}
	notificationHandler := notification.Handler{
		AuctionResultRepo: notificationstore.AuctionResultRepo{Redis: rdb},
		Sender: notification.EventSender{
			HttpClient:  biddingHttpClient,
			EventLogger: eventLogger,
		},
	}

	adUnitsMatcher := &auctionstore.AdUnitsMatcher{
		DB:    db,
		Cache: config.NewMemoryCacheOf[[]auction.AdUnit](10 * time.Minute),
	}
	auctionHandler := sdkapi.AuctionHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.AuctionRequest, *schema.AuctionRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geoCoder,
		},
		SegmentMatcher: &segmentMatcher,
		AuctionBuilder: &auction.Builder{
			ConfigFetcher: configFetcher,
			LineItemsMatcher: &auctionstore.LineItemsMatcher{
				DB:    db,
				Cache: config.NewMemoryCacheOf[[]auction.LineItem](10 * time.Minute),
			},
		},
		AuctionBuilderV2: &auction.BuilderV2{
			ConfigFetcher:  configFetcher,
			AdUnitsMatcher: adUnitsMatcher,
		},
		EventLogger: eventLogger,
	}
	configHandler := sdkapi.ConfigHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.ConfigRequest, *schema.ConfigRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geoCoder,
		},
		SegmentMatcher:            &segmentMatcher,
		AdapterInitConfigsFetcher: &sdkapistore.AdapterInitConfigsFetcher{DB: db},
		EventLogger:               eventLogger,
	}
	biddingHandler := sdkapi.BiddingHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.BiddingRequest, *schema.BiddingRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geoCoder,
		},
		BiddingBuilder: &bidding.Builder{
			AdaptersBuilder:     adapters_builder.BuildBiddingAdapters(biddingHttpClient),
			NotificationHandler: notificationHandler,
		},
		AdaptersConfigBuilder: &adapters_builder.AdaptersConfigBuilder{
			ConfigurationFetcher: &adapterstore.ConfigurationFetcher{
				DB:    db,
				Cache: config.NewMemoryCacheOf[adapter.RawConfigsMap](10 * time.Minute),
			},
		},
		AdUnitsMatcher: adUnitsMatcher,
		EventLogger:    eventLogger,
	}
	statsHandler := sdkapi.StatsHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.StatsRequest, *schema.StatsRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geoCoder,
		},
		EventLogger:         eventLogger,
		NotificationHandler: notificationHandler,
	}
	showHandler := sdkapi.ShowHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.ShowRequest, *schema.ShowRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geoCoder,
		},
		EventLogger:         eventLogger,
		NotificationHandler: notificationHandler,
	}
	clickHandler := sdkapi.ClickHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.ClickRequest, *schema.ClickRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geoCoder,
		},
		EventLogger: eventLogger,
	}
	rewardHandler := sdkapi.RewardHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.RewardRequest, *schema.RewardRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geoCoder,
		},
		EventLogger: eventLogger,
	}
	lossHandler := sdkapi.LossHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.LossRequest, *schema.LossRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geoCoder,
		},
		EventLogger:         eventLogger,
		NotificationHandler: notificationHandler,
	}
	winHandler := sdkapi.WinHandler{
		BaseHandler: &sdkapi.BaseHandler[schema.WinRequest, *schema.WinRequest]{
			AppFetcher:    appFetcher,
			ConfigFetcher: configFetcher,
			Geocoder:      geoCoder,
		},
		EventLogger:         eventLogger,
		NotificationHandler: notificationHandler,
	}

	e := config.Echo()

	g := e.Group("")
	config.UseCommonMiddleware(g, "bidon-sdkapi", logger)
	g.Use(sdkapi.CheckBidonHeader)

	e.Use(echoprometheus.NewMiddleware("sdkapi"))  // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	g.POST("/config", configHandler.Handle)
	g.POST("/auction/:ad_type", auctionHandler.Handle)
	g.POST("/bidding/:ad_type", biddingHandler.Handle)
	g.POST("/stats/:ad_type", statsHandler.Handle)
	g.POST("/show/:ad_type", showHandler.Handle)
	g.POST("/click/:ad_type", clickHandler.Handle)
	g.POST("/reward/:ad_type", rewardHandler.Handle)
	g.POST("/loss/:ad_type", lossHandler.Handle)
	g.POST("/win/:ad_type", winHandler.Handle)

	// Legacy endpoints
	g.POST("/:ad_type/auction", auctionHandler.Handle)
	g.POST("/:ad_type/stats", statsHandler.Handle)
	g.POST("/:ad_type/show", showHandler.Handle)
	g.POST("/:ad_type/click", clickHandler.Handle)
	g.POST("/:ad_type/reward", rewardHandler.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	addr := fmt.Sprintf(":%s", port)

	go func() {
		err := e.Start(addr)
		if !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatalf("failed to start http server: %v", err)
		}
		e.Logger.Warn(err)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Errorf("failed to gracefully shutdown http server: %v", err)
	}
}
