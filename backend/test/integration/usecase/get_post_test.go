// Package usecase provides integration tests for use cases.
// These tests use real PostgreSQL and Redis to verify complete use case flows.
package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	redisclient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"fuck_boss/backend/internal/application/content"
	domaincontent "fuck_boss/backend/internal/domain/content"
	"fuck_boss/backend/internal/domain/shared"
	"fuck_boss/backend/internal/infrastructure/persistence/postgres"
	"fuck_boss/backend/internal/infrastructure/persistence/redis"
	apperrors "fuck_boss/backend/pkg/errors"
)

// GetPostUseCaseTestSuite is the test suite for GetPostUseCase integration tests.
type GetPostUseCaseTestSuite struct {
	suite.Suite
	db          *sql.DB
	redisClient *redisclient.Client
	useCase     *content.GetPostUseCase
	ctx         context.Context
}

// SetupSuite runs once before all tests in the suite.
func (s *GetPostUseCaseTestSuite) SetupSuite() {
	// Setup PostgreSQL
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://test_user:test_password@localhost:5433/test_db?sslmode=disable"
	}

	var err error
	s.db, err = sql.Open("postgres", dsn)
	require.NoError(s.T(), err, "Failed to connect to test database")

	err = s.waitForDB()
	require.NoError(s.T(), err, "Database is not ready")

	err = s.runMigrations()
	require.NoError(s.T(), err, "Failed to run migrations")

	// Setup Redis
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6380"
	}

	s.redisClient = redisclient.NewClient(&redisclient.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	err = s.waitForRedis()
	require.NoError(s.T(), err, "Redis is not ready")

	// Create repositories
	postRepo := postgres.NewPostRepository(s.db)
	cacheRepo := redis.NewCacheRepository(s.redisClient)

	// Create use case
	s.useCase = content.NewGetPostUseCase(postRepo, cacheRepo)

	// Create context
	s.ctx = context.Background()
}

// TearDownSuite runs once after all tests in the suite.
func (s *GetPostUseCaseTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
}

// SetupTest runs before each test.
func (s *GetPostUseCaseTestSuite) SetupTest() {
	// Clean up database
	_, err := s.db.ExecContext(s.ctx, "TRUNCATE TABLE posts CASCADE")
	if err != nil {
		s.T().Logf("Failed to truncate posts table: %v", err)
	}

	// Clean up Redis
	err = s.redisClient.FlushDB(s.ctx).Err()
	if err != nil {
		s.T().Logf("Failed to flush Redis database: %v", err)
	}
}

// TearDownTest runs after each test.
func (s *GetPostUseCaseTestSuite) TearDownTest() {
	// Clean up database
	_, err := s.db.ExecContext(s.ctx, "TRUNCATE TABLE posts CASCADE")
	if err != nil {
		s.T().Logf("Failed to truncate posts table: %v", err)
	}

	// Clean up Redis
	err = s.redisClient.FlushDB(s.ctx).Err()
	if err != nil {
		s.T().Logf("Failed to flush Redis database: %v", err)
	}
}

// waitForDB waits for PostgreSQL to be ready.
func (s *GetPostUseCaseTestSuite) waitForDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		err := s.db.PingContext(ctx)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("database not ready: %w", ctx.Err())
		case <-time.After(1 * time.Second):
			// Retry
		}
	}
}

// waitForRedis waits for Redis to be ready.
func (s *GetPostUseCaseTestSuite) waitForRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		err := s.redisClient.Ping(ctx).Err()
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

// runMigrations runs database migrations.
func (s *GetPostUseCaseTestSuite) runMigrations() error {
	ctx := context.Background()

	// Create cities table
	citiesSQL := `
	CREATE TABLE IF NOT EXISTS cities (
		code VARCHAR(50) PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		pinyin VARCHAR(100),
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_cities_name ON cities(name);
	`

	_, err := s.db.ExecContext(ctx, citiesSQL)
	if err != nil {
		return fmt.Errorf("failed to create cities table: %w", err)
	}

	// Insert test cities
	insertCitiesSQL := `
	INSERT INTO cities (code, name, pinyin) VALUES
		('beijing', '北京', 'beijing'),
		('shanghai', '上海', 'shanghai'),
		('guangzhou', '广州', 'guangzhou'),
		('shenzhen', '深圳', 'shenzhen'),
		('hangzhou', '杭州', 'hangzhou')
	ON CONFLICT (code) DO NOTHING;
	`

	_, err = s.db.ExecContext(ctx, insertCitiesSQL)
	if err != nil {
		return fmt.Errorf("failed to insert test cities: %w", err)
	}

	// Create posts table
	postsSQL := `
	CREATE TABLE IF NOT EXISTS posts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		company_name VARCHAR(100) NOT NULL,
		city_code VARCHAR(50) NOT NULL,
		city_name VARCHAR(100) NOT NULL,
		content TEXT NOT NULL,
		occurred_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		FOREIGN KEY (city_code) REFERENCES cities(code)
	);

	CREATE INDEX IF NOT EXISTS idx_posts_city_code ON posts(city_code);
	CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_posts_company_name ON posts(company_name);
	CREATE INDEX IF NOT EXISTS idx_posts_fulltext ON posts USING GIN(to_tsvector('simple', company_name || ' ' || content));
	`

	_, err = s.db.ExecContext(ctx, postsSQL)
	if err != nil {
		return fmt.Errorf("failed to create posts table: %w", err)
	}

	return nil
}

