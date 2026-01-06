// Package content_test provides unit tests for content use cases.
// These tests use mocked dependencies to isolate the use case logic.
package content_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"fuck_boss/backend/internal/application/content"
	"fuck_boss/backend/internal/application/dto"
	domaincontent "fuck_boss/backend/internal/domain/content"
	"fuck_boss/backend/internal/domain/shared"
	apperrors "fuck_boss/backend/pkg/errors"
)

// TestGetPostUseCase_Execute_CacheHit tests cache hit scenario.
func TestGetPostUseCase_Execute_CacheHit(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := content.NewGetPostUseCase(mockRepo, mockCache)

	ctx := context.Background()
	postID := "550e8400-e29b-41d4-a716-446655440000"

	// Create cached data
	cachedResult := &dto.PostDTO{
		ID:        postID,
		Company:   "测试公司",
		CityCode:  "beijing",
		CityName:  "北京",
		Content:   "测试内容",
		CreatedAt: time.Now(),
	}
	cachedData, _ := json.Marshal(cachedResult)

	// Setup expectations
	mockCache.On("Get", ctx, "post:"+postID).Return(string(cachedData), nil)

	// Execute
	result, err := uc.Execute(ctx, postID)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, postID, result.ID)
	assert.Equal(t, "测试公司", result.Company)

	// Verify repository was not called
	mockRepo.AssertNotCalled(t, "FindByID")
	mockCache.AssertExpectations(t)
}

// TestGetPostUseCase_Execute_CacheMiss tests cache miss scenario.
func TestGetPostUseCase_Execute_CacheMiss(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := content.NewGetPostUseCase(mockRepo, mockCache)

	ctx := context.Background()
	postID := "550e8400-e29b-41d4-a716-446655440000"

	// Create test post
	company, _ := domaincontent.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条测试内容，用于验证获取功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	// Setup expectations
	mockCache.On("Get", ctx, "post:"+postID).Return("", errors.New("cache miss"))
	postIDVO, _ := domaincontent.NewPostID(postID)
	mockRepo.On("FindByID", ctx, postIDVO).Return(post, nil)
	mockCache.On("Set", ctx, "post:"+postID, mock.AnythingOfType("string"), 10*time.Minute).Return(nil)

	// Execute
	result, err := uc.Execute(ctx, postID)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, post.ID().String(), result.ID)
	assert.Equal(t, post.Company().String(), result.Company)

	// Verify all expectations
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestGetPostUseCase_Execute_ValidationError tests validation errors.
func TestGetPostUseCase_Execute_ValidationError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := content.NewGetPostUseCase(mockRepo, mockCache)

	ctx := context.Background()

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
		t.Run(tc.name, func(t *testing.T) {
			// Execute
			result, err := uc.Execute(ctx, tc.postID)

			// Assertions
			require.Error(t, err)
			assert.Nil(t, result)
			assert.True(t, apperrors.IsValidationError(err))
			assert.Contains(t, err.Error(), tc.wantErr)

			// Verify no methods were called
			mockRepo.AssertNotCalled(t, "FindByID")
			mockCache.AssertNotCalled(t, "Get")
		})
	}
}

// TestGetPostUseCase_Execute_NotFound tests not found error.
func TestGetPostUseCase_Execute_NotFound(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := content.NewGetPostUseCase(mockRepo, mockCache)

	ctx := context.Background()
	postID := "550e8400-e29b-41d4-a716-446655440000"

	// Setup expectations
	mockCache.On("Get", ctx, "post:"+postID).Return("", errors.New("cache miss"))
	postIDVO, _ := domaincontent.NewPostID(postID)
	mockRepo.On("FindByID", ctx, postIDVO).Return(nil, apperrors.NewNotFoundError("post"))

	// Execute
	result, err := uc.Execute(ctx, postID)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, apperrors.IsNotFoundError(err))

	// Verify cache was not updated
	mockCache.AssertNotCalled(t, "Set")
	mockRepo.AssertExpectations(t)
}

