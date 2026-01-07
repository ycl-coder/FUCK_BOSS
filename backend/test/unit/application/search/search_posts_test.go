// Package search_test provides unit tests for search use cases.
// These tests use mocked dependencies to isolate the use case logic.
package search_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"fuck_boss/backend/internal/application/dto"
	"fuck_boss/backend/internal/application/search"
	domaincontent "fuck_boss/backend/internal/domain/content"
	"fuck_boss/backend/internal/domain/shared"
	apperrors "fuck_boss/backend/pkg/errors"
)

// MockPostRepository is a mock implementation of PostRepository.
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Save(ctx context.Context, post *domaincontent.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) FindByID(ctx context.Context, id domaincontent.PostID) (*domaincontent.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domaincontent.Post), args.Error(1)
}

func (m *MockPostRepository) FindByCity(ctx context.Context, city shared.City, page, pageSize int) ([]*domaincontent.Post, int, error) {
	args := m.Called(ctx, city, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domaincontent.Post), args.Int(1), args.Error(2)
}

func (m *MockPostRepository) Search(ctx context.Context, keyword string, city *shared.City, page, pageSize int) ([]*domaincontent.Post, int, error) {
	args := m.Called(ctx, keyword, city, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domaincontent.Post), args.Int(1), args.Error(2)
}

func (m *MockPostRepository) FindAll(ctx context.Context, page, pageSize int) ([]*domaincontent.Post, int, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domaincontent.Post), args.Int(1), args.Error(2)
}

// MockCacheRepository is a mock implementation of CacheRepository.
type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheRepository) DeleteByPattern(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}

// TestSearchPostsUseCase_Execute_CacheHit tests cache hit scenario.
func TestSearchPostsUseCase_Execute_CacheHit(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()
	query := search.SearchPostsQuery{
		Keyword:  "测试",
		CityCode: nil,
		Page:     1,
		PageSize: 20,
	}

	// Create cached data
	cachedResult := &dto.PostsListDTO{
		Posts: []*dto.PostDTO{
			{
				ID:        "test-id-1",
				Company:   "测试公司1",
				CityCode:  "beijing",
				CityName:  "北京",
				Content:   "测试内容1",
				CreatedAt: time.Now(),
			},
		},
		Total:    1,
		Page:     1,
		PageSize: 20,
	}
	cachedData, _ := json.Marshal(cachedResult)

	// Setup expectations
	mockCache.On("Get", ctx, "search:测试:page:1").Return(string(cachedData), nil)

	// Execute
	result, err := uc.Execute(ctx, query)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, 1, len(result.Posts))

	// Verify repository was not called
	mockRepo.AssertNotCalled(t, "Search")
	mockCache.AssertExpectations(t)
}

// TestSearchPostsUseCase_Execute_CacheMiss tests cache miss scenario.
func TestSearchPostsUseCase_Execute_CacheMiss(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()
	query := search.SearchPostsQuery{
		Keyword:  "测试",
		CityCode: nil,
		Page:     1,
		PageSize: 20,
	}

	// Create test posts
	company, _ := domaincontent.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条测试内容，用于验证搜索功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	// Setup expectations
	mockCache.On("Get", ctx, "search:测试:page:1").Return("", errors.New("cache miss"))
	mockRepo.On("Search", ctx, "测试", (*shared.City)(nil), 1, 20).
		Return([]*domaincontent.Post{post}, 1, nil)
	mockCache.On("Set", ctx, "search:测试:page:1", mock.AnythingOfType("string"), 5*time.Minute).Return(nil)

	// Execute
	result, err := uc.Execute(ctx, query)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, 1, len(result.Posts))

	// Verify all expectations
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestSearchPostsUseCase_Execute_ValidationError tests validation errors.
func TestSearchPostsUseCase_Execute_ValidationError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()

	testCases := []struct {
		name    string
		query   search.SearchPostsQuery
		wantErr string
	}{
		{
			name: "empty keyword",
			query: search.SearchPostsQuery{
				Keyword:  "",
				CityCode: nil,
				Page:     1,
				PageSize: 20,
			},
			wantErr: "keyword is required",
		},
		{
			name: "keyword too short",
			query: search.SearchPostsQuery{
				Keyword:  "a",
				CityCode: nil,
				Page:     1,
				PageSize: 20,
			},
			wantErr: "keyword must be at least 2 characters",
		},
		{
			name: "keyword with only whitespace",
			query: search.SearchPostsQuery{
				Keyword:  "  ",
				CityCode: nil,
				Page:     1,
				PageSize: 20,
			},
			wantErr: "keyword is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute
			result, err := uc.Execute(ctx, tc.query)

			// Assertions
			require.Error(t, err)
			assert.Nil(t, result)
			assert.True(t, apperrors.IsValidationError(err))
			assert.Contains(t, err.Error(), tc.wantErr)

			// Verify no methods were called
			mockRepo.AssertNotCalled(t, "Search")
			mockCache.AssertNotCalled(t, "Get")
		})
	}
}

