// Package cache provides integration tests for Redis cache and rate limiter.
// These tests use a real Redis instance to verify cache and rate limiter implementations.
package cache

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	redisclient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"fuck_boss/backend/internal/infrastructure/persistence/redis"
	apperrors "fuck_boss/backend/pkg/errors"
)

// RedisCacheTestSuite is the test suite for Redis Cache integration tests.
type RedisCacheTestSuite struct {
	suite.Suite
	client      *redisclient.Client
	cacheRepo   *redis.CacheRepository
	rateLimiter *redis.RateLimiter
	ctx         context.Context
}

// SetupSuite runs once before all tests in the suite.
func (s *RedisCacheTestSuite) SetupSuite() {
	// Get Redis connection string from environment or use default
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6380"
	}

	// Create Redis client
	s.client = redisclient.NewClient(&redisclient.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Wait for Redis to be ready
	err := s.waitForRedis()
	require.NoError(s.T(), err, "Redis is not ready")

	// Create cache repository
	s.cacheRepo = redis.NewCacheRepository(s.client)

	// Create rate limiter
	s.rateLimiter = redis.NewRateLimiter(s.client)

	// Create context
	s.ctx = context.Background()
}

// TearDownSuite runs once after all tests in the suite.
func (s *RedisCacheTestSuite) TearDownSuite() {
	if s.client != nil {
		s.client.Close()
	}
}

// SetupTest runs before each test.
func (s *RedisCacheTestSuite) SetupTest() {
	// Clean up any existing test data before each test
	// Use FLUSHDB to clear all keys in the current database
	err := s.client.FlushDB(s.ctx).Err()
	if err != nil {
		s.T().Logf("Failed to flush Redis database: %v", err)
	}
}

// TearDownTest runs after each test.
func (s *RedisCacheTestSuite) TearDownTest() {
	// Clean up test data after each test
	// err := s.client.FlushDB(s.ctx).Err()
	// if err != nil {
	// 	s.T().Logf("Failed to flush Redis database: %v", err)
	// }
}

// waitForRedis waits for Redis to be ready.
func (s *RedisCacheTestSuite) waitForRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		err := s.client.Ping(ctx).Err()
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("redis not ready: %w", ctx.Err())
		case <-time.After(1 * time.Second):
			// Retry
		}
	}
}

// TestCacheRepository_Get tests the Get method.
func (s *RedisCacheTestSuite) TestCacheRepository_Get() {
	// Set a value first
	err := s.cacheRepo.Set(s.ctx, "test:key", "test:value", 5*time.Minute)
	s.Require().NoError(err)

	// Get the value
	value, err := s.cacheRepo.Get(s.ctx, "test:key")
	s.Require().NoError(err)
	s.Equal("test:value", value)
}

// TestCacheRepository_Get_NotFound tests Get with non-existent key.
func (s *RedisCacheTestSuite) TestCacheRepository_Get_NotFound() {
	// Try to get a non-existent key
	value, err := s.cacheRepo.Get(s.ctx, "test:nonexistent")
	s.Require().Error(err)
	s.Empty(value)
	s.True(apperrors.IsNotFoundError(err), "Error should be NotFoundError")
}

// TestCacheRepository_Set tests the Set method.
func (s *RedisCacheTestSuite) TestCacheRepository_Set() {
	// Set a value
	err := s.cacheRepo.Set(s.ctx, "test:key", "test:value", 5*time.Minute)
	s.Require().NoError(err)

	// Verify it was set
	value, err := s.client.Get(s.ctx, "test:key").Result()
	s.Require().NoError(err)
	s.Equal("test:value", value)
}

// TestCacheRepository_Set_TTL tests TTL expiration.
func (s *RedisCacheTestSuite) TestCacheRepository_Set_TTL() {
	// Set a value with short TTL
	err := s.cacheRepo.Set(s.ctx, "test:ttl", "test:value", 2*time.Second)
	s.Require().NoError(err)

	// Verify it exists
	value, err := s.cacheRepo.Get(s.ctx, "test:ttl")
	s.Require().NoError(err)
	s.Equal("test:value", value)

	// Wait for TTL to expire
	time.Sleep(3 * time.Second)

	// Verify it's gone
	_, err = s.cacheRepo.Get(s.ctx, "test:ttl")
	s.Require().Error(err)
	s.True(apperrors.IsNotFoundError(err), "Key should be expired")
}

// TestCacheRepository_Delete tests the Delete method.
func (s *RedisCacheTestSuite) TestCacheRepository_Delete() {
	// Set a value
	err := s.cacheRepo.Set(s.ctx, "test:delete", "test:value", 5*time.Minute)
	s.Require().NoError(err)

	// Delete it
	err = s.cacheRepo.Delete(s.ctx, "test:delete")
	s.Require().NoError(err)

	// Verify it's gone
	_, err = s.cacheRepo.Get(s.ctx, "test:delete")
	s.Require().Error(err)
	s.True(apperrors.IsNotFoundError(err), "Key should be deleted")
}

