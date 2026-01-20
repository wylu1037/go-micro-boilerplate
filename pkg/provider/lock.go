package provider

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// DistributedLocker defines the interface for distributed locking
type DistributedLocker interface {
	// Lock acquires a distributed lock. ttl is in seconds.
	// Returns an UnlockFunc that must be called to release the lock.
	Lock(ctx context.Context, key string, ttl int) (UnlockFunc, error)
}

type UnlockFunc func(ctx context.Context) error

type etcdLocker struct {
	client *clientv3.Client
}

// NewDistributedLocker creates a new distributed locker.
// If valid client is nil, it returns a no-op locker (useful for simple local dev without etcd).
func NewDistributedLocker(client *clientv3.Client) DistributedLocker {
	if client == nil {
		return &noopLocker{}
	}
	return &etcdLocker{client: client}
}

func (l *etcdLocker) Lock(ctx context.Context, key string, ttl int) (UnlockFunc, error) {
	// Create a session for the lock. The lock will be released if the session expires.
	// concurrency.NewSession starts a background goroutine that keeps the lease alive (Watchdog behavior)
	// as long as the session is active. 'ttl' is the time-to-live if the client disconnects.
	sess, err := concurrency.NewSession(l.client, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd session: %w", err)
	}

	mu := concurrency.NewMutex(sess, "/dlocks/"+key)
	if err := mu.Lock(ctx); err != nil {
		_ = sess.Close()
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return func(ctx context.Context) error {
		defer func() { _ = sess.Close() }()
		return mu.Unlock(ctx)
	}, nil
}

type noopLocker struct{}

func (l *noopLocker) Lock(ctx context.Context, key string, ttl int) (UnlockFunc, error) {
	return func(ctx context.Context) error { return nil }, nil
}
