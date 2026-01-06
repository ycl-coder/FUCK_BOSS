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
	"fuck_boss/backend/internal/infrastructure/persistence/postgres"
	"fuck_boss/backend/internal/infrastructure/persistence/redis"
	apperrors "fuck_boss/backend/pkg/errors"
)

// CreatePostUseCaseTestSuite is the test suite for CreatePostUseCase integration tests.
type CreatePostUseCaseTestSuite struct {
	suite.Suite
	db          *sql.DB
	redisClient *redisclient.Client
	useCase     *content.CreatePostUseCase
	ctx         context.Context
}

// SetupSuite runs once before all tests in the suite.
func (s *CreatePostUseCaseTestSuite) SetupSuite() {
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
	rateLimiter := redis.NewRateLimiter(s.redisClient)

	// Create use case
	s.useCase = content.NewCreatePostUseCase(postRepo, cacheRepo, rateLimiter)

	// Create context
	s.ctx = context.Background()
}

// TearDownSuite runs once after all tests in the suite.
func (s *CreatePostUseCaseTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
}

// SetupTest runs before each test.
func (s *CreatePostUseCaseTestSuite) SetupTest() {
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
func (s *CreatePostUseCaseTestSuite) TearDownTest() {
	// Clean up database
	// _, err := s.db.ExecContext(s.ctx, "TRUNCATE TABLE posts CASCADE")
	// if err != nil {
	// 	s.T().Logf("Failed to truncate posts table: %v", err)
	// }

	// Clean up Redis
	// err = s.redisClient.FlushDB(s.ctx).Err()
	// if err != nil {
	// 	s.T().Logf("Failed to flush Redis database: %v", err)
	// }
}

// waitForDB waits for PostgreSQL to be ready.
func (s *CreatePostUseCaseTestSuite) waitForDB() error {
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
func (s *CreatePostUseCaseTestSuite) waitForRedis() error {
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
func (s *CreatePostUseCaseTestSuite) runMigrations() error {
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

	// Insert test cities (use ON CONFLICT to handle duplicates)
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

	// Create posts table with foreign key constraint
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

// TestCreatePostUseCase_Execute_Success tests successful post creation.
func (s *CreatePostUseCaseTestSuite) TestCreatePostUseCase_Execute_Success() {
	cmd := content.CreatePostCommand{
		Company:    "测试公司",
		CityCode:   "beijing",
		CityName:   "北京",
		Content:    "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP:   "127.0.0.1",
		OccurredAt: nil,
	}

	// Execute
	result, err := s.useCase.Execute(s.ctx, cmd)

	// Assertions
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	assert.NotEmpty(s.T(), result.ID)
	assert.Equal(s.T(), "测试公司", result.Company)
	assert.Equal(s.T(), "beijing", result.CityCode)
	assert.Equal(s.T(), "北京", result.CityName)
	assert.Equal(s.T(), cmd.Content, result.Content)
	assert.NotZero(s.T(), result.CreatedAt)

	// Verify saved in database
	var dbID, dbCompany, dbCityCode, dbCityName, dbContent string
	var dbCreatedAt time.Time
	err = s.db.QueryRowContext(s.ctx,
		"SELECT id, company_name, city_code, city_name, content, created_at FROM posts WHERE id = $1",
		result.ID,
	).Scan(&dbID, &dbCompany, &dbCityCode, &dbCityName, &dbContent, &dbCreatedAt)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), result.ID, dbID)
	assert.Equal(s.T(), result.Company, dbCompany)
	assert.Equal(s.T(), result.CityCode, dbCityCode)
	assert.Equal(s.T(), result.CityName, dbCityName)
	assert.Equal(s.T(), result.Content, dbContent)
}

// TestCreatePostUseCase_Execute_CacheCleared tests that cache is cleared after post creation.
func (s *CreatePostUseCaseTestSuite) TestCreatePostUseCase_Execute_CacheCleared() {
	// Setup: Create some cache entries that should be cleared
	cacheKey1 := "posts:city:beijing:page:1"
	cacheKey2 := "posts:city:beijing:page:2"
	cacheKey3 := "posts:city:shanghai:page:1" // Should not be cleared

	err := s.redisClient.Set(s.ctx, cacheKey1, "cached_data_1", time.Hour).Err()
	require.NoError(s.T(), err)

	err = s.redisClient.Set(s.ctx, cacheKey2, "cached_data_2", time.Hour).Err()
	require.NoError(s.T(), err)

	err = s.redisClient.Set(s.ctx, cacheKey3, "cached_data_3", time.Hour).Err()
	require.NoError(s.T(), err)

	// Verify cache exists
	val1, err := s.redisClient.Get(s.ctx, cacheKey1).Result()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "cached_data_1", val1)

	// Create post
	cmd := content.CreatePostCommand{
		Company:  "测试公司",
		CityCode: "beijing",
		CityName: "北京",
		Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP: "127.0.0.1",
	}

	_, err = s.useCase.Execute(s.ctx, cmd)
	require.NoError(s.T(), err)

	// Verify beijing cache is cleared
	val1, err = s.redisClient.Get(s.ctx, cacheKey1).Result()
	assert.Error(s.T(), err) // Should be deleted
	assert.Empty(s.T(), val1)

	val2, err := s.redisClient.Get(s.ctx, cacheKey2).Result()
	assert.Error(s.T(), err) // Should be deleted
	assert.Empty(s.T(), val2)

	// Verify shanghai cache is still there
	val3, err := s.redisClient.Get(s.ctx, cacheKey3).Result()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "cached_data_3", val3)
}