// TestGetPostUseCase_Execute_RepositoryError tests repository error.
func TestGetPostUseCase_Execute_RepositoryError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := content.NewGetPostUseCase(mockRepo, mockCache)

	ctx := context.Background()
	postID := "550e8400-e29b-41d4-a716-446655440000"

	// Setup expectations
	mockCache.On("Get", ctx, "post:"+postID).Return("", errors.New("cache miss"))
	postIDVO, _ := domaincontent.NewPostID(postID)
	mockRepo.On("FindByID", ctx, postIDVO).Return(nil, errors.New("database connection failed"))

	// Execute
	result, err := uc.Execute(ctx, postID)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, apperrors.IsDatabaseError(err))

	// Verify cache was not updated
	mockCache.AssertNotCalled(t, "Set")
	mockRepo.AssertExpectations(t)
}

// TestGetPostUseCase_Execute_CacheError tests cache error handling (should fallback to database).
func TestGetPostUseCase_Execute_CacheError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := content.NewGetPostUseCase(mockRepo, mockCache)

	ctx := context.Background()
	postID := "550e8400-e29b-41d4-a716-446655440000"

	// Create test post
	company, _ := domaincontent.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条测试内容，用于验证获取功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	// Setup expectations - cache error but should fallback to database
	mockCache.On("Get", ctx, "post:"+postID).Return("", errors.New("redis connection failed"))
	postIDVO, _ := domaincontent.NewPostID(postID)
	mockRepo.On("FindByID", ctx, postIDVO).Return(post, nil)
	mockCache.On("Set", ctx, "post:"+postID, mock.AnythingOfType("string"), 10*time.Minute).Return(nil)

	// Execute
	result, err := uc.Execute(ctx, postID)

	// Assertions - should succeed despite cache error
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, post.ID().String(), result.ID)

	// Verify all expectations
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestGetPostUseCase_Execute_CacheDeserializationError tests cache deserialization error (should fallback to database).
func TestGetPostUseCase_Execute_CacheDeserializationError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := content.NewGetPostUseCase(mockRepo, mockCache)

	ctx := context.Background()
	postID := "550e8400-e29b-41d4-a716-446655440000"

	// Create test post
	company, _ := domaincontent.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条测试内容，用于验证获取功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	// Setup expectations - invalid JSON in cache
	mockCache.On("Get", ctx, "post:"+postID).Return("invalid json", nil)
	postIDVO, _ := domaincontent.NewPostID(postID)
	mockRepo.On("FindByID", ctx, postIDVO).Return(post, nil)
	mockCache.On("Set", ctx, "post:"+postID, mock.AnythingOfType("string"), 10*time.Minute).Return(nil)

	// Execute
	result, err := uc.Execute(ctx, postID)

	// Assertions - should succeed despite deserialization error
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, post.ID().String(), result.ID)

	// Verify all expectations
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestGetPostUseCase_Execute_CacheSetError tests cache set error (should not fail).
func TestGetPostUseCase_Execute_CacheSetError(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockPostRepository)
	mockCache := new(MockCacheRepository)

	// Create use case
	uc := content.NewGetPostUseCase(mockRepo, mockCache)

	ctx := context.Background()
	postID := "550e8400-e29b-41d4-a716-446655440000"

	// Create test post
	company, _ := domaincontent.NewCompanyName("测试公司")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := domaincontent.NewContent("这是一条测试内容，用于验证获取功能。内容应该足够长以满足最小长度要求。")
	post, _ := domaincontent.NewPost(company, city, postContent)

	// Setup expectations - cache set fails but should not affect result
	mockCache.On("Get", ctx, "post:"+postID).Return("", errors.New("cache miss"))
	postIDVO, _ := domaincontent.NewPostID(postID)
	mockRepo.On("FindByID", ctx, postIDVO).Return(post, nil)
	mockCache.On("Set", ctx, "post:"+postID, mock.AnythingOfType("string"), 10*time.Minute).
		Return(errors.New("redis connection failed"))

	// Execute
	result, err := uc.Execute(ctx, postID)

	// Assertions - should succeed despite cache set error
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, post.ID().String(), result.ID)

	// Verify all expectations
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
