package config

import (
	"errors"
	"net/http"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func Echo(logger *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Debug = Env != ProdEnv
	e.HTTPErrorHandler = HTTPErrorHandler

	e.Use(middleware.RequestID())
	e.Use(echoRequestLogger(logger))
	e.Use(echoBodyDump(logger))
	e.Use(middleware.Recover())

	e.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))

	return e
}

// HTTPErrorHandler is the default error handler for Bidon services.
// Errors that do not wrap echo.HTTPError are considered unexpected and are sent to Sentry.
func HTTPErrorHandler(err error, c echo.Context) {
	var herr *echo.HTTPError
	if !errors.As(err, &herr) {
		hub := sentryecho.GetHubFromContext(c)
		if hub != nil {
			hub.CaptureException(err)
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

func echoBodyDump(logger *zap.Logger) echo.MiddlewareFunc {
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
			reqBody, _ := c.Get("reqBody").([]byte)
			resBody, _ := c.Get("resBody").([]byte)

			logger.Info("request",
				zap.String("id", v.RequestID),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.String("user_agent", v.UserAgent),
				zap.Int("status", v.Status),
				zap.NamedError("error", v.Error),
				zap.Duration("latency", v.Latency),
				zap.ByteString("request_body", reqBody),
				zap.ByteString("response_body", resBody),
			)

			return nil
		},
	})
}
