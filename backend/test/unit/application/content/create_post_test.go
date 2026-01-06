// Package content_test provides unit tests for content use cases.
// These tests use mocked dependencies to isolate the use case logic.
package content_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"fuck_boss/backend/internal/application/content"
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

// MockRateLimiter is a mock implementation of RateLimiter.
type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	args := m.Called(ctx, key, limit, window)
	return args.Bool(0), args.Error(1)
}

// TestCreatePostUseCase_Execute_Success tests successful post creation.
func TestCreatePostUseCase_Execute_Success(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)
	mockRateLimiter := new(MockRateLimiter)

	// Create use case
	uc := content.NewCreatePostUseCase(mockRepo, mockCache, mockRateLimiter)

	ctx := context.Background()
	cmd := content.CreatePostCommand{
		Company:    "测试公司",
		CityCode:   "beijing",
		CityName:   "北京",
		Content:    "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP:   "127.0.0.1",
		OccurredAt: nil,
	}

	// Setup expectations
	mockRateLimiter.On("Allow", ctx, mock.AnythingOfType("string"), 3, time.Hour).Return(true, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*content.Post")).Return(nil)
	mockCache.On("DeleteByPattern", ctx, "posts:city:beijing:*").Return(nil)

	// Execute
	result, err := uc.Execute(ctx, cmd)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, "测试公司", result.Company)
	assert.Equal(t, "beijing", result.CityCode)
	assert.Equal(t, "北京", result.CityName)
	assert.Equal(t, cmd.Content, result.Content)
	assert.NotZero(t, result.CreatedAt)

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockRateLimiter.AssertExpectations(t)
}

// TestCreatePostUseCase_Execute_ValidationError tests validation errors.
func TestCreatePostUseCase_Execute_ValidationError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)
	mockRateLimiter := new(MockRateLimiter)

	// Create use case
	uc := content.NewCreatePostUseCase(mockRepo, mockCache, mockRateLimiter)

	ctx := context.Background()

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
			name: "empty city name",
			cmd: content.CreatePostCommand{
				Company:  "测试公司",
				CityCode: "beijing",
				CityName: "",
				Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
				ClientIP: "127.0.0.1",
			},
			wantErr: "city name is required",
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
		{
			name: "empty client IP",
			cmd: content.CreatePostCommand{
				Company:  "测试公司",
				CityCode: "beijing",
				CityName: "北京",
				Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
				ClientIP: "",
			},
			wantErr: "client IP is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := uc.Execute(ctx, tc.cmd)

			require.Error(t, err)
			assert.Nil(t, result)
			assert.True(t, apperrors.IsValidationError(err), "Error should be ValidationError")
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}

	// Verify no methods were called on mocks
	mockRepo.AssertNotCalled(t, "Save")
	mockCache.AssertNotCalled(t, "DeleteByPattern")
	mockRateLimiter.AssertNotCalled(t, "Allow")
}

// TestCreatePostUseCase_Execute_RateLimitExceeded tests rate limit exceeded.
func TestCreatePostUseCase_Execute_RateLimitExceeded(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)
	mockRateLimiter := new(MockRateLimiter)

	// Create use case
	uc := content.NewCreatePostUseCase(mockRepo, mockCache, mockRateLimiter)

	ctx := context.Background()
	cmd := content.CreatePostCommand{
		Company:  "测试公司",
		CityCode: "beijing",
		CityName: "北京",
		Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP: "127.0.0.1",
	}

	// Setup expectations - rate limit exceeded
	mockRateLimiter.On("Allow", ctx, mock.AnythingOfType("string"), 3, time.Hour).Return(false, nil)

	// Execute
	result, err := uc.Execute(ctx, cmd)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, apperrors.IsRateLimitError(err), "Error should be RateLimitError")

	// Verify repository and cache were not called
	mockRepo.AssertNotCalled(t, "Save")
	mockCache.AssertNotCalled(t, "DeleteByPattern")
	mockRateLimiter.AssertExpectations(t)
}

// TestCreatePostUseCase_Execute_RateLimitCheckError tests rate limit check error.
func TestCreatePostUseCase_Execute_RateLimitCheckError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)
	mockRateLimiter := new(MockRateLimiter)

	// Create use case
	uc := content.NewCreatePostUseCase(mockRepo, mockCache, mockRateLimiter)

	ctx := context.Background()
	cmd := content.CreatePostCommand{
		Company:  "测试公司",
		CityCode: "beijing",
		CityName: "北京",
		Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP: "127.0.0.1",
	}

	// Setup expectations - rate limit check fails
	mockRateLimiter.On("Allow", ctx, mock.AnythingOfType("string"), 3, time.Hour).
		Return(false, errors.New("redis connection failed"))

	// Execute
	result, err := uc.Execute(ctx, cmd)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, apperrors.IsDatabaseError(err), "Error should be DatabaseError")

	// Verify repository and cache were not called
	mockRepo.AssertNotCalled(t, "Save")
	mockCache.AssertNotCalled(t, "DeleteByPattern")
	mockRateLimiter.AssertExpectations(t)
}

