package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisPinger struct {
	client redis.Cmdable
}

func NewRedisPinger(rdb redis.Cmdable) Pinger {
	return &RedisPinger{
		client: rdb,
	}
}

func (r *RedisPinger) Ping(ctx context.Context) error {
	if r.client == nil {
		return nil
	}

	cmd := r.client.Ping(ctx)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}
