package config

import (
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/zap"
)

func Echo() *echo.Echo {
	e := echo.New()
	e.Debug = Debug()
	e.HTTPErrorHandler = HTTPErrorHandler
	e.Validator = &echoValidator{validate: validator.New()}

	e.GET("/health_checks", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	return e
}

type Middleware struct {
	Service               string
	Logger                *zap.Logger
	LogRequestAndResponse bool
}

func UseCommonMiddleware(g *echo.Group, config Middleware) {
	g.Use(otelecho.Middleware(config.Service))

	g.Use(middleware.RequestID())
	g.Use(echoRequestLogger(config.Logger))
	if config.LogRequestAndResponse {
		g.Use(echoBodyDump())
	}
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

func echoBodyDump() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		c.Set("reqBody", reqBody)
		c.Set("resBody", resBody)
	})
}

func echoRequestLogger(logger *zap.Logger) echo.MiddlewareFunc {
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
			fields := []zap.Field{
				zap.String("id", v.RequestID),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.String("user_agent", v.UserAgent),
				zap.Int("status", v.Status),
				zap.NamedError("error", v.Error),
				zap.Duration("latency", v.Latency),
			}

			reqBody, ok := c.Get("reqBody").([]byte)
			if ok {
				fields = append(fields, zap.ByteString("request_body", reqBody))
			}
			resBody, ok := c.Get("resBody").([]byte)
			if ok {
				fields = append(fields, zap.ByteString("response_body", resBody))
			}

			logger.Info("request", fields...)

			return nil
		},
	})
}
