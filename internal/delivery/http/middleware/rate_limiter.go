package middleware

import (
	"context"
	"e-wallet/internal/delivery/http/handler"
	"e-wallet/internal/domain/repository"
	"e-wallet/internal/infrastructure/logger"
	apperrors "e-wallet/pkg/errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	cache  repository.CacheRepository
	limit  int
	window time.Duration
}

func NewRateLimiter(cache repository.CacheRepository, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		cache:  cache,
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		allowed, err := rl.allow(c.Request.Context(), clientIP)
		if err != nil {
			logger.Error.Printf("[RateLimiter]: Error checking rate limit for IP %s: %v", clientIP, err)
			c.Next()
			return
		}

		if !allowed {
			logger.Warning.Printf("[RateLimiter]: Rate limit exceeded for IP: %s", clientIP)
			handler.HandleError(c, apperrors.ErrRateLimitExceeded)
			c.Abort()
			return
		}

		c.Next()
	}
}

// allow checks if request is allowed using Redis atomic counter
func (rl *RateLimiter) allow(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", ip)

	count, err := rl.cache.Incr(ctx, key)
	if err != nil {
		logger.Error.Printf("[RateLimiter]: Error incrementing rate limit for IP %s: %v", ip, err)
		return true, nil
	}

	if count == 1 {
		if err := rl.cache.Expire(ctx, key, rl.window); err != nil {
			logger.Error.Printf("[RateLimiter]: Failed to set expiration for IP %s: %v", ip, err)
		}
	}

	if count > int64(rl.limit) {
		logger.Warning.Printf("[RateLimiter]: Rate limit exceeded for IP %s: %d/%d", ip, count, rl.limit)
		return false, nil
	}

	logger.Debug.Printf("[RateLimiter]: IP %s - %d/%d requests in window", ip, count, rl.limit)

	return true, nil
}