// TestSearchPostsUseCase_Execute_WithCityFilter tests search with city filter.
func TestSearchPostsUseCase_Execute_WithCityFilter(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()
	cityCode := "beijing"
	query := search.SearchPostsQuery{
		Keyword:  "测试",
		CityCode: &cityCode,
		Page:     1,
		PageSize: 20,
	}

	// Create test post
	company, _ := domaincontent.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条测试内容，用于验证搜索功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	// Setup expectations
	mockCache.On("Get", ctx, "search:测试:city:beijing:page:1").Return("", errors.New("cache miss"))
	mockRepo.On("Search", ctx, "测试", mock.AnythingOfType("*shared.City"), 1, 20).
		Return([]*domaincontent.Post{post}, 1, nil)
	mockCache.On("Set", ctx, "search:测试:city:beijing:page:1", mock.AnythingOfType("string"), 5*time.Minute).Return(nil)

	// Execute
	result, err := uc.Execute(ctx, query)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)

	// Verify all expectations
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestSearchPostsUseCase_Execute_KeywordNormalization tests keyword normalization in cache key.
func TestSearchPostsUseCase_Execute_KeywordNormalization(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()

	// Test with uppercase and whitespace
	query := search.SearchPostsQuery{
		Keyword:  "  TEST  ",
		CityCode: nil,
		Page:     1,
		PageSize: 20,
	}

	// Create test post
	company, _ := domaincontent.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条测试内容，用于验证搜索功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	// Setup expectations - cache key should be normalized (lowercase, trimmed)
	mockCache.On("Get", ctx, "search:test:page:1").Return("", errors.New("cache miss"))
	mockRepo.On("Search", ctx, "  TEST  ", (*shared.City)(nil), 1, 20).
		Return([]*domaincontent.Post{post}, 1, nil)
	mockCache.On("Set", ctx, "search:test:page:1", mock.AnythingOfType("string"), 5*time.Minute).Return(nil)

	// Execute
	result, err := uc.Execute(ctx, query)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify cache key was normalized
	mockCache.AssertExpectations(t)
}

// TestSearchPostsUseCase_Execute_RepositoryError tests repository error.
func TestSearchPostsUseCase_Execute_RepositoryError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()
	query := search.SearchPostsQuery{
		Keyword:  "测试",
		CityCode: nil,
		Page:     1,
		PageSize: 20,
	}

	// Setup expectations
	mockCache.On("Get", ctx, "search:测试:page:1").Return("", errors.New("cache miss"))
	mockRepo.On("Search", ctx, "测试", (*shared.City)(nil), 1, 20).
		Return(nil, 0, errors.New("database connection failed"))

	// Execute
	result, err := uc.Execute(ctx, query)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, apperrors.IsDatabaseError(err))

	// Verify cache was not updated
	mockCache.AssertNotCalled(t, "Set")
	mockRepo.AssertExpectations(t)
}

