package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	adapterstore "github.com/bidon-io/bidon-backend/internal/adapter/store"
	"github.com/bidon-io/bidon-backend/internal/auction"
	auctionstore "github.com/bidon-io/bidon-backend/internal/auction/store"
	"github.com/bidon-io/bidon-backend/internal/bidding"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters_builder"
	dbpkg "github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/notification"
	notificationstore "github.com/bidon-io/bidon-backend/internal/notification/store"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/event/engine"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	sdkapistore "github.com/bidon-io/bidon-backend/internal/sdkapi/store"
	v1 "github.com/bidon-io/bidon-backend/internal/sdkapi/v1"
	v2 "github.com/bidon-io/bidon-backend/internal/sdkapi/v2"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/openapi"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentstore "github.com/bidon-io/bidon-backend/internal/segment/store"

	"github.com/bool64/cache"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/oschwald/maxminddb-golang"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

	geoCoder := &geocoder.Geocoder{
		DB:        db,
		MaxMindDB: maxMindDB,
		Cache:     config.NewMemoryCacheOf[*dbpkg.Country](cache.UnlimitedTTL), // We don't update countries
	}
	configFetcher := &auctionstore.ConfigFetcher{
		DB:    db,
		Cache: config.NewMemoryCacheOf[*auction.Config](10 * time.Minute),
	}
	appFetcher := &sdkapistore.AppFetcher{
		DB:    db,
		Cache: config.NewMemoryCacheOf[sdkapi.App](10 * time.Minute),
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: &segmentstore.SegmentFetcher{
			DB:    db,
			Cache: config.NewMemoryCacheOf[[]segment.Segment](10 * time.Minute),
		},
	}
	biddingHttpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: otelhttp.NewTransport(&http.Transport{
			MaxConnsPerHost:     200,
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200, // TODO: Move to config
		}),
	}
	notificationHandler := notification.Handler{
		AuctionResultRepo: notificationstore.AuctionResultRepo{Redis: rdb},
		Sender: notification.EventSender{
			HttpClient:  biddingHttpClient,
			EventLogger: eventLogger,
		},
	}
	notificationHandlerV2 := notification.HandlerV2{
		AuctionResultRepo: notificationstore.AuctionResultV2Repo{Redis: rdb},
		Sender: notification.EventSender{
			HttpClient:  biddingHttpClient,
			EventLogger: eventLogger,
		},
	}
	adUnitsMatcher := &auctionstore.AdUnitsMatcher{
		DB:    db,
		Cache: config.NewMemoryCacheOf[[]auction.AdUnit](10 * time.Minute),
	}
	biddingBuilder := &bidding.Builder{
		AdaptersBuilder:     adapters_builder.BuildBiddingAdapters(biddingHttpClient),
		NotificationHandler: notificationHandler,
	}
	biddingAdaptersCfgBuilder := &adapters_builder.AdaptersConfigBuilder{
		ConfigurationFetcher: &adapterstore.ConfigurationFetcher{
			DB:    db,
			Cache: config.NewMemoryCacheOf[adapter.RawConfigsMap](10 * time.Minute),
		},
	}
	lineItemsMatcher := &auctionstore.LineItemsMatcher{
		DB:    db,
		Cache: config.NewMemoryCacheOf[[]auction.LineItem](10 * time.Minute),
	}
	adapterInitConfigsFetcher := &sdkapistore.AdapterInitConfigsFetcher{DB: db}
	configurationFetcher := &adapterstore.ConfigurationFetcher{
		DB:    db,
		Cache: config.NewMemoryCacheOf[adapter.RawConfigsMap](10 * time.Minute),
	}

	e := config.Echo()

	v1Group := e.Group("")
	config.UseCommonMiddleware(v1Group, "bidon-sdkapi", logger)
	v1Group.Use(sdkapi.CheckBidonHeader)

	routerV1 := v1.Router{
		ConfigFetcher:             configFetcher,
		AppFetcher:                appFetcher,
		SegmentMatcher:            segmentMatcher,
		BiddingBuilder:            biddingBuilder,
		BiddingAdaptersCfgBuilder: biddingAdaptersCfgBuilder,
		AdUnitsMatcher:            adUnitsMatcher,
		NotificationHandler:       notificationHandler,
		GeoCoder:                  geoCoder,
		EventLogger:               eventLogger,
		LineItemsMatcher:          lineItemsMatcher,
		AdapterInitConfigsFetcher: adapterInitConfigsFetcher,
		ConfigurationFetcher:      configurationFetcher,
	}
	routerV1.RegisterRoutes(v1Group)

	v2Group := e.Group("")
	config.UseCommonMiddleware(v2Group, "bidon-sdkapi", logger)
	v2Group.Use(sdkapi.CheckBidonHeader)
	routerV2 := v2.Router{
		ConfigFetcher:             configFetcher,
		AppFetcher:                appFetcher,
		SegmentMatcher:            segmentMatcher,
		BiddingBuilder:            biddingBuilder,
		BiddingAdaptersCfgBuilder: biddingAdaptersCfgBuilder,
		AdUnitsMatcher:            adUnitsMatcher,
		NotificationHandler:       notificationHandlerV2,
		GeoCoder:                  geoCoder,
		EventLogger:               eventLogger,
		LineItemsMatcher:          lineItemsMatcher,
		AdapterInitConfigsFetcher: adapterInitConfigsFetcher,
		ConfigurationFetcher:      configurationFetcher,
	}
	routerV2.RegisterRoutes(v2Group)

	docsWebServer := http.FileServer(http.FS(openapi.FS))
	e.GET("/docs/*", echo.WrapHandler(http.StripPrefix("/docs/", docsWebServer)))

	e.Use(echoprometheus.NewMiddleware("sdkapi"))  // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

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
