package config

import (
	"time"

	"github.com/bool64/cache"
)

func NewMemoryCacheOf[T any](ttl time.Duration) *cache.FailoverOf[T] {
	return cache.NewFailoverOf[T](func(cfg *cache.FailoverConfigOf[T]) {
		cfg.MaxStaleness = 30 * time.Second   // How long we can serve stale value, update function will be called in background
		cfg.FailedUpdateTTL = 5 * time.Second // Can we serve from cache if update function returned error
		cfg.BackendConfig.TimeToLive = ttl    // TTL
	})
}
