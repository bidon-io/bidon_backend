package config

import (
	"time"

	"github.com/bool64/cache"
)

func NewMemoryCacheOf[T any](staleness, failedUpdateTTL, ttl time.Duration) *cache.FailoverOf[T] {
	return cache.NewFailoverOf[T](func(cfg *cache.FailoverConfigOf[T]) {
		cfg.MaxStaleness = staleness          // How long we can serve stale value, update function will be called in background
		cfg.FailedUpdateTTL = failedUpdateTTL // Can we serve from cache if update function returned error
		cfg.BackendConfig.TimeToLive = ttl    // TTL
	})
}
