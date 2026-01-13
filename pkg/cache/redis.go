package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

// Client wraps redis.Client with additional methods
type Client struct {
	*redis.Client
}

// NewClient creates a new Redis client
func NewClient(cfg config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{Client: client}, nil
}

// Close closes the Redis client
func (c *Client) Close() error {
	return c.Client.Close()
}

// HealthCheck performs a health check on Redis
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.Ping(ctx).Err()
}

// SetJSON sets a value as JSON with expiration
func (c *Client) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return c.Set(ctx, key, data, expiration).Err()
}

// GetJSON gets a JSON value and unmarshals it
func (c *Client) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := c.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// SetNX sets a value only if it doesn't exist (for distributed locks)
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.Client.SetNX(ctx, key, value, expiration).Result()
}

// Lock acquires a distributed lock
func (c *Client) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return c.SetNX(ctx, key, "locked", ttl)
}

// Unlock releases a distributed lock
func (c *Client) Unlock(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}

// IncrWithExpiry increments a key and sets expiry if it's a new key
func (c *Client) IncrWithExpiry(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	pipe := c.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}
