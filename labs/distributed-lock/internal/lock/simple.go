package lock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Simple struct {
	rdb *redis.Client
}

func NewSimple(rdb *redis.Client) *Simple {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{Addr: "redis-master:6379"})
	}
	return &Simple{
		rdb: rdb,
	}
}

func (l *Simple) Get(ctx context.Context, key string, value string) (string, error) {
	val, err := l.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	if val != value {
		return "", ErrNotAvailable
	}
	if val == "" {
		return "", ErrNotFound
	}
	return val, nil
}

func (l *Simple) Set(ctx context.Context, key string, value string, expireTime time.Duration) error {
	success, err := l.rdb.SetNX(ctx, key, value, expireTime).Result()
	if err != nil {
		return err
	}
	if !success {
		return ErrNotAvailable
	}
	return nil
}

func (l *Simple) Delete(ctx context.Context, key string, _ string) error {
	_, err := l.rdb.Del(ctx, key).Result()
	return err
}