// TestSearchPostsUseCase_Execute_CacheError tests cache error handling (should fallback to database).
func TestSearchPostsUseCase_Execute_CacheError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()
	query := search.SearchPostsQuery{
		Keyword:  "测试",
		CityCode: nil,
		Page:     1,
		PageSize: 20,
	}

	// Create test post
	company, _ := domaincontent.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条测试内容，用于验证搜索功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	// Setup expectations - cache error but should fallback to database
	mockCache.On("Get", ctx, "search:测试:page:1").Return("", errors.New("redis connection failed"))
	mockRepo.On("Search", ctx, "测试", (*shared.City)(nil), 1, 20).
		Return([]*domaincontent.Post{post}, 1, nil)
	mockCache.On("Set", ctx, "search:测试:page:1", mock.AnythingOfType("string"), 5*time.Minute).Return(nil)

	// Execute
	result, err := uc.Execute(ctx, query)

	// Assertions - should succeed despite cache error
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)

	// Verify all expectations
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestSearchPostsUseCase_Execute_Pagination tests pagination.
func TestSearchPostsUseCase_Execute_Pagination(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()

	testCases := []struct {
		name     string
		query    search.SearchPostsQuery
		expected struct {
			page     int
			pageSize int
		}
	}{
		{
			name: "default pagination",
			query: search.SearchPostsQuery{
				Keyword:  "测试",
				CityCode: nil,
				Page:     0, // Will default to 1
				PageSize: 0, // Will default to 20
			},
			expected: struct {
				page     int
				pageSize int
			}{page: 1, pageSize: 20},
		},
		{
			name: "custom pagination",
			query: search.SearchPostsQuery{
				Keyword:  "测试",
				CityCode: nil,
				Page:     2,
				PageSize: 10,
			},
			expected: struct {
				page     int
				pageSize int
			}{page: 2, pageSize: 10},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup expectations
			mockCache.On("Get", ctx, mock.AnythingOfType("string")).Return("", errors.New("cache miss"))
			mockRepo.On("Search", ctx, "测试", (*shared.City)(nil), tc.expected.page, tc.expected.pageSize).
				Return([]*domaincontent.Post{}, 0, nil)
			mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string"), 5*time.Minute).Return(nil)

			// Execute
			result, err := uc.Execute(ctx, tc.query)

			// Assertions
			require.NoError(t, err)
			assert.Equal(t, tc.expected.page, result.Page)
			assert.Equal(t, tc.expected.pageSize, result.PageSize)
		})
	}
}

// TestSearchPostsUseCase_Execute_EmptyResult tests empty result.
func TestSearchPostsUseCase_Execute_EmptyResult(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()
	query := search.SearchPostsQuery{
		Keyword:  "不存在",
		CityCode: nil,
		Page:     1,
		PageSize: 20,
	}

	// Setup expectations
	mockCache.On("Get", ctx, "search:不存在:page:1").Return("", errors.New("cache miss"))
	mockRepo.On("Search", ctx, "不存在", (*shared.City)(nil), 1, 20).
		Return([]*domaincontent.Post{}, 0, nil)
	mockCache.On("Set", ctx, "search:不存在:page:1", mock.AnythingOfType("string"), 5*time.Minute).Return(nil)

	// Execute
	result, err := uc.Execute(ctx, query)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, 0, len(result.Posts))

	// Verify all expectations
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestSearchPostsUseCase_Execute_InvalidCityCode tests invalid city code.
func TestSearchPostsUseCase_Execute_InvalidCityCode(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := search.NewSearchPostsUseCase(mockRepo, mockCache)

	ctx := context.Background()
	cityCode := "invalid-city"
	query := search.SearchPostsQuery{
		Keyword:  "测试",
		CityCode: &cityCode,
		Page:     1,
		PageSize: 20,
	}

	// Setup expectations - city code will be used as name (fallback)
	mockCache.On("Get", ctx, "search:测试:city:invalid-city:page:1").Return("", errors.New("cache miss"))
	mockRepo.On("Search", ctx, "测试", mock.AnythingOfType("*shared.City"), 1, 20).
		Return([]*domaincontent.Post{}, 0, nil)
	mockCache.On("Set", ctx, "search:测试:city:invalid-city:page:1", mock.AnythingOfType("string"), 5*time.Minute).Return(nil)

	// Execute
	result, err := uc.Execute(ctx, query)

	// Assertions - should succeed (city code used as name fallback)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify all expectations
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
