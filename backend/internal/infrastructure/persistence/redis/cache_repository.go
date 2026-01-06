// Package redis provides Redis implementation of cache repository.
// It implements the CacheRepository interface defined in Application Layer.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	apperrors "fuck_boss/backend/pkg/errors"
)

// CacheRepository is the Redis implementation of cache.CacheRepository.
type CacheRepository struct {
	// client is the Redis client.
	client *redis.Client
}

// NewCacheRepository creates a new CacheRepository instance.
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		client: client,
	}
}

// Get retrieves a value from cache by key.
// Returns the value if found, or an error if not found or operation fails.
func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", apperrors.NewValidationError("cache key cannot be empty")
	}

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", apperrors.NewNotFoundError("cache key")
		}
		return "", apperrors.NewDatabaseErrorWithCause("failed to get cache", err)
	}

	return val, nil
}

// Set stores a value in cache with the specified TTL.
// If the key already exists, it will be overwritten.
// Returns an error if the operation fails.
func (r *CacheRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if key == "" {
		return apperrors.NewValidationError("cache key cannot be empty")
	}

	if ttl < 0 {
		return apperrors.NewValidationError("TTL cannot be negative")
	}

	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return apperrors.NewDatabaseErrorWithCause("failed to set cache", err)
	}

	return nil
}

// Delete removes a key from cache.
// Returns an error if the operation fails.
func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	if key == "" {
		return apperrors.NewValidationError("cache key cannot be empty")
	}

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return apperrors.NewDatabaseErrorWithCause("failed to delete cache", err)
	}

	return nil
}

// DeleteByPattern removes all keys matching the pattern.
// Pattern supports wildcards (e.g., "posts:city:*").
// Returns an error if the operation fails.
// Note: This method uses SCAN instead of KEYS to avoid blocking Redis.
func (r *CacheRepository) DeleteByPattern(ctx context.Context, pattern string) error {
	if pattern == "" {
		return apperrors.NewValidationError("pattern cannot be empty")
	}

	// Use SCAN to find all keys matching the pattern
	// This is safer than KEYS command which can block Redis
	// SCAN iterates through keys without blocking the server
	var keys []string
	var cursor uint64 = 0

	for {
		var err error
		var batch []string
		batch, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return apperrors.NewDatabaseErrorWithCause("failed to scan cache keys", err)
		}

		keys = append(keys, batch...)

		// If cursor is 0, we've scanned all keys
		if cursor == 0 {
			break
		}
	}

	// If no keys found, return early
	if len(keys) == 0 {
		return nil
	}

	// Delete all keys in batch
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		return apperrors.NewDatabaseErrorWithCause("failed to delete cache keys by pattern", err)
	}

	return nil
}

// Ping checks if Redis is available.
// This is useful for health checks.
func (r *CacheRepository) Ping(ctx context.Context) error {
	err := r.client.Ping(ctx).Err()
	if err != nil {
		return apperrors.NewDatabaseErrorWithCause("redis ping failed", err)
	}
	return nil
}

// Close closes the Redis client connection.
// This should be called when the application shuts down.
func (r *CacheRepository) Close() error {
	return r.client.Close()
}
