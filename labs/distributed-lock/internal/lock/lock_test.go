package lock_test

import (
	"context"
	"testing"
	"time"

	"github.com/ericbutera/system-design/labs/distributed-lock/internal/lock"
	"github.com/ericbutera/system-design/labs/distributed-lock/internal/test"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testSetup struct {
	lock lock.Lock
	rdb  *redis.Client
	ctx  context.Context
}

func setup(t *testing.T) *testSetup {
	t.Helper()

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	require.NoError(t, test.FlushDB(ctx, rdb))

	// lock := lock.NewRedlock(rdb)
	//lock := lock.NewSimple(rdb)
	lock := lock.NewRedisLua(rdb)

	return &testSetup{
		rdb:  rdb,
		ctx:  ctx,
		lock: lock,
	}
}

func TestSet(t *testing.T) {
	s := setup(t)

	key := "test-key"
	value := "test-value"
	expireTime := 5 * time.Second

	err := s.lock.Set(s.ctx, key, value, expireTime)
	assert.NoError(t, err)

	actual, err := s.lock.Get(s.ctx, key, value)
	require.NoError(t, err)
	assert.Equal(t, value, actual)
}

func TestSetWithExisting(t *testing.T) {
	s := setup(t)

	key := "test-key"
	expireTime := 5 * time.Second

	err := s.lock.Set(s.ctx, key, "first", expireTime)
	require.NoError(t, err)

	s.lock.Set(s.ctx, key, "second", expireTime)
	require.NoError(t, err)

	actual, err := s.lock.Get(s.ctx, key, "first")
	require.NoError(t, err)
	assert.Equal(t, "first", actual)
}

func TestExpiredCanLock(t *testing.T) {
	s := setup(t)

	key := "test-key"
	expireTime := 10 * time.Millisecond

	err := s.lock.Set(s.ctx, key, "first", expireTime)
	require.NoError(t, err)

	time.Sleep(20 * time.Millisecond)

	s.lock.Set(s.ctx, key, "second", expireTime)
	require.NoError(t, err)

	actual, err := s.lock.Get(s.ctx, key, "second")
	require.NoError(t, err)
	assert.Equal(t, "second", actual)
}
