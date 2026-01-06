// Package ratelimit provides rate limiter interface for application layer.
// This interface is defined in Application Layer to follow Dependency Inversion Principle.
package ratelimit

import (
	"context"
	"time"
)

// RateLimiter defines the interface for rate limiting operations.
// Implementations are in Infrastructure Layer (e.g., Redis).
type RateLimiter interface {
	// Allow checks if a request is allowed within the rate limit.
	// key: unique identifier for the rate limit (e.g., "rate_limit:post:127.0.0.1:2026-01-06-14")
	// limit: maximum number of requests allowed
	// window: time window for the rate limit (e.g., 1 hour)
	// Returns true if the request is allowed, false if rate limit exceeded, and an error if operation fails.
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}
