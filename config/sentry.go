package config

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

type SentryConfig struct {
	ClientOptions sentry.ClientOptions
	FlushTimeout  time.Duration
}

func Sentry() SentryConfig {
	env := GetEnv()
	return SentryConfig{
		ClientOptions: sentry.ClientOptions{
			Dsn:              os.Getenv("SENTRY_DSN"),
			Debug:            env != ProdEnv,
			AttachStacktrace: true,
			EnableTracing:    true,
			TracesSampleRate: 1.0,
			SendDefaultPII:   true,
			Environment:      env,
		},
		FlushTimeout: 2 * time.Second,
	}
}
