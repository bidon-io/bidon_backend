package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/bool64/cache"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/oschwald/maxminddb-golang"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc/reflection"

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
	grpcserver "github.com/bidon-io/bidon-backend/internal/sdkapi/grpc"
	sdkapistore "github.com/bidon-io/bidon-backend/internal/sdkapi/store"
	v2 "github.com/bidon-io/bidon-backend/internal/sdkapi/v2"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/v2/openapi"
	"github.com/bidon-io/bidon-backend/internal/segment"
	segmentstore "github.com/bidon-io/bidon-backend/internal/segment/store"
	"github.com/bidon-io/bidon-backend/pkg/clock"
	pb "github.com/bidon-io/bidon-backend/pkg/proto/org/bidon/proto/v1"
)

var cpus = runtime.GOMAXPROCS(0)

func main() {
	config.ConfigureOTel()
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("prometheus.New(): %v", err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("bidon-sdkapi")

	logger, err := config.NewLogger()
	if err != nil {
		log.Fatalf("config.NewLogger(): %v", err)
	}
	defer logger.Sync() //nolint:errcheck

	sentryConf := config.Sentry()
	err = sentry.Init(sentryConf.ClientOptions)
	if err != nil {
		log.Fatalf("sentry.Init(%+v): %v", sentryConf.ClientOptions, err)
	}
	defer sentry.Flush(sentryConf.FlushTimeout)

	dbURL := os.Getenv("DATABASE_REPLICA_URL")
	dbConfig := dbpkg.Config{
		MaxOpenConns:    10 * cpus,
		MaxIdleConns:    5 * cpus,
		ConnMaxLifetime: 15 * time.Minute,
		ReadOnly:        true,
	}
	db, err := dbpkg.Open(dbURL, dbpkg.WithConfig(dbConfig))
	if err != nil {
		log.Fatalf("db.Open(%v): %v", dbURL, err)
	}

	redisClusterAddrs := os.Getenv("REDIS_CLUSTER")
	if redisClusterAddrs == "" {
		log.Fatalf("REDIS_CLUSTER is not set")
	}
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    strings.Split(redisClusterAddrs, ","),
		PoolSize: 10 * cpus,
	})

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
				log.Printf("kgo.Client.Flush(): %v", err)
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
	auctionCache := config.NewRedisCacheOf[*auction.Config](rdb, 10*time.Minute, "auction_configs")
	err = auctionCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for auctionCache: %v", err)
	}
	configFetcher := &auctionstore.ConfigFetcher{
		DB:    db,
		Cache: auctionCache,
	}
	appCache := config.NewRedisCacheOf[sdkapi.App](rdb, 10*time.Minute, "apps")
	err = appCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for appCache: %v", err)
	}
	appFetcher := &sdkapistore.AppFetcher{
		DB:    db,
		Cache: appCache,
	}
	segmentCache := config.NewRedisCacheOf[[]segment.Segment](rdb, 10*time.Minute, "segments")
	err = segmentCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for segmentCache: %v", err)
	}
	segmentMatcher := &segment.Matcher{
		Fetcher: &segmentstore.SegmentFetcher{
			DB:    db,
			Cache: segmentCache,
		},
	}
	biddingHTTPClient := &http.Client{
		Timeout: 4 * time.Second,
		Transport: otelhttp.NewTransport(&http.Transport{
			MaxConnsPerHost:     30 * cpus,
			MaxIdleConns:        30 * cpus,
			MaxIdleConnsPerHost: 30 * cpus,
		}),
	}
	notificationHandler := notification.Handler{
		AuctionResultRepo: notificationstore.AuctionResultRepo{Redis: rdb},
		Sender: notification.EventSender{
			HttpClient:  biddingHTTPClient,
			EventLogger: eventLogger,
		},
	}
	adUnitsCache := config.NewRedisCacheOf[[]auction.AdUnit](rdb, 10*time.Minute, "ad_units")
	err = adUnitsCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for adUnitsCache: %v", err)
	}
	adUnitsMatcher := &auctionstore.AdUnitsMatcher{
		DB:    db,
		Cache: adUnitsCache,
	}
	biddingBuilder := &bidding.Builder{
		AdaptersBuilder:     adapters_builder.BuildBiddingAdapters(biddingHTTPClient),
		NotificationHandler: notificationHandler,
		BidCacher:           &bidding.BidCache{Redis: rdb, Clock: clock.New()},
	}
	biddingAdaptersCfgCache := config.NewRedisCacheOf[adapter.RawConfigsMap](rdb, 10*time.Minute, "bidding_adapters_cfg")
	err = biddingAdaptersCfgCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for biddingAdaptersCfgCache: %v", err)
	}
	demandCfg := config.NewDemandConfig()
	biddingAdaptersCfgBuilder := adapters_builder.NewAdaptersConfigBuilder(
		&adapterstore.ConfigurationFetcher{
			DB:    db,
			Cache: biddingAdaptersCfgCache,
		},
		demandCfg,
	)
	lineItemsCache := config.NewRedisCacheOf[[]dbpkg.LineItem](rdb, 10*time.Minute, "line_items")
	err = lineItemsCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for lineItemsCache: %v", err)
	}
	profilesCache := config.NewRedisCacheOf[[]dbpkg.AppDemandProfile](rdb, 10*time.Minute, "app_demand_profiles")
	err = profilesCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for profilesCache: %v", err)
	}
	amazonSlotsCache := config.NewRedisCacheOf[[]sdkapi.AmazonSlot](rdb, 10*time.Minute, "amazon_slots")
	err = amazonSlotsCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for amazonSlotsCache: %v", err)
	}

	adapterInitConfigsFetcher := &sdkapistore.AdapterInitConfigsFetcher{DB: db, ProfilesCache: profilesCache, AmazonSlotsCache: amazonSlotsCache, LineItemsCache: lineItemsCache}
	configsCache := config.NewRedisCacheOf[adapter.RawConfigsMap](rdb, 10*time.Minute, "configs")
	err = configsCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for configsCache: %v", err)
	}
	configurationFetcher := &adapterstore.ConfigurationFetcher{
		DB:    db,
		Cache: configsCache,
	}
	adUnitLookupCache := config.NewRedisCacheOf[*dbpkg.LineItem](rdb, 10*time.Minute, "ad_unit_lookup")
	err = adUnitLookupCache.Monitor(meter)
	if err != nil {
		log.Fatalf("Unable to register observer for adUnitLookupCache: %v", err)
	}
	adUnitLookup := &sdkapistore.AdUnitLookup{
		DB:    db,
		Cache: adUnitLookupCache,
	}
	auctionService := &auction.Service{
		ConfigFetcher:      configFetcher,
		SegmentMatcher:     segmentMatcher,
		AdapterKeysFetcher: adapterInitConfigsFetcher,
		AuctionBuilder: &auction.Builder{
			AdUnitsMatcher:               adUnitsMatcher,
			BiddingBuilder:               biddingBuilder,
			BiddingAdaptersConfigBuilder: biddingAdaptersCfgBuilder,
		},
		EventLogger: eventLogger,
	}

	e := config.Echo()

	v2Group := e.Group("")
	config.UseCommonMiddleware(v2Group, config.Middleware{
		Service:               "bidon-sdkapi",
		Logger:                logger,
		LogRequestAndResponse: true,
	})
	v2Group.Use(sdkapi.CheckBidonHeader)
	routerV2 := v2.Router{
		ConfigFetcher:             configFetcher,
		AppFetcher:                appFetcher,
		SegmentMatcher:            segmentMatcher,
		BiddingBuilder:            biddingBuilder,
		AdUnitsMatcher:            adUnitsMatcher,
		NotificationHandler:       notificationHandler,
		GeoCoder:                  geoCoder,
		EventLogger:               eventLogger,
		AdapterInitConfigsFetcher: adapterInitConfigsFetcher,
		ConfigurationFetcher:      configurationFetcher,
		AuctionService:            auctionService,
		AdUnitLookup:              adUnitLookup,
	}
	routerV2.RegisterRoutes(v2Group)

	docsWebServer := http.FileServer(http.FS(openapi.FS))
	e.GET("/docs/*", echo.WrapHandler(http.StripPrefix("/docs/", docsWebServer)))

	e.Use(echoprometheus.NewMiddleware("sdkapi"))  // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	config.UseHealthCheckHandler(e, config.HealthCheckParams{
		"db":    db,
		"redis": config.NewRedisPinger(rdb),
		"kafka": eventLogger.Engine,
	})

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

	grpcServer := config.NewGRPCServer(logger)
	go func() {
		grpcPort := os.Getenv("GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "50051"
		}
		grpcAddr := fmt.Sprintf(":%s", grpcPort)

		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcAddr, err)
		}

		server := grpcserver.NewServer(auctionService, appFetcher, geoCoder)
		pb.RegisterBiddingServiceServer(grpcServer, server)
		if os.Getenv("ENVIRONMENT") == "development" {
			reflection.Register(grpcServer)
		}

		log.Printf("gRPC server is listening on %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Errorf("failed to gracefully shutdown http server: %v", err)
	}

	grpcServer.GracefulStop()
}
