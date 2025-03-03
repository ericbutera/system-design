package lock

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotAvailable = errors.New("already reserved")
	ErrNotFound     = errors.New("not found")
)

type Lock interface {
	Get(ctx context.Context, key string, value string) (string, error)
	Set(ctx context.Context, key string, value string, expireTime time.Duration) error
	Delete(ctx context.Context, key string, value string) error
}
