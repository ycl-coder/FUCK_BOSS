// Package cache provides cache repository interface for application layer.
// This interface is defined in Application Layer to follow Dependency Inversion Principle.
package cache

import (
	"context"
	"time"
)

// CacheRepository defines the interface for cache operations.
// Implementations are in Infrastructure Layer (e.g., Redis).
type CacheRepository interface {
	// Get retrieves a value from cache by key.
	// Returns the value if found, or an error if not found or operation fails.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value in cache with the specified TTL.
	// If the key already exists, it will be overwritten.
	// Returns an error if the operation fails.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete removes a key from cache.
	// Returns an error if the operation fails.
	Delete(ctx context.Context, key string) error

	// DeleteByPattern removes all keys matching the pattern.
	// Pattern supports wildcards (e.g., "posts:city:*").
	// Returns an error if the operation fails.
	DeleteByPattern(ctx context.Context, pattern string) error
}
