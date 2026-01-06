// Package redis provides Redis implementation of rate limiter.
// It implements the RateLimiter interface defined in Application Layer.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	apperrors "fuck_boss/backend/pkg/errors"
)

// RateLimiter is the Redis implementation of ratelimit.RateLimiter.
// It uses a sliding window algorithm to implement rate limiting.
type RateLimiter struct {
	// client is the Redis client.
	client *redis.Client
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{
		client: client,
	}
}

// Allow checks if a request is allowed within the rate limit using sliding window algorithm.
// Algorithm:
// 1. Use Redis INCR to increment the counter for the key
// 2. If the counter is 1 (first request), set EXPIRE to window duration
// 3. If counter <= limit, allow the request
// 4. If counter > limit, deny the request
//
// This implements a sliding window: each key has its own window that starts
// from the first request and expires after the window duration.
func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	if key == "" {
		return false, apperrors.NewValidationError("rate limit key cannot be empty")
	}

	if limit <= 0 {
		return false, apperrors.NewValidationError("rate limit must be greater than 0")
	}

	if window <= 0 {
		return false, apperrors.NewValidationError("rate limit window must be greater than 0")
	}

	// Use Redis pipeline to execute INCR and EXPIRE atomically
	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return false, apperrors.NewDatabaseErrorWithCause("failed to check rate limit", err)
	}

	// Get the incremented value
	count := incr.Val()

	// If count is 1, it means this is the first request in the window
	// The EXPIRE command will set the expiration
	// If count > limit, rate limit exceeded
	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}

// GetRemaining returns the remaining number of requests allowed for the key.
// Returns -1 if the key doesn't exist or an error occurs.
func (r *RateLimiter) GetRemaining(ctx context.Context, key string, limit int) (int, error) {
	if key == "" {
		return -1, apperrors.NewValidationError("rate limit key cannot be empty")
	}

	if limit <= 0 {
		return -1, apperrors.NewValidationError("rate limit must be greater than 0")
	}

	val, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			// Key doesn't exist, all requests are available
			return limit, nil
		}
		return -1, apperrors.NewDatabaseErrorWithCause("failed to get rate limit count", err)
	}

	remaining := limit - int(val)
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// Reset resets the rate limit counter for the key.
func (r *RateLimiter) Reset(ctx context.Context, key string) error {
	if key == "" {
		return apperrors.NewValidationError("rate limit key cannot be empty")
	}

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return apperrors.NewDatabaseErrorWithCause("failed to reset rate limit", err)
	}

	return nil
}
