package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo-contrib/echoprometheus"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/bidon-io/bidon-backend/cmd/bidon-admin/web"
	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/auth"
	adminecho "github.com/bidon-io/bidon-backend/internal/admin/echo"
	adminstore "github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bwmarrin/snowflake"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	config.ConfigureOTel()

	logger, err := config.NewLogger()
	if err != nil {
		log.Fatalf("config.NewLogger(): %v", err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Printf("logger.Sync(): %v", err)
		}
	}()

	sentryConf := config.Sentry()
	err = sentry.Init(sentryConf.ClientOptions)
	if err != nil {
		log.Fatalf("sentry.Init(%+v): %v", sentryConf.ClientOptions, err)
	}
	defer sentry.Flush(sentryConf.FlushTimeout)

	dbURL := os.Getenv("DATABASE_URL")
	snowflakeNode, err := prepareSnowflakeNode()
	if err != nil {
		log.Fatalf("prepareSnowflakeNode(): %v", err)
	}
	db, err := db.Open(dbURL, db.WithSnowflakeNode(snowflakeNode))
	if err != nil {
		log.Fatalf("db.Open(%v): %v", dbURL, err)
	}

	e := config.Echo()
	configureCORS(e)

	store := adminstore.New(db)
	authConfig := auth.Config{
		SecretKey:         []byte(os.Getenv("APP_SECRET")),
		SuperUserLogin:    []byte(os.Getenv("SUPERUSER_LOGIN")),
		SuperUserPassword: []byte(os.Getenv("SUPERUSER_PASSWORD")),
	}

	if config.Env == config.ProdEnv {
		redisURL := os.Getenv("REDIS_URL")
		opts, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("redis.ParseURL(%v): %v", redisURL, err)
		}
		rdb := redis.NewClient(opts)
		authConfig.SessionStore = goredisstore.New(rdb)
	}
	authService := auth.NewAuthService(store.UserRepo, authConfig)
	adminService := admin.NewService(store)

	authGroup := e.Group("/auth")
	config.UseCommonMiddleware(authGroup, "bidon-admin", logger)
	adminecho.RegisterAuthService(authGroup, authService)

	apiGroup := e.Group("/api")
	config.UseCommonMiddleware(apiGroup, "bidon-admin", logger)
	adminecho.UseAuthorization(apiGroup, authService)
	adminecho.RegisterAdminService(apiGroup, adminService)

	e.Use(echoprometheus.NewMiddleware("admin"))   // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

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

func configureCORS(e *echo.Echo) {
	if os.Getenv("ENVIRONMENT") == "development" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		}))
	}
}

func prepareSnowflakeNode() (*snowflake.Node, error) {
	snowflakeNodeIDStr := os.Getenv("SNOWFLAKE_NODE_ID")
	if snowflakeNodeIDStr == "" {
		return nil, errors.New("env var SNOWFLAKE_NODE_ID is not set or empty")
	}
	snowflakeNodeID, err := strconv.ParseInt(snowflakeNodeIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SNOWFLAKE_NODE_ID: %v", err)
	}
	node, err := snowflake.NewNode(snowflakeNodeID)
	if err != nil {
		return nil, fmt.Errorf("snowflake.NewNode(%v): %v", snowflakeNodeID, err)
	}

	return node, nil
}
