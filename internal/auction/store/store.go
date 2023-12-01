package store

import (
	"context"
)

type cache[T any] interface {
	Get(context.Context, []byte, func(ctx context.Context) (T, error)) (T, error)
}
