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

// ListPostsUseCaseTestSuite is the test suite for ListPostsUseCase integration tests.
type ListPostsUseCaseTestSuite struct {
	suite.Suite
	db          *sql.DB
	redisClient *redisclient.Client
	useCase     *content.ListPostsUseCase
	ctx         context.Context
}

// SetupSuite runs once before all tests in the suite.
func (s *ListPostsUseCaseTestSuite) SetupSuite() {
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
	s.useCase = content.NewListPostsUseCase(postRepo, cacheRepo)

	// Create context
	s.ctx = context.Background()
}

// TearDownSuite runs once after all tests in the suite.
func (s *ListPostsUseCaseTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
}

// SetupTest runs before each test.
func (s *ListPostsUseCaseTestSuite) SetupTest() {
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
func (s *ListPostsUseCaseTestSuite) TearDownTest() {
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
func (s *ListPostsUseCaseTestSuite) waitForDB() error {
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
func (s *ListPostsUseCaseTestSuite) waitForRedis() error {
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
func (s *ListPostsUseCaseTestSuite) runMigrations() error {
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

// TestListPostsUseCase_Execute_Success tests successful list query.
func (s *ListPostsUseCaseTestSuite) TestListPostsUseCase_Execute_Success() {
	// Create test posts
	repo := postgres.NewPostRepository(s.db)
	company1, _ := domaincontent.NewCompanyName("测试公司1")
	company2, _ := domaincontent.NewCompanyName("测试公司2")
	city, _ := shared.NewCity("beijing", "北京")
	content1, _ := domaincontent.NewContent("这是一条测试内容1，用于验证列表查询功能。内容应该足够长以满足最小长度要求。")
	content2, _ := domaincontent.NewContent("这是一条测试内容2，用于验证列表查询功能。内容应该足够长以满足最小长度要求。")

	post1, _ := domaincontent.NewPost(company1, city, content1)
	post2, _ := domaincontent.NewPost(company2, city, content2)

	err := repo.Save(s.ctx, post1)
	require.NoError(s.T(), err)

	err = repo.Save(s.ctx, post2)
	require.NoError(s.T(), err)

	// Query
	query := content.ListPostsQuery{
		CityCode: "beijing",
		Page:     1,
		PageSize: 20,
	}

	result, err := s.useCase.Execute(s.ctx, query)

	// Assertions
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	assert.GreaterOrEqual(s.T(), result.Total, 2)
	assert.GreaterOrEqual(s.T(), len(result.Posts), 2)
	assert.Equal(s.T(), 1, result.Page)
	assert.Equal(s.T(), 20, result.PageSize)
}

// TestListPostsUseCase_Execute_CacheHit tests cache hit scenario.
func (s *ListPostsUseCaseTestSuite) TestListPostsUseCase_Execute_CacheHit() {
	// Create test post
	repo := postgres.NewPostRepository(s.db)
	company, _ := domaincontent.NewCompanyName("缓存测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条缓存测试内容，用于验证缓存命中功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	err := repo.Save(s.ctx, post)
	require.NoError(s.T(), err)

	// First query (cache miss)
	query := content.ListPostsQuery{
		CityCode: "beijing",
		Page:     1,
		PageSize: 20,
	}

	result1, err := s.useCase.Execute(s.ctx, query)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result1)

	// Verify cache was set
	cacheKey := "posts:city:beijing:page:1"
	cachedData, err := s.redisClient.Get(s.ctx, cacheKey).Result()
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), cachedData)

	// Second query (should hit cache)
	result2, err := s.useCase.Execute(s.ctx, query)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result2)

	// Verify results are the same
	assert.Equal(s.T(), result1.Total, result2.Total)
	assert.Equal(s.T(), len(result1.Posts), len(result2.Posts))
	if len(result1.Posts) > 0 && len(result2.Posts) > 0 {
		assert.Equal(s.T(), result1.Posts[0].ID, result2.Posts[0].ID)
	}
}

// TestListPostsUseCase_Execute_Pagination tests pagination.
func (s *ListPostsUseCaseTestSuite) TestListPostsUseCase_Execute_Pagination() {
	// Create multiple test posts
	repo := postgres.NewPostRepository(s.db)
	city, _ := shared.NewCity("beijing", "北京")

	for i := 0; i < 5; i++ {
		company, _ := domaincontent.NewCompanyName(fmt.Sprintf("分页测试公司%d", i))
		postContent, _ := domaincontent.NewContent(fmt.Sprintf("这是第%d条分页测试内容，用于验证分页功能。内容应该足够长以满足最小长度要求。", i))
		post, _ := domaincontent.NewPost(company, city, postContent)
		err := repo.Save(s.ctx, post)
		require.NoError(s.T(), err)
	}

	// Query page 1
	query1 := content.ListPostsQuery{
		CityCode: "beijing",
		Page:     1,
		PageSize: 2,
	}

	result1, err := s.useCase.Execute(s.ctx, query1)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, result1.Page)
	assert.Equal(s.T(), 2, result1.PageSize)
	assert.LessOrEqual(s.T(), len(result1.Posts), 2)

	// Query page 2
	query2 := content.ListPostsQuery{
		CityCode: "beijing",
		Page:     2,
		PageSize: 2,
	}

	result2, err := s.useCase.Execute(s.ctx, query2)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, result2.Page)
	assert.Equal(s.T(), 2, result2.PageSize)
	assert.LessOrEqual(s.T(), len(result2.Posts), 2)

	// Verify total is consistent
	assert.Equal(s.T(), result1.Total, result2.Total)
}

