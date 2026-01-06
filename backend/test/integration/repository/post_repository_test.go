// Package repository provides integration tests for PostgreSQL repositories.
// These tests use a real PostgreSQL database to verify repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"fuck_boss/backend/internal/domain/content"
	"fuck_boss/backend/internal/domain/shared"
	"fuck_boss/backend/internal/infrastructure/persistence/postgres"
	apperrors "fuck_boss/backend/pkg/errors"
)

// PostRepositoryTestSuite is the test suite for PostRepository integration tests.
type PostRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo *postgres.PostRepository
	ctx  context.Context
}

// SetupSuite runs once before all tests in the suite.
func (s *PostRepositoryTestSuite) SetupSuite() {
	// Get database connection string from environment or use default
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://test_user:test_password@localhost:5433/test_db?sslmode=disable"
	}

	// Connect to database
	var err error
	s.db, err = sql.Open("postgres", dsn)
	require.NoError(s.T(), err, "Failed to connect to test database")

	// Wait for database to be ready
	err = s.waitForDB()
	require.NoError(s.T(), err, "Database is not ready")

	// Run migrations
	err = s.runMigrations()
	require.NoError(s.T(), err, "Failed to run migrations")

	// Create repository
	s.repo = postgres.NewPostRepository(s.db)

	// Create context
	s.ctx = context.Background()
}

