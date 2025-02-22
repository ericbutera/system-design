package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

const Key = "key"

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis-master:6379",
	})

	ctx := context.Background()

	val := time.Now().String()

	lo.Must0(rdb.Set(ctx, Key, val, 0).Err())
	slog.Info("writing", "key", Key, "value", val)

	out := rdb.Get(ctx, Key).Val()
	slog.Info("read", "value", out)

	time.Sleep(15 * time.Minute)
}
