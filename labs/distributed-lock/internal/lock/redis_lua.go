package lock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLua struct {
	rdb *redis.Client
}

func NewRedisLua(rdb *redis.Client) *RedisLua {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{Addr: "redis-master:6379"})
	}
	return &RedisLua{
		rdb: rdb,
	}
}

func (l *RedisLua) Get(ctx context.Context, key string, value string) (string, error) {
	script := redis.NewScript(`
                local val = redis.call("get", KEYS[1])
                if val == false then
                        return redis.error_reply("not found")
                end
                if val ~= ARGV[1] then
                        return redis.error_reply("not available")
                end
                return val
        `)

	val, err := script.Run(ctx, l.rdb, []string{key}, value).Result()
	if err != nil {
		if err.Error() == "not found" {
			return "", ErrNotFound
		} else if err.Error() == "not available" {
			return "", ErrNotAvailable
		}
		return "", err
	}

	return val.(string), nil
}

func (l *RedisLua) Set(ctx context.Context, key string, value string, expireTime time.Duration) error {
	success, err := l.rdb.SetNX(ctx, key, value, expireTime).Result()
	if err != nil {
		return err
	}
	if !success {
		return ErrNotAvailable
	}
	return nil
}

func (l *RedisLua) Delete(ctx context.Context, key string, value string) error {
	script := redis.NewScript(`
                if redis.call("get", KEYS[1]) == ARGV[1] then
                        return redis.call("del", KEYS[1])
                else
                        return 0
                end
        `)

	result, err := script.Run(ctx, l.rdb, []string{key}, value).Result()
	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		return ErrNotAvailable
	}
	return nil
}