// TestCacheRepository_DeleteByPattern tests the DeleteByPattern method.
func (s *RedisCacheTestSuite) TestCacheRepository_DeleteByPattern() {
	// Set multiple keys with pattern
	keys := []string{
		"posts:city:beijing:page:1",
		"posts:city:beijing:page:2",
		"posts:city:shanghai:page:1",
		"post:123",
	}

	for _, key := range keys {
		err := s.cacheRepo.Set(s.ctx, key, "value", 5*time.Minute)
		s.Require().NoError(err)
	}

	// Delete by pattern
	err := s.cacheRepo.DeleteByPattern(s.ctx, "posts:city:beijing:*")
	s.Require().NoError(err)

	// Verify beijing keys are deleted
	_, err = s.cacheRepo.Get(s.ctx, "posts:city:beijing:page:1")
	s.Require().Error(err)
	s.True(apperrors.IsNotFoundError(err))

	_, err = s.cacheRepo.Get(s.ctx, "posts:city:beijing:page:2")
	s.Require().Error(err)
	s.True(apperrors.IsNotFoundError(err))

	// Verify other keys still exist
	value, err := s.cacheRepo.Get(s.ctx, "posts:city:shanghai:page:1")
	s.Require().NoError(err)
	s.Equal("value", value)

	value, err = s.cacheRepo.Get(s.ctx, "post:123")
	s.Require().NoError(err)
	s.Equal("value", value)
}

// TestCacheRepository_Validation tests input validation.
func (s *RedisCacheTestSuite) TestCacheRepository_Validation() {
	// Test empty key
	_, err := s.cacheRepo.Get(s.ctx, "")
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))

	err = s.cacheRepo.Set(s.ctx, "", "value", 5*time.Minute)
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))

	err = s.cacheRepo.Delete(s.ctx, "")
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))

	err = s.cacheRepo.DeleteByPattern(s.ctx, "")
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))

	// Test negative TTL
	err = s.cacheRepo.Set(s.ctx, "key", "value", -1*time.Second)
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))
}

// TestRateLimiter_Allow tests the Allow method within limit.
func (s *RedisCacheTestSuite) TestRateLimiter_Allow() {
	key := "rate_limit:test:1"
	limit := 3
	window := 1 * time.Hour

	// First 3 requests should be allowed
	for i := 0; i < limit; i++ {
		allowed, err := s.rateLimiter.Allow(s.ctx, key, limit, window)
		s.Require().NoError(err)
		s.True(allowed, fmt.Sprintf("Request %d should be allowed", i+1))
	}

	// 4th request should be denied
	allowed, err := s.rateLimiter.Allow(s.ctx, key, limit, window)
	s.Require().NoError(err)
	s.False(allowed, "Request 4 should be denied")
}

// TestRateLimiter_Allow_WindowExpiration tests window expiration.
func (s *RedisCacheTestSuite) TestRateLimiter_Allow_WindowExpiration() {
	key := "rate_limit:test:2"
	limit := 2
	window := 2 * time.Second

	// Use up the limit
	for i := 0; i < limit; i++ {
		allowed, err := s.rateLimiter.Allow(s.ctx, key, limit, window)
		s.Require().NoError(err)
		s.True(allowed)
	}

	// Next request should be denied
	allowed, err := s.rateLimiter.Allow(s.ctx, key, limit, window)
	s.Require().NoError(err)
	s.False(allowed)

	// Wait for window to expire
	time.Sleep(3 * time.Second)

	// Request should be allowed again
	allowed, err = s.rateLimiter.Allow(s.ctx, key, limit, window)
	s.Require().NoError(err)
	s.True(allowed, "Request should be allowed after window expiration")
}

// TestRateLimiter_GetRemaining tests the GetRemaining method.
func (s *RedisCacheTestSuite) TestRateLimiter_GetRemaining() {
	key := "rate_limit:test:3"
	limit := 5

	// Initially, all requests are available
	remaining, err := s.rateLimiter.GetRemaining(s.ctx, key, limit)
	s.Require().NoError(err)
	s.Equal(limit, remaining)

	// Use 2 requests
	for i := 0; i < 2; i++ {
		_, err := s.rateLimiter.Allow(s.ctx, key, limit, 1*time.Hour)
		s.Require().NoError(err)
	}

	// Check remaining
	remaining, err = s.rateLimiter.GetRemaining(s.ctx, key, limit)
	s.Require().NoError(err)
	s.Equal(3, remaining)
}

// TestRateLimiter_Reset tests the Reset method.
func (s *RedisCacheTestSuite) TestRateLimiter_Reset() {
	key := "rate_limit:test:4"
	limit := 3

	// Use up the limit
	for i := 0; i < limit; i++ {
		_, err := s.rateLimiter.Allow(s.ctx, key, limit, 1*time.Hour)
		s.Require().NoError(err)
	}

	// Verify limit is exceeded
	allowed, err := s.rateLimiter.Allow(s.ctx, key, limit, 1*time.Hour)
	s.Require().NoError(err)
	s.False(allowed)

	// Reset
	err = s.rateLimiter.Reset(s.ctx, key)
	s.Require().NoError(err)

	// Verify requests are allowed again
	allowed, err = s.rateLimiter.Allow(s.ctx, key, limit, 1*time.Hour)
	s.Require().NoError(err)
	s.True(allowed)
}

// TestRateLimiter_Validation tests input validation.
func (s *RedisCacheTestSuite) TestRateLimiter_Validation() {
	// Test empty key
	_, err := s.rateLimiter.Allow(s.ctx, "", 3, 1*time.Hour)
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))

	// Test zero limit
	_, err = s.rateLimiter.Allow(s.ctx, "key", 0, 1*time.Hour)
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))

	// Test negative limit
	_, err = s.rateLimiter.Allow(s.ctx, "key", -1, 1*time.Hour)
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))

	// Test zero window
	_, err = s.rateLimiter.Allow(s.ctx, "key", 3, 0)
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))

	// Test negative window
	_, err = s.rateLimiter.Allow(s.ctx, "key", 3, -1*time.Second)
	s.Require().Error(err)
	s.True(apperrors.IsValidationError(err))
}

// TestRedisCacheSuite runs all tests in the suite.
func TestRedisCacheSuite(t *testing.T) {
	suite.Run(t, new(RedisCacheTestSuite))
}