// TestCreatePostUseCase_Execute_InvalidValueObjects tests invalid value objects.
func TestCreatePostUseCase_Execute_InvalidValueObjects(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)
	mockRateLimiter := new(MockRateLimiter)

	// Create use case
	uc := content.NewCreatePostUseCase(mockRepo, mockCache, mockRateLimiter)

	ctx := context.Background()

	testCases := []struct {
		name    string
		cmd     content.CreatePostCommand
		wantErr string
	}{
		{
			name: "invalid company name (too long)",
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
			name: "invalid content (too short)",
			cmd: content.CreatePostCommand{
				Company:  "测试公司",
				CityCode: "beijing",
				CityName: "北京",
				Content:  "太短", // Less than 10 characters
				ClientIP: "127.0.0.1",
			},
			wantErr: "invalid content",
		},
		{
			name: "invalid city (empty name after trim)",
			cmd: content.CreatePostCommand{
				Company:  "测试公司",
				CityCode: "beijing",
				CityName: "   ", // Only whitespace, will be trimmed to empty
				Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
				ClientIP: "127.0.0.1",
			},
			wantErr: "invalid city",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup expectations - rate limit passes
			mockRateLimiter.On("Allow", ctx, mock.AnythingOfType("string"), 3, time.Hour).Return(true, nil)

			result, err := uc.Execute(ctx, tc.cmd)

			require.Error(t, err)
			assert.Nil(t, result)
			assert.True(t, apperrors.IsValidationError(err), "Error should be ValidationError")
			assert.Contains(t, err.Error(), tc.wantErr)

			// Verify repository was not called
			mockRepo.AssertNotCalled(t, "Save")
			mockCache.AssertNotCalled(t, "DeleteByPattern")
		})
	}
}

// TestCreatePostUseCase_Execute_RepositoryError tests repository save error.
func TestCreatePostUseCase_Execute_RepositoryError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)
	mockRateLimiter := new(MockRateLimiter)

	// Create use case
	uc := content.NewCreatePostUseCase(mockRepo, mockCache, mockRateLimiter)

	ctx := context.Background()
	cmd := content.CreatePostCommand{
		Company:  "测试公司",
		CityCode: "beijing",
		CityName: "北京",
		Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP: "127.0.0.1",
	}

	// Setup expectations
	mockRateLimiter.On("Allow", ctx, mock.AnythingOfType("string"), 3, time.Hour).Return(true, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*content.Post")).
		Return(errors.New("database connection failed"))

	// Execute
	result, err := uc.Execute(ctx, cmd)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)

	// Verify cache was not called (because save failed)
	mockCache.AssertNotCalled(t, "DeleteByPattern")
	mockRepo.AssertExpectations(t)
	mockRateLimiter.AssertExpectations(t)
}

// TestCreatePostUseCase_Execute_CacheError tests cache deletion error (should not fail).
func TestCreatePostUseCase_Execute_CacheError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)
	mockRateLimiter := new(MockRateLimiter)

	// Create use case
	uc := content.NewCreatePostUseCase(mockRepo, mockCache, mockRateLimiter)

	ctx := context.Background()
	cmd := content.CreatePostCommand{
		Company:  "测试公司",
		CityCode: "beijing",
		CityName: "北京",
		Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP: "127.0.0.1",
	}

	// Setup expectations
	mockRateLimiter.On("Allow", ctx, mock.AnythingOfType("string"), 3, time.Hour).Return(true, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*content.Post")).Return(nil)
	// Cache deletion fails, but should not cause the operation to fail
	mockCache.On("DeleteByPattern", ctx, "posts:city:beijing:*").
		Return(errors.New("redis connection failed"))

	// Execute
	result, err := uc.Execute(ctx, cmd)

	// Assertions - should succeed despite cache error
	require.NoError(t, err, "Cache error should not fail the operation")
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ID)

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockRateLimiter.AssertExpectations(t)
}

// TestCreatePostUseCase_Execute_RateLimitKeyFormat tests rate limit key format.
// We test this indirectly by checking the key passed to rate limiter.
func TestCreatePostUseCase_Execute_RateLimitKeyFormat(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)
	mockRateLimiter := new(MockRateLimiter)

	// Create use case
	uc := content.NewCreatePostUseCase(mockRepo, mockCache, mockRateLimiter)

	ctx := context.Background()
	cmd := content.CreatePostCommand{
		Company:  "测试公司",
		CityCode: "beijing",
		CityName: "北京",
		Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		ClientIP: "127.0.0.1",
	}

	// Setup expectations - capture the rate limit key
	var capturedKey string
	mockRateLimiter.On("Allow", ctx, mock.AnythingOfType("string"), 3, time.Hour).
		Return(true, nil).
		Run(func(args mock.Arguments) {
			capturedKey = args.String(1) // key is the second argument
		})
	mockRepo.On("Save", ctx, mock.AnythingOfType("*content.Post")).Return(nil)
	mockCache.On("DeleteByPattern", ctx, "posts:city:beijing:*").Return(nil)

	// Execute
	_, err := uc.Execute(ctx, cmd)
	require.NoError(t, err)

	// Verify rate limit key format: "rate_limit:post:{ip}:{hour}"
	assert.Contains(t, capturedKey, "rate_limit:post:")
	assert.Contains(t, capturedKey, cmd.ClientIP)
	assert.Contains(t, capturedKey, ":")
}