// TearDownSuite runs once after all tests in the suite.
func (s *PostRepositoryTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

// SetupTest runs before each test.
func (s *PostRepositoryTestSuite) SetupTest() {
	// Clean up any existing test data before each test
	_, err := s.db.ExecContext(s.ctx, "TRUNCATE TABLE posts CASCADE")
	if err != nil {
		s.T().Logf("Failed to truncate posts table: %v", err)
	}
}

// TearDownTest runs after each test.
func (s *PostRepositoryTestSuite) TearDownTest() {
	// Clean up test data after each test
	// Note: Comment out the TRUNCATE below if you want to inspect test data
	// After inspection, uncomment it to ensure test isolation
	// _, err := s.db.ExecContext(s.ctx, "TRUNCATE TABLE posts CASCADE")
	// if err != nil {
	// 	s.T().Logf("Failed to truncate posts table: %v", err)
	// }
}

// waitForDB waits for the database to be ready.
func (s *PostRepositoryTestSuite) waitForDB() error {
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

// runMigrations runs database migrations.
func (s *PostRepositoryTestSuite) runMigrations() error {
	// Create context for migration (s.ctx is not set yet in SetupSuite)
	ctx := context.Background()

	// Read migration file
	migrationSQL := `
		CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			company_name VARCHAR(100) NOT NULL,
			city_code VARCHAR(50) NOT NULL,
			city_name VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			occurred_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_posts_city_code ON posts(city_code);
		CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_posts_company_name ON posts(company_name);
		CREATE INDEX IF NOT EXISTS idx_posts_search ON posts USING GIN(
			to_tsvector('simple', company_name || ' ' || content)
		);
	`

	_, err := s.db.ExecContext(ctx, migrationSQL)
	return err
}

// TestPostRepository_Save tests the Save method.
func (s *PostRepositoryTestSuite) TestPostRepository_Save() {
	// Create test post
	company, err := content.NewCompanyName("测试公司")
	s.Require().NoError(err)

	city, err := shared.NewCity("beijing", "北京")
	s.Require().NoError(err)

	postContent, err := content.NewContent("这是一条测试内容，用于验证 Save 方法的功能。内容应该足够长以满足最小长度要求。")
	s.Require().NoError(err)

	post, err := content.NewPost(company, city, postContent)
	s.Require().NoError(err)
	s.Require().NotNil(post)

	// Save post
	err = s.repo.Save(s.ctx, post)
	s.Require().NoError(err)

	// Verify saved
	found, err := s.repo.FindByID(s.ctx, post.ID())
	s.Require().NoError(err)
	s.Require().NotNil(found)
	s.Equal(post.ID().String(), found.ID().String())
	s.Equal(post.Company().String(), found.Company().String())
	s.Equal(post.City().Code(), found.City().Code())
	s.Equal(post.City().Name(), found.City().Name())
	s.Equal(post.Content().String(), found.Content().String())
}

// TestPostRepository_Save_Update tests updating an existing post.
func (s *PostRepositoryTestSuite) TestPostRepository_Save_Update() {
	// Create and save initial post
	company1, _ := content.NewCompanyName("公司A")
	city1, _ := shared.NewCity("beijing", "北京")
	content1, _ := content.NewContent("这是初始内容，用于测试更新功能。内容应该足够长以满足最小长度要求。")
	post1, _ := content.NewPost(company1, city1, content1)

	err := s.repo.Save(s.ctx, post1)
	s.Require().NoError(err)

	// Update post with new content
	company2, _ := content.NewCompanyName("公司B")
	city2, _ := shared.NewCity("shanghai", "上海")
	content2, _ := content.NewContent("这是更新后的内容，用于验证 Save 方法能够更新已存在的记录。内容应该足够长以满足最小长度要求。")
	post2, _ := content.NewPostFromDB(post1.ID(), company2, city2, content2, post1.CreatedAt())

	err = s.repo.Save(s.ctx, post2)
	s.Require().NoError(err)

	// Verify updated
	found, err := s.repo.FindByID(s.ctx, post1.ID())
	s.Require().NoError(err)
	s.Require().NotNil(found)
	s.Equal("公司B", found.Company().String())
	s.Equal("shanghai", found.City().Code())
	s.Equal("上海", found.City().Name())
}

// TestPostRepository_FindByID tests the FindByID method.
func (s *PostRepositoryTestSuite) TestPostRepository_FindByID() {
	// Create and save test post
	company, _ := content.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent("这是用于测试 FindByID 方法的内容。内容应该足够长以满足最小长度要求。")
	post, _ := content.NewPost(company, city, postContent)

	err := s.repo.Save(s.ctx, post)
	s.Require().NoError(err)

	// Find by ID
	found, err := s.repo.FindByID(s.ctx, post.ID())
	s.Require().NoError(err)
	s.Require().NotNil(found)
	s.Equal(post.ID().String(), found.ID().String())
}

// TestPostRepository_FindByID_NotFound tests FindByID with non-existent ID.
func (s *PostRepositoryTestSuite) TestPostRepository_FindByID_NotFound() {
	// Generate a non-existent ID
	nonExistentID := content.GeneratePostID()

	// Try to find
	found, err := s.repo.FindByID(s.ctx, nonExistentID)
	s.Require().Error(err)
	s.Require().Nil(found)
	s.True(apperrors.IsNotFoundError(err), "Error should be NotFoundError")
}

// TestPostRepository_FindByCity tests the FindByCity method.
func (s *PostRepositoryTestSuite) TestPostRepository_FindByCity() {
	// Create test posts in different cities
	beijing, _ := shared.NewCity("beijing", "北京")
	shanghai, _ := shared.NewCity("shanghai", "上海")

	// Create posts in Beijing
	for i := 0; i < 5; i++ {
		company, _ := content.NewCompanyName(fmt.Sprintf("北京公司%d", i))
		postContent, _ := content.NewContent(fmt.Sprintf("这是北京的第%d条测试内容。内容应该足够长以满足最小长度要求。", i))
		post, _ := content.NewPost(company, beijing, postContent)
		err := s.repo.Save(s.ctx, post)
		s.Require().NoError(err)
		// Add small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// Create posts in Shanghai
	for i := 0; i < 3; i++ {
		company, _ := content.NewCompanyName(fmt.Sprintf("上海公司%d", i))
		postContent, _ := content.NewContent(fmt.Sprintf("这是上海的第%d条测试内容。内容应该足够长以满足最小长度要求。", i))
		post, _ := content.NewPost(company, shanghai, postContent)
		err := s.repo.Save(s.ctx, post)
		s.Require().NoError(err)
		time.Sleep(10 * time.Millisecond)
	}

	// Find posts in Beijing
	posts, total, err := s.repo.FindByCity(s.ctx, beijing, 1, 10)
	s.Require().NoError(err)
	s.Equal(5, total)
	s.Len(posts, 5)

	// Verify all posts are from Beijing
	for _, post := range posts {
		s.Equal("beijing", post.City().Code())
		s.Equal("北京", post.City().Name())
	}

	// Verify ordering (should be DESC by created_at)
	for i := 0; i < len(posts)-1; i++ {
		s.True(posts[i].CreatedAt().After(posts[i+1].CreatedAt()) || posts[i].CreatedAt().Equal(posts[i+1].CreatedAt()),
			"Posts should be ordered by created_at DESC")
	}
}

// TestPostRepository_FindByCity_Pagination tests pagination in FindByCity.
func (s *PostRepositoryTestSuite) TestPostRepository_FindByCity_Pagination() {
	beijing, _ := shared.NewCity("beijing", "北京")

	// Create 15 posts
	for i := 0; i < 15; i++ {
		company, _ := content.NewCompanyName(fmt.Sprintf("公司%d", i))
		postContent, _ := content.NewContent(fmt.Sprintf("这是第%d条测试内容。内容应该足够长以满足最小长度要求。", i))
		post, _ := content.NewPost(company, beijing, postContent)
		err := s.repo.Save(s.ctx, post)
		s.Require().NoError(err)
		time.Sleep(10 * time.Millisecond)
	}

	// Test first page
	posts1, total1, err := s.repo.FindByCity(s.ctx, beijing, 1, 10)
	s.Require().NoError(err)
	s.Equal(15, total1)
	s.Len(posts1, 10)

	// Test second page
	posts2, total2, err := s.repo.FindByCity(s.ctx, beijing, 2, 10)
	s.Require().NoError(err)
	s.Equal(15, total2)
	s.Len(posts2, 5)

	// Verify no overlap
	ids1 := make(map[string]bool)
	for _, post := range posts1 {
		ids1[post.ID().String()] = true
	}
	for _, post := range posts2 {
		s.False(ids1[post.ID().String()], "Posts should not overlap between pages")
	}
}

// TestPostRepository_Search tests the Search method.
func (s *PostRepositoryTestSuite) TestPostRepository_Search() {
	beijing, _ := shared.NewCity("beijing", "北京")
	shanghai, _ := shared.NewCity("shanghai", "上海")

	// Create posts with different keywords
	company1, _ := content.NewCompanyName("阿里巴巴")
	content1, _ := content.NewContent("这是一条关于阿里巴巴的测试内容。内容应该足够长以满足最小长度要求。")
	post1, _ := content.NewPost(company1, beijing, content1)
	s.repo.Save(s.ctx, post1)

	company2, _ := content.NewCompanyName("腾讯公司")
	content2, _ := content.NewContent("这是一条关于腾讯的测试内容。内容应该足够长以满足最小长度要求。")
	post2, _ := content.NewPost(company2, shanghai, content2)
	s.repo.Save(s.ctx, post2)

	company3, _ := content.NewCompanyName("百度公司")
	content3, _ := content.NewContent("这是一条关于百度的测试内容。内容应该足够长以满足最小长度要求。")
	post3, _ := content.NewPost(company3, beijing, content3)
	s.repo.Save(s.ctx, post3)

	// Search for "阿里巴巴"
	posts, total, err := s.repo.Search(s.ctx, "阿里巴巴", nil, 1, 10)
	s.Require().NoError(err)
	s.GreaterOrEqual(total, 1)
	s.GreaterOrEqual(len(posts), 1)
	s.Contains(posts[0].Company().String(), "阿里巴巴")
}

// TestPostRepository_Search_WithCityFilter tests Search with city filter.
func (s *PostRepositoryTestSuite) TestPostRepository_Search_WithCityFilter() {
	beijing, _ := shared.NewCity("beijing", "北京")
	shanghai, _ := shared.NewCity("shanghai", "上海")

	// Create posts in different cities with same keyword
	company1, _ := content.NewCompanyName("测试公司")
	content1, _ := content.NewContent("这是一条测试内容，包含关键词：测试。内容应该足够长以满足最小长度要求。")
	post1, _ := content.NewPost(company1, beijing, content1)
	s.repo.Save(s.ctx, post1)

	company2, _ := content.NewCompanyName("测试公司")
	content2, _ := content.NewContent("这是一条测试内容，包含关键词：测试。内容应该足够长以满足最小长度要求。")
	post2, _ := content.NewPost(company2, shanghai, content2)
	s.repo.Save(s.ctx, post2)

	// Search with city filter
	posts, total, err := s.repo.Search(s.ctx, "测试", &beijing, 1, 10)
	s.Require().NoError(err)
	s.GreaterOrEqual(total, 1)
	s.GreaterOrEqual(len(posts), 1)

	// Verify all results are from Beijing
	for _, post := range posts {
		s.Equal("beijing", post.City().Code())
	}
}

// TestPostRepository_Search_Pagination tests pagination in Search.
func (s *PostRepositoryTestSuite) TestPostRepository_Search_Pagination() {
	beijing, _ := shared.NewCity("beijing", "北京")

	// Create multiple posts with same keyword
	for i := 0; i < 12; i++ {
		company, _ := content.NewCompanyName(fmt.Sprintf("测试公司%d", i))
		postContent, _ := content.NewContent(fmt.Sprintf("这是第%d条测试内容，包含关键词：测试。内容应该足够长以满足最小长度要求。", i))
		post, _ := content.NewPost(company, beijing, postContent)
		s.repo.Save(s.ctx, post)
		time.Sleep(10 * time.Millisecond)
	}

	// Test first page
	posts1, total1, err := s.repo.Search(s.ctx, "测试", &beijing, 1, 5)
	s.Require().NoError(err)
	s.GreaterOrEqual(total1, 12)
	s.Len(posts1, 5)

	// Test second page
	posts2, total2, err := s.repo.Search(s.ctx, "测试", &beijing, 2, 5)
	s.Require().NoError(err)
	s.Equal(total1, total2)
	s.Len(posts2, 5)
}

// TestPostRepository_ErrorHandling tests error handling.
func (s *PostRepositoryTestSuite) TestPostRepository_ErrorHandling() {
	// Test with invalid context (cancelled)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	company, _ := content.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent("这是测试内容。内容应该足够长以满足最小长度要求。")
	post, _ := content.NewPost(company, city, postContent)

	err := s.repo.Save(ctx, post)
	s.Require().Error(err)
}

// TestPostRepositorySuite runs all tests in the suite.
func TestPostRepositorySuite(t *testing.T) {
	suite.Run(t, new(PostRepositoryTestSuite))
}