// TestCreatePostUseCase_Execute_RateLimit tests rate limiting.
func (s *CreatePostUseCaseTestSuite) TestCreatePostUseCase_Execute_RateLimit() {
	cmd := content.CreatePostCommand{
		Company:  "测试公司",
		CityCode: "beijing",
		CityName: "北京",
		Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP: "192.168.1.100", // Use unique IP for this test
	}

	// Create 3 posts (within limit)
	for i := 0; i < 3; i++ {
		cmd.Company = fmt.Sprintf("测试公司%d", i)
		result, err := s.useCase.Execute(s.ctx, cmd)
		require.NoError(s.T(), err, "Post %d should succeed", i)
		assert.NotEmpty(s.T(), result.ID)
	}

	// 4th post should fail (rate limit exceeded)
	cmd.Company = "测试公司4"
	result, err := s.useCase.Execute(s.ctx, cmd)
	require.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.True(s.T(), apperrors.IsRateLimitError(err), "Error should be RateLimitError")
	assert.Contains(s.T(), err.Error(), "rate limit exceeded")
}

// TestCreatePostUseCase_Execute_ValidationError tests validation errors.
func (s *CreatePostUseCaseTestSuite) TestCreatePostUseCase_Execute_ValidationError() {
	testCases := []struct {
		name    string
		cmd     content.CreatePostCommand
		wantErr string
	}{
		{
			name: "empty company",
			cmd: content.CreatePostCommand{
				Company:  "",
				CityCode: "beijing",
				CityName: "北京",
				Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
				ClientIP: "127.0.0.1",
			},
			wantErr: "company name is required",
		},
		{
			name: "empty city code",
			cmd: content.CreatePostCommand{
				Company:  "测试公司",
				CityCode: "",
				CityName: "北京",
				Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
				ClientIP: "127.0.0.1",
			},
			wantErr: "city code is required",
		},
		{
			name: "empty content",
			cmd: content.CreatePostCommand{
				Company:  "测试公司",
				CityCode: "beijing",
				CityName: "北京",
				Content:  "",
				ClientIP: "127.0.0.1",
			},
			wantErr: "content is required",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result, err := s.useCase.Execute(s.ctx, tc.cmd)

			require.Error(s.T(), err)
			assert.Nil(s.T(), result)
			assert.True(s.T(), apperrors.IsValidationError(err))
			assert.Contains(s.T(), err.Error(), tc.wantErr)
		})
	}
}

// TestCreatePostUseCase_Execute_InvalidValueObjects tests invalid value objects.
func (s *CreatePostUseCaseTestSuite) TestCreatePostUseCase_Execute_InvalidValueObjects() {
	testCases := []struct {
		name    string
		cmd     content.CreatePostCommand
		wantErr string
	}{
		{
			name: "company name too long",
			cmd: content.CreatePostCommand{
				Company:  string(make([]byte, 101)), // 101 characters
				CityCode: "beijing",
				CityName: "北京",
				Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
				ClientIP: "127.0.0.1",
			},
			wantErr: "invalid company name",
		},
		{
			name: "content too short",
			cmd: content.CreatePostCommand{
				Company:  "测试公司",
				CityCode: "beijing",
				CityName: "北京",
				Content:  "太短", // Less than 10 characters
				ClientIP: "127.0.0.1",
			},
			wantErr: "invalid content",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result, err := s.useCase.Execute(s.ctx, tc.cmd)

			require.Error(s.T(), err)
			assert.Nil(s.T(), result)
			assert.True(s.T(), apperrors.IsValidationError(err))
			assert.Contains(s.T(), err.Error(), tc.wantErr)
		})
	}
}

// TestCreatePostUseCase_Execute_CompleteFlow tests the complete flow with all components.
func (s *CreatePostUseCaseTestSuite) TestCreatePostUseCase_Execute_CompleteFlow() {
	// Setup: Create cache entries
	cacheKey := "posts:city:beijing:page:1"
	err := s.redisClient.Set(s.ctx, cacheKey, "cached_data", time.Hour).Err()
	require.NoError(s.T(), err)

	// Create post
	cmd := content.CreatePostCommand{
		Company:  "完整流程测试公司",
		CityCode: "beijing",
		CityName: "北京",
		Content:  "这是一条完整的流程测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP: "192.168.1.200",
	}

	result, err := s.useCase.Execute(s.ctx, cmd)

	// Verify success
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	assert.NotEmpty(s.T(), result.ID)

	// Verify saved in database
	var count int
	err = s.db.QueryRowContext(s.ctx,
		"SELECT COUNT(*) FROM posts WHERE id = $1",
		result.ID,
	).Scan(&count)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, count)

	// Verify cache is cleared
	val, err := s.redisClient.Get(s.ctx, cacheKey).Result()
	assert.Error(s.T(), err) // Should be deleted
	assert.Empty(s.T(), val)

	// Verify rate limit key exists
	rateLimitKey := fmt.Sprintf("rate_limit:post:%s:%s", cmd.ClientIP, time.Now().Format("2006-01-02-15"))
	val, err = s.redisClient.Get(s.ctx, rateLimitKey).Result()
	// Rate limit key might not exist if window expired, but if it exists, it should have a value
	if err == nil {
		assert.NotEmpty(s.T(), val)
	}
}

// TestCreatePostUseCase_Execute_Suite runs the test suite.
func TestCreatePostUseCase_Execute_Suite(t *testing.T) {
	suite.Run(t, new(CreatePostUseCaseTestSuite))
}
