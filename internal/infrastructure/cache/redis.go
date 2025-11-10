package cache

import (
	"context"
	"e-wallet/internal/infrastructure/config"
	"e-wallet/internal/infrastructure/logger"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps redis client
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		DisableIdentity: true,
		MinIdleConns:    1,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Suppress stderr during ping to avoid connection error spam
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	pingErr := client.Ping(ctx).Err()

	// Restore stderr
	_ = w.Close()
	os.Stderr = oldStderr

	go func() {
		buf := make([]byte, 1024)
		for {
			if n, _ := r.Read(buf); n == 0 {
				break
			}
		}
		_ = r.Close()
	}()

	if pingErr != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", pingErr)
	}

	logger.Info.Println("Successfully connected to Redis")

	return &RedisClient{client: client}, nil
}

// Get retrieves a value from cache
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set stores a value in cache with expiration
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Delete removes a value from cache
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Incr atomically increments a key and returns the new value
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Expire sets a timeout on a key
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}