// TestGetPostUseCase_Execute_Success tests successful post retrieval.
func (s *GetPostUseCaseTestSuite) TestGetPostUseCase_Execute_Success() {
	// Create test post
	repo := postgres.NewPostRepository(s.db)
	company, _ := domaincontent.NewCompanyName("详情测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条详情测试内容，用于验证获取功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	err := repo.Save(s.ctx, post)
	require.NoError(s.T(), err)

	postID := post.ID().String()

	// Query
	result, err := s.useCase.Execute(s.ctx, postID)

	// Assertions
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	assert.Equal(s.T(), postID, result.ID)
	assert.Equal(s.T(), "详情测试公司", result.Company)
	assert.Equal(s.T(), "beijing", result.CityCode)
	assert.Equal(s.T(), "北京", result.CityName)
	assert.Equal(s.T(), postContent.String(), result.Content)
	assert.NotZero(s.T(), result.CreatedAt)
}

// TestGetPostUseCase_Execute_CacheHit tests cache hit scenario.
func (s *GetPostUseCaseTestSuite) TestGetPostUseCase_Execute_CacheHit() {
	// Create test post
	repo := postgres.NewPostRepository(s.db)
	company, _ := domaincontent.NewCompanyName("缓存命中测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条缓存命中测试内容，用于验证缓存命中功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	err := repo.Save(s.ctx, post)
	require.NoError(s.T(), err)

	postID := post.ID().String()

	// First query (cache miss)
	result1, err := s.useCase.Execute(s.ctx, postID)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result1)

	// Verify cache was set
	cacheKey := "post:" + postID
	cachedData, err := s.redisClient.Get(s.ctx, cacheKey).Result()
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), cachedData)

	// Second query (should hit cache)
	result2, err := s.useCase.Execute(s.ctx, postID)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result2)

	// Verify results are the same
	assert.Equal(s.T(), result1.ID, result2.ID)
	assert.Equal(s.T(), result1.Company, result2.Company)
	assert.Equal(s.T(), result1.Content, result2.Content)
}

// TestGetPostUseCase_Execute_NotFound tests not found error.
func (s *GetPostUseCaseTestSuite) TestGetPostUseCase_Execute_NotFound() {
	// Use a non-existent post ID
	postID := "550e8400-e29b-41d4-a716-446655440000"

	// Query
	result, err := s.useCase.Execute(s.ctx, postID)

	// Assertions
	require.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.True(s.T(), apperrors.IsNotFoundError(err))
	assert.Contains(s.T(), err.Error(), "not found")
}

// TestGetPostUseCase_Execute_ValidationError tests validation error.
func (s *GetPostUseCaseTestSuite) TestGetPostUseCase_Execute_ValidationError() {
	testCases := []struct {
		name    string
		postID  string
		wantErr string
	}{
		{
			name:    "empty post ID",
			postID:  "",
			wantErr: "post ID is required",
		},
		{
			name:    "invalid UUID format",
			postID:  "invalid-uuid",
			wantErr: "invalid post ID",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result, err := s.useCase.Execute(s.ctx, tc.postID)

			require.Error(s.T(), err)
			assert.Nil(s.T(), result)
			assert.True(s.T(), apperrors.IsValidationError(err))
			assert.Contains(s.T(), err.Error(), tc.wantErr)
		})
	}
}

// TestGetPostUseCase_Execute_CacheTTL tests cache TTL.
func (s *GetPostUseCaseTestSuite) TestGetPostUseCase_Execute_CacheTTL() {
	// Create test post
	repo := postgres.NewPostRepository(s.db)
	company, _ := domaincontent.NewCompanyName("TTL测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条TTL测试内容，用于验证缓存TTL功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	err := repo.Save(s.ctx, post)
	require.NoError(s.T(), err)

	postID := post.ID().String()

	// Query (cache miss)
	_, err = s.useCase.Execute(s.ctx, postID)
	require.NoError(s.T(), err)

	// Verify cache TTL (should be 10 minutes)
	cacheKey := "post:" + postID
	ttl, err := s.redisClient.TTL(s.ctx, cacheKey).Result()
	require.NoError(s.T(), err)
	assert.Greater(s.T(), ttl, time.Duration(0))
	assert.LessOrEqual(s.T(), ttl, 10*time.Minute+time.Second) // Allow 1 second tolerance
}

// TestGetPostUseCase_Execute_Suite runs the test suite.
func TestGetPostUseCase_Execute_Suite(t *testing.T) {
	suite.Run(t, new(GetPostUseCaseTestSuite))
}

