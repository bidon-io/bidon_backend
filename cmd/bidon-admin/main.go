package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/bidon-io/bidon-backend/cmd/bidon-admin/web"
	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/api"
	"github.com/bidon-io/bidon-backend/internal/admin/auth"
	adminecho "github.com/bidon-io/bidon-backend/internal/admin/echo"
	"github.com/bidon-io/bidon-backend/internal/admin/openapi"
	adminstore "github.com/bidon-io/bidon-backend/internal/admin/store"
	dbpkg "github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bwmarrin/snowflake"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	dbConfig := dbpkg.Config{
		MaxOpenConns:    5 * runtime.GOMAXPROCS(0),
		MaxIdleConns:    1 * runtime.GOMAXPROCS(0),
		ConnMaxLifetime: 15 * time.Minute,
		ReadOnly:        false,
	}
	db, err := dbpkg.Open(dbURL, dbpkg.WithConfig(dbConfig), dbpkg.WithSnowflakeNode(snowflakeNode))
	if err != nil {
		log.Fatalf("dbpkg.Open(%v): %v", dbURL, err)
	}

	e := config.Echo()
	configureCORS(e)

	store := adminstore.New(db)
	authConfig := auth.Config{
		SecretKey:         []byte(os.Getenv("APP_SECRET")),
		SuperUserLogin:    []byte(os.Getenv("SUPERUSER_LOGIN")),
		SuperUserPassword: []byte(os.Getenv("SUPERUSER_PASSWORD")),
	}

	if config.GetEnv() == config.ProdEnv {
		if authConfig.SessionStore, err = gormstore.New(db.DB); err != nil {
			log.Fatal(err)
		}
	}
	authService := auth.NewAuthService(store.UserRepo, store.APIKeyRepo, authConfig)
	adminService := admin.NewService(store)

	g := e.Group("")
	config.UseCommonMiddleware(g, "bidon-admin", logger)
	adminecho.UseAuthorization(g, authService)
	serv := adminecho.NewServer(adminService, authService)
	api.RegisterHandlers(g, serv)

	e.Use(echoprometheus.NewMiddleware("admin"))   // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	oapiWebServer := http.FileServer(http.FS(openapi.FS))
	e.GET("/docs/*", echo.WrapHandler(http.StripPrefix("/docs/", oapiWebServer)))

	uiFileSystem, _ := fs.Sub(web.FS, "ui")
	uiWebServer := echo.WrapHandler(http.FileServer(http.FS(uiFileSystem)))
	e.GET("/", uiWebServer)
	e.GET("/*", uiWebServer, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			file, err := uiFileSystem.Open(strings.TrimPrefix(c.Request().URL.Path, "/"))
			if err != nil {
				c.Request().URL.Path = "/"
				return next(c)
			}
			err = file.Close()
			if err != nil {
				c.Logger().Warnf("Web server file.Close(): %v", err)
			}

			return next(c)
		}
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Errorf("failed to gracefully shutdown http server: %v", err)
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
