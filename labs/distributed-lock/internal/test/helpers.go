package test

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func FlushDB(ctx context.Context, rdb *redis.Client) error {
	keys, err := rdb.Keys(ctx, "test-*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		slog.Debug("Flushing keys", "keys", keys)
		return rdb.Del(ctx, keys...).Err()
	}
	return nil
}