// TestListPostsUseCase_Execute_EmptyResult tests empty result.
func (s *ListPostsUseCaseTestSuite) TestListPostsUseCase_Execute_EmptyResult() {
	query := content.ListPostsQuery{
		CityCode: "hangzhou", // City with no posts
		Page:     1,
		PageSize: 20,
	}

	result, err := s.useCase.Execute(s.ctx, query)

	// Assertions
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	assert.Equal(s.T(), 0, result.Total)
	assert.Equal(s.T(), 0, len(result.Posts))
}

// TestListPostsUseCase_Execute_ValidationError tests validation error.
func (s *ListPostsUseCaseTestSuite) TestListPostsUseCase_Execute_ValidationError() {
	query := content.ListPostsQuery{
		CityCode: "", // Empty city code
		Page:     1,
		PageSize: 20,
	}

	result, err := s.useCase.Execute(s.ctx, query)

	// Assertions
	require.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.True(s.T(), apperrors.IsValidationError(err))
	assert.Contains(s.T(), err.Error(), "city code is required")
}

// TestListPostsUseCase_Execute_DefaultPagination tests default pagination values.
func (s *ListPostsUseCaseTestSuite) TestListPostsUseCase_Execute_DefaultPagination() {
	query := content.ListPostsQuery{
		CityCode: "beijing",
		Page:     0, // Will default to 1
		PageSize: 0, // Will default to 20
	}

	result, err := s.useCase.Execute(s.ctx, query)

	// Assertions
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	assert.Equal(s.T(), 1, result.Page)
	assert.Equal(s.T(), 20, result.PageSize)
}

// TestListPostsUseCase_Execute_CacheTTL tests cache TTL based on city popularity.
func (s *ListPostsUseCaseTestSuite) TestListPostsUseCase_Execute_CacheTTL() {
	// Create test post
	repo := postgres.NewPostRepository(s.db)
	company, _ := domaincontent.NewCompanyName("TTL测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条TTL测试内容，用于验证缓存TTL功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	err := repo.Save(s.ctx, post)
	require.NoError(s.T(), err)

	// Query (cache miss)
	query := content.ListPostsQuery{
		CityCode: "beijing",
		Page:     1,
		PageSize: 20,
	}

	_, err = s.useCase.Execute(s.ctx, query)
	require.NoError(s.T(), err)

	// Verify cache TTL (beijing is popular city, should be 5 minutes)
	cacheKey := "posts:city:beijing:page:1"
	ttl, err := s.redisClient.TTL(s.ctx, cacheKey).Result()
	require.NoError(s.T(), err)
	assert.Greater(s.T(), ttl, time.Duration(0))
	assert.LessOrEqual(s.T(), ttl, 5*time.Minute+time.Second) // Allow 1 second tolerance
}

// TestListPostsUseCase_Execute_Suite runs the test suite.
func TestListPostsUseCase_Execute_Suite(t *testing.T) {
	suite.Run(t, new(ListPostsUseCaseTestSuite))
}
