package shop_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/ericbutera/system-design/labs/distributed-lock/internal/lock"
	"github.com/ericbutera/system-design/labs/distributed-lock/internal/shop"
	"github.com/ericbutera/system-design/labs/distributed-lock/internal/test"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testSetup struct {
	shop *shop.Shop
	lock lock.Lock
	rdb  *redis.Client
	ctx  context.Context
}

func setup(t *testing.T) *testSetup {
	t.Helper()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	require.NoError(t, test.FlushDB(ctx, rdb))
	lock := lock.NewRedisLua(rdb)
	//lock := lock.NewSimple(rdb)
	//lock := lock.NewRedlock(rdb)

	shop := shop.New(lock, shop.WithKeyFormat("test-inventory:%s"))

	return &testSetup{
		rdb:  rdb,
		ctx:  ctx,
		lock: lock,
		shop: shop,
	}
}

func TestReserve(t *testing.T) {
	s := setup(t)

	key := "test-inventory:test-reserve"
	itemID := "test-reserve"
	userID := "user-123"

	err := s.shop.Reserve(s.ctx, itemID, userID)
	require.NoError(t, err)

	//assert.Equal(t, userID, s.rdb.Get(s.ctx, key).Val())
	actual, err := s.lock.Get(s.ctx, key, userID)
	require.NoError(t, err)
	assert.Equal(t, userID, actual)
}

func TestReserveNotAvailable(t *testing.T) {
	s := setup(t)

	itemID := "test-reserve-not-available"
	user1 := "user-1"
	user2 := "user-2"

	err := s.shop.Reserve(s.ctx, itemID, user1)
	require.NoError(t, err)

	err = s.shop.Reserve(s.ctx, itemID, user2)
	assert.ErrorIs(t, err, shop.ErrReserved)
}

func TestConfirm(t *testing.T) {
	s := setup(t)

	itemID := "test-confirm"
	userID := "user-123"

	err := s.shop.Reserve(s.ctx, itemID, userID)
	require.NoError(t, err)

	err = s.shop.Confirm(s.ctx, itemID, userID)
	require.NoError(t, err)

	key := "test-inventory:test-item-123"
	assert.Empty(t, s.rdb.Get(s.ctx, key).Val())
}

func TestConfirmNotAvailable(t *testing.T) {
	s := setup(t)

	itemID := "test-confirm-not-available"
	user1 := "user-1"
	user2 := "user-2"

	err := s.shop.Reserve(s.ctx, itemID, user1)
	require.NoError(t, err)

	err = s.shop.Confirm(s.ctx, itemID, user2)
	assert.ErrorIs(t, err, shop.ErrReserved)
}
