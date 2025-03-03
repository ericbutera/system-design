// Work In Progress
package lock

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/redis/go-redis/v9"
)

type redlock struct {
	rdb  *redis.Client
	pool redsyncredis.Pool
	rs   *redsync.Redsync
}

// func NewRedlock(rdb *redis.Client) *redlock {
// 	if rdb == nil {
// 		rdb = redis.NewClient(&redis.Options{
// 			Addr: "redis-master:6379",
// 		})
// 	}
// 	pool := goredis.NewPool(rdb)
// 	rs := redsync.New(pool)
// 	return &redlock{
// 		rdb:  rdb,
// 		pool: pool,
// 		rs:   rs,
// 	}
// }

func (l *redlock) Get(ctx context.Context, key string, value string) (string, error) {
	mutex := l.rs.NewMutex(key, redsync.WithValue(value))
	val := mutex.Value()
	// if val != value {
	// 	return "", ErrNotAvailable
	// }
	if val == "" {
		return "", ErrNotFound
	}
	return val, nil
}

func (l *redlock) Set(ctx context.Context, key string, value string, expireTime time.Duration) error {
	mutex := l.rs.NewMutex(
		key,
		redsync.WithExpiry(expireTime),
		redsync.WithValue(value),
	)
	err := mutex.Lock()
	if err != nil {
		var errTaken *redsync.ErrTaken
		if errors.As(err, &errTaken) {
			return ErrNotAvailable
		}
		return err
	}
	return nil
}

func (l *redlock) Delete(ctx context.Context, key string, value string) error {
	mutex := l.rs.NewMutex(key, redsync.WithValue(value))
	ok, err := mutex.Unlock()
	if err != nil {
		if strings.Contains(err.Error(), "lock was already expired") {
			return nil
		}
		return err
	}
	if !ok {
		return fmt.Errorf("unable to unlock mutex for key %s", key)
	}
	return nil
}
