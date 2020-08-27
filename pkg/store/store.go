package store

import (
	"context"
	"time"
)

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expire time.Duration) error
	Delete(ctx context.Context, key string) error
	Keys() ([]string, error)
	Close() error
	Type() string
}