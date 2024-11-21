package config

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type RedisCache[T any] struct {
	cache *cache.Cache
	ttl   time.Duration
}

// NewRedisCacheOf initializes the RedisCache with a TTL, and local TinyLFU cache.
func NewRedisCacheOf[T any](client *redis.Client, ttl time.Duration) *RedisCache[T] {
	localCache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, ttl),
	})

	return &RedisCache[T]{
		cache: localCache,
		ttl:   ttl,
	}
}

// Get retrieves a value from the cache or loads it using the load function if not found.
func (c *RedisCache[T]) Get(ctx context.Context, key []byte, load func(ctx context.Context) (T, error)) (T, error) {
	var result T

	err := c.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   string(key),
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
