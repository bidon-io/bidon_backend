package config

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func Echo() *echo.Echo {
	e := echo.New()
	e.Debug = GetEnv() != ProdEnv
	e.HTTPErrorHandler = HTTPErrorHandler
	e.Validator = &echoValidator{validate: validator.New()}

	e.GET("/health_checks", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	return e
}

func UseCommonMiddleware(g *echo.Group, service string, logger *slog.Logger) {
	g.Use(otelecho.Middleware(service))

	g.Use(middleware.RequestID())
	g.Use(echoRequestLogger(logger))
	g.Use(echoBodyDump(logger))
	g.Use(middleware.Recover())

	g.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))
}

// HTTPErrorHandler is the default error handler for Bidon services.
// Errors that do not wrap echo.HTTPError are considered unexpected and are sent to Sentry.
func HTTPErrorHandler(err error, c echo.Context) {
	// Error handler can be called from middleware, so we need to check if we already handled the error
	if c.Response().Committed {
		return
	}

	var herr *echo.HTTPError
	if !errors.As(err, &herr) {
		hub := sentryecho.GetHubFromContext(c)
		if hub != nil {
			client, scope := hub.Client(), hub.Scope()
			client.CaptureException(
				err,
				&sentry.EventHint{Context: c.Request().Context()},
				scope,
			)
		}

		var message string
		if c.Echo().Debug {
			message = err.Error()
		} else {
			message = http.StatusText(http.StatusInternalServerError)
		}

		herr = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: message,
		}
	}

	response := map[string]any{
		"error": map[string]any{
			"code":    herr.Code,
			"message": herr.Message,
		},
	}

	if err := c.JSON(herr.Code, response); err != nil {
		c.Logger().Error(err)
	}
}

type echoValidator struct {
	validate *validator.Validate
}

func (v *echoValidator) Validate(i any) error {
	if err := v.validate.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return nil
}

func echoBodyDump(logger *slog.Logger) echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		c.Set("reqBody", reqBody)
		c.Set("resBody", resBody)
	})
}

func echoRequestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID: true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogUserAgent: true,
		LogStatus:    true,
		LogError:     true,
		LogLatency:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			reqBody, _ := c.Get("reqBody").([]byte)
			resBody, _ := c.Get("resBody").([]byte)

			logger.Info("request",
				slog.String("id", v.RequestID),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("host", v.Host),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.String("user_agent", v.UserAgent),
				slog.Int("status", v.Status),
				slog.Any("error", v.Error),
				slog.Duration("latency", v.Latency),
				slog.Any("request_body", reqBody),
				slog.Any("response_body", resBody),
			)

			return nil
		},
	})
}
