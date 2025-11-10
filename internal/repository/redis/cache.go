package redis

import (
	"context"
	"e-wallet/internal/infrastructure/cache"
	"e-wallet/internal/infrastructure/logger"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *cache.RedisClient
}

func NewCacheRepository(client *cache.RedisClient) *CacheRepository {
	return &CacheRepository{
		client: client,
	}
}

func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key)
}
func (r *CacheRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration)
}
func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Delete(ctx, key)
}
func (r *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	val, err := r.client.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return val != "", nil
}

func (r *CacheRepository) Incr(ctx context.Context, key string) (int64, error) {
	logger.Debug.Printf("[CacheRepository.Incr]: key=%s", key)
	val, err := r.client.Incr(ctx, key)
	if err != nil {
		logger.Error.Printf("[CacheRepository.Incr]: FAILED key=%s err=%v", key, err)
	} else {
		logger.Debug.Printf("[CacheRepository.Incr]: SUCCESS key=%s value=%d", key, val)
	}
	return val, err
}

func (r *CacheRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	logger.Debug.Printf("[CacheRepository.Expire]: key=%s ttl=%v", key, expiration)
	err := r.client.Expire(ctx, key, expiration)
	if err != nil {
		logger.Error.Printf("[CacheRepository.Expire]: FAILED key=%s err=%v", key, err)
	} else {
		logger.Debug.Printf("[CacheRepository.Expire]: SUCCESS key=%s", key)
	}
	return err
}
