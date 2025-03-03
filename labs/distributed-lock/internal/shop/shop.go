package shop

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ericbutera/system-design/labs/distributed-lock/internal/lock"
	"github.com/go-redsync/redsync/v4"
)

const (
	DefaultLockTTL   = 5 * time.Second
	DefaultKeyFormat = "inventory:%s"
)

var (
	ErrReserved = errors.New("item is already reserved")
)

type Shop struct {
	lock      lock.Lock
	KeyFormat string
	LockTTL   time.Duration
}

func New(lock lock.Lock, opts ...Option) *Shop {
	s := &Shop{
		lock:      lock,
		KeyFormat: DefaultKeyFormat,
		LockTTL:   DefaultLockTTL,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Option func(*Shop)

func WithKeyFormat(f string) Option {
	return func(s *Shop) {
		s.KeyFormat = f
	}
}

func WithLockTTL(t time.Duration) Option {
	return func(s *Shop) {
		s.LockTTL = t
	}
}

func (s *Shop) lockName(itemID string) string {
	return fmt.Sprintf(s.KeyFormat, itemID)
}

func (s *Shop) Reserve(ctx context.Context, itemID string, userID string) error {
	key := s.lockName(itemID)
	slog.Debug("Reserving item", "itemID", itemID, "userID", userID)
	if err := s.lock.Set(ctx, key, userID, s.LockTTL); err != nil {
		slog.Debug("Failed to reserve item", "itemID", itemID, "userID", userID, "err", err)
		if errors.Is(err, lock.ErrNotAvailable) {
			return ErrReserved
		}
		return err
	}

	slog.Debug("Reserved item", "itemID", itemID, "userID", userID, "expiry", time.Now().Add(s.LockTTL))
	return nil
}
func (s *Shop) Confirm(ctx context.Context, itemID string, userID string) error {
	key := s.lockName(itemID)

	// Retrieve the current lock value for the item
	current, err := s.lock.Get(ctx, key, userID)
	if err != nil {
		// If the lock is already taken by someone else
		var errTaken *redsync.ErrTaken
		if errors.As(err, &errTaken) {
			return ErrReserved
		}

		// Other errors
		if errors.Is(err, lock.ErrNotAvailable) {
			return ErrReserved
		}

		return err
	}

	// Ensure the current lock value matches the userID
	if current != userID {
		slog.Debug("User mismatch on confirm", "itemID", itemID, "expected", userID, "found", current)
		return ErrReserved
	}

	// Release the lock if userID matches
	if err := s.lock.Delete(ctx, key, userID); err != nil {
		return err
	}

	slog.Debug("Item confirmed and lock released", "itemID", itemID, "userID", userID)
	return nil
}
