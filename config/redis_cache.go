package config

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

const CacheNamespace = "cache"

type RedisCache[T any] struct {
	Label string
	cache *cache.Cache
	ttl   time.Duration
}

// NewRedisCacheOf initializes the RedisCache with a TTL, and local TinyLFU cache.
func NewRedisCacheOf[T any](client *redis.ClusterClient, ttl time.Duration, label string) *RedisCache[T] {
	localCache := cache.New(&cache.Options{
		Redis:        client,
		StatsEnabled: true,
		LocalCache:   cache.NewTinyLFU(1000, ttl),
	})

	return &RedisCache[T]{
		Label: label,
		cache: localCache,
		ttl:   ttl,
	}
}

// Get retrieves a value from the cache or loads it using the load function if not found.
func (c *RedisCache[T]) Get(ctx context.Context, key []byte, load func(ctx context.Context) (T, error)) (T, error) {
	var result T

	err := c.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   c.namespacedKey(key),
		Value: &result,
		TTL:   c.ttl,
		Do: func(*cache.Item) (any, error) {
			return load(ctx)
		},
	})

	if err != nil {
		var zero T
		return zero, err
	}

	return result, nil
}

// Monitor sets up OpenTelemetry monitoring for the RedisCache.
func (c *RedisCache[T]) Monitor(meter metric.Meter) error {
	counter, err := meter.Int64ObservableCounter("cache.stats", metric.WithDescription("Cache stats"))
	if err != nil {
		return err
	}

	_, err = meter.RegisterCallback(
		func(ctx context.Context, observer metric.Observer) error {
			stats := c.cache.Stats()
			cacheLabelAttr := attribute.String("label", c.Label)

			observer.ObserveInt64(counter, int64(stats.Hits), metric.WithAttributes(cacheLabelAttr, attribute.String("type", "hit")))
			observer.ObserveInt64(counter, int64(stats.Misses), metric.WithAttributes(cacheLabelAttr, attribute.String("type", "miss")))

			return nil
		},
		counter,
	)

	return err
}

func (c *RedisCache[T]) namespacedKey(key []byte) string {
	return CacheNamespace + ":" + string(key)
}
