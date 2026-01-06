// Package grpc_test provides unit tests for gRPC handlers.
// These tests use mocked UseCases to isolate the handler logic.
package grpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	contentv1 "fuck_boss/backend/api/proto/content/v1"
	"fuck_boss/backend/internal/application/content"
	"fuck_boss/backend/internal/application/dto"
	"fuck_boss/backend/internal/application/search"
	apperrors "fuck_boss/backend/pkg/errors"
	grpchandler "fuck_boss/backend/internal/presentation/grpc"
)

// MockCreatePostUseCase is a mock implementation of CreatePostUseCase.
type MockCreatePostUseCase struct {
	mock.Mock
}

func (m *MockCreatePostUseCase) Execute(ctx context.Context, cmd content.CreatePostCommand) (*dto.PostDTO, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PostDTO), args.Error(1)
}

// MockListPostsUseCase is a mock implementation of ListPostsUseCase.
type MockListPostsUseCase struct {
	mock.Mock
}

func (m *MockListPostsUseCase) Execute(ctx context.Context, query content.ListPostsQuery) (*dto.PostsListDTO, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PostsListDTO), args.Error(1)
}

// MockGetPostUseCase is a mock implementation of GetPostUseCase.
type MockGetPostUseCase struct {
	mock.Mock
}

func (m *MockGetPostUseCase) Execute(ctx context.Context, postID string) (*dto.PostDTO, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PostDTO), args.Error(1)
}

// MockSearchPostsUseCase is a mock implementation of SearchPostsUseCase.
type MockSearchPostsUseCase struct {
	mock.Mock
}

func (m *MockSearchPostsUseCase) Execute(ctx context.Context, query search.SearchPostsQuery) (*dto.PostsListDTO, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PostsListDTO), args.Error(1)
}

// TestContentService_CreatePost_Success tests successful post creation.
func TestContentService_CreatePost_Success(t *testing.T) {
	// Setup mocks
	mockCreate := new(MockCreatePostUseCase)
	mockList := new(MockListPostsUseCase)
	mockGet := new(MockGetPostUseCase)
	mockSearch := new(MockSearchPostsUseCase)

	// Create service
	service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

	// Create context with peer info (for client IP extraction)
	ctx := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("192.168.1.100"),
			Port: 12345,
		},
	})

	// Create request
	req := &contentv1.CreatePostRequest{
		Company:    "测试公司",
		CityCode:   "beijing",
		CityName:   "北京",
		Content:    "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		OccurredAt: 0, // Optional
	}

	// Create expected DTO
	expectedDTO := &dto.PostDTO{
		ID:        "test-post-id",
		Company:   "测试公司",
		CityCode:  "beijing",
		CityName:  "北京",
		Content:   "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		CreatedAt: time.Now(),
	}

	// Setup expectations
	mockCreate.On("Execute", ctx, mock.MatchedBy(func(cmd content.CreatePostCommand) bool {
		return cmd.Company == req.Company &&
			cmd.CityCode == req.CityCode &&
			cmd.CityName == req.CityName &&
			cmd.Content == req.Content &&
			cmd.ClientIP == "192.168.1.100"
	})).Return(expectedDTO, nil)

	// Execute
	resp, err := service.CreatePost(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, expectedDTO.ID, resp.PostId)
	assert.Equal(t, expectedDTO.CreatedAt.Unix(), resp.CreatedAt)

	// Verify mock was called
	mockCreate.AssertExpectations(t)
}

// TestContentService_CreatePost_WithOccurredAt tests post creation with occurred_at.
func TestContentService_CreatePost_WithOccurredAt(t *testing.T) {
	// Setup mocks
	mockCreate := new(MockCreatePostUseCase)
	mockList := new(MockListPostsUseCase)
	mockGet := new(MockGetPostUseCase)
	mockSearch := new(MockSearchPostsUseCase)

	// Create service
	service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

	// Create context
	ctx := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("192.168.1.100"),
			Port: 12345,
		},
	})

	// Create request with occurred_at
	occurredAt := time.Now().Add(-24 * time.Hour)
	req := &contentv1.CreatePostRequest{
		Company:    "测试公司",
		CityCode:   "beijing",
		CityName:   "北京",
		Content:    "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		OccurredAt: occurredAt.Unix(),
	}

	// Create expected DTO
	expectedDTO := &dto.PostDTO{
		ID:        "test-post-id",
		Company:   "测试公司",
		CityCode:  "beijing",
		CityName:  "北京",
		Content:   "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		CreatedAt: time.Now(),
	}

	// Setup expectations
	mockCreate.On("Execute", ctx, mock.MatchedBy(func(cmd content.CreatePostCommand) bool {
		return cmd.Company == req.Company &&
			cmd.CityCode == req.CityCode &&
			cmd.CityName == req.CityName &&
			cmd.Content == req.Content &&
			cmd.OccurredAt != nil &&
			cmd.OccurredAt.Unix() == occurredAt.Unix()
	})).Return(expectedDTO, nil)

	// Execute
	resp, err := service.CreatePost(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, expectedDTO.ID, resp.PostId)

	// Verify mock was called
	mockCreate.AssertExpectations(t)
}

// TestContentService_CreatePost_ValidationError tests validation error handling.
func TestContentService_CreatePost_ValidationError(t *testing.T) {
	// Setup mocks
	mockCreate := new(MockCreatePostUseCase)
	mockList := new(MockListPostsUseCase)
	mockGet := new(MockGetPostUseCase)
	mockSearch := new(MockSearchPostsUseCase)

	// Create service
	service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

	// Create context
	ctx := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("192.168.1.100"),
			Port: 12345,
		},
	})

	// Create request
	req := &contentv1.CreatePostRequest{
		Company:    "", // Empty company name
		CityCode:   "beijing",
		CityName:   "北京",
		Content:    "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		OccurredAt: 0,
	}

	// Setup expectations
	mockCreate.On("Execute", ctx, mock.Anything).Return(nil, apperrors.NewValidationError("company name is required"))

	// Execute
	resp, err := service.CreatePost(ctx, req)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, resp)

	// Verify it's a gRPC InvalidArgument error
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "company name is required")

	// Verify mock was called
	mockCreate.AssertExpectations(t)
}

// TestContentService_CreatePost_RateLimitError tests rate limit error handling.
func TestContentService_CreatePost_RateLimitError(t *testing.T) {
	// Setup mocks
	mockCreate := new(MockCreatePostUseCase)
	mockList := new(MockListPostsUseCase)
	mockGet := new(MockGetPostUseCase)
	mockSearch := new(MockSearchPostsUseCase)

	// Create service
	service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

	// Create context
	ctx := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("192.168.1.100"),
			Port: 12345,
		},
	})

	// Create request
	req := &contentv1.CreatePostRequest{
		Company:    "测试公司",
		CityCode:   "beijing",
		CityName:   "北京",
		Content:    "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
		OccurredAt: 0,
	}

	// Setup expectations
	mockCreate.On("Execute", ctx, mock.Anything).Return(nil, apperrors.NewRateLimitError("rate limit exceeded"))

	// Execute
	resp, err := service.CreatePost(ctx, req)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, resp)

	// Verify it's a gRPC ResourceExhausted error
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.ResourceExhausted, st.Code())
	assert.Contains(t, st.Message(), "rate limit exceeded")

	// Verify mock was called
	mockCreate.AssertExpectations(t)
}

// TestContentService_ListPosts_Success tests successful post listing.
func TestContentService_ListPosts_Success(t *testing.T) {
	// Setup mocks
	mockCreate := new(MockCreatePostUseCase)
	mockList := new(MockListPostsUseCase)
	mockGet := new(MockGetPostUseCase)
	mockSearch := new(MockSearchPostsUseCase)

	// Create service
	service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

	// Create context
	ctx := context.Background()

	// Create request
	req := &contentv1.ListPostsRequest{
		CityCode: "beijing",
		Page:     1,
		PageSize: 20,
	}

	// Create expected DTO
	expectedDTO := &dto.PostsListDTO{
		Posts: []*dto.PostDTO{
			{
				ID:        "post-1",
				Company:   "公司A",
				CityCode:  "beijing",
				CityName:  "北京",
				Content:   "内容A",
				CreatedAt: time.Now(),
			},
		},
		Total:    1,
		Page:     1,
		PageSize: 20,
	}

	// Setup expectations
	mockList.On("Execute", ctx, content.ListPostsQuery{
		CityCode: req.CityCode,
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}).Return(expectedDTO, nil)

	// Execute
	resp, err := service.ListPosts(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.Total)
	assert.Equal(t, int32(1), resp.Page)
	assert.Equal(t, int32(20), resp.PageSize)
	assert.Len(t, resp.Posts, 1)
	assert.Equal(t, "post-1", resp.Posts[0].Id)

	// Verify mock was called
	mockList.AssertExpectations(t)
}

// TestContentService_GetPost_Success tests successful post retrieval.
func TestContentService_GetPost_Success(t *testing.T) {
	// Setup mocks
	mockCreate := new(MockCreatePostUseCase)
	mockList := new(MockListPostsUseCase)
	mockGet := new(MockGetPostUseCase)
	mockSearch := new(MockSearchPostsUseCase)

	// Create service
	service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

	// Create context
	ctx := context.Background()

	// Create request
	req := &contentv1.GetPostRequest{
		PostId: "test-post-id",
	}

	// Create expected DTO
	expectedDTO := &dto.PostDTO{
		ID:        "test-post-id",
		Company:   "测试公司",
		CityCode:  "beijing",
		CityName:  "北京",
		Content:   "测试内容",
		CreatedAt: time.Now(),
	}

	// Setup expectations
	mockGet.On("Execute", ctx, req.PostId).Return(expectedDTO, nil)

	// Execute
	resp, err := service.GetPost(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Post)
	assert.Equal(t, expectedDTO.ID, resp.Post.Id)
	assert.Equal(t, expectedDTO.Company, resp.Post.Company)
	assert.Equal(t, expectedDTO.CityCode, resp.Post.CityCode)
	assert.Equal(t, expectedDTO.CityName, resp.Post.CityName)
	assert.Equal(t, expectedDTO.Content, resp.Post.Content)

	// Verify mock was called
	mockGet.AssertExpectations(t)
}

// TestContentService_GetPost_NotFound tests not found error handling.
func TestContentService_GetPost_NotFound(t *testing.T) {
	// Setup mocks
	mockCreate := new(MockCreatePostUseCase)
	mockList := new(MockListPostsUseCase)
	mockGet := new(MockGetPostUseCase)
	mockSearch := new(MockSearchPostsUseCase)

	// Create service
	service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

	// Create context
	ctx := context.Background()

	// Create request
	req := &contentv1.GetPostRequest{
		PostId: "non-existent-id",
	}

	// Setup expectations
	mockGet.On("Execute", ctx, req.PostId).Return(nil, apperrors.NewNotFoundError("post not found"))

	// Execute
	resp, err := service.GetPost(ctx, req)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, resp)

	// Verify it's a gRPC NotFound error
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "post not found")

	// Verify mock was called
	mockGet.AssertExpectations(t)
}

// TestContentService_SearchPosts_Success tests successful post searching.
func TestContentService_SearchPosts_Success(t *testing.T) {
	// Setup mocks
	mockCreate := new(MockCreatePostUseCase)
	mockList := new(MockListPostsUseCase)
	mockGet := new(MockGetPostUseCase)
	mockSearch := new(MockSearchPostsUseCase)

	// Create service
	service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

	// Create context
	ctx := context.Background()

	// Create request
	req := &contentv1.SearchPostsRequest{
		Keyword:  "测试",
		CityCode: "beijing",
		Page:     1,
		PageSize: 20,
	}

	// Create expected DTO
	expectedDTO := &dto.PostsListDTO{
		Posts: []*dto.PostDTO{
			{
				ID:        "post-1",
				Company:   "测试公司",
				CityCode:  "beijing",
				CityName:  "北京",
				Content:   "测试内容",
				CreatedAt: time.Now(),
			},
		},
		Total:    1,
		Page:     1,
		PageSize: 20,
	}

	// Setup expectations
	mockSearch.On("Execute", ctx, search.SearchPostsQuery{
		Keyword:  req.Keyword,
		CityCode: &req.CityCode,
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}).Return(expectedDTO, nil)

	// Execute
	resp, err := service.SearchPosts(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.Total)
	assert.Len(t, resp.Posts, 1)

	// Verify mock was called
	mockSearch.AssertExpectations(t)
}

// TestContentService_SearchPosts_WithoutCity tests search without city filter.
func TestContentService_SearchPosts_WithoutCity(t *testing.T) {
	// Setup mocks
	mockCreate := new(MockCreatePostUseCase)
	mockList := new(MockListPostsUseCase)
	mockGet := new(MockGetPostUseCase)
	mockSearch := new(MockSearchPostsUseCase)

	// Create service
	service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

	// Create context
	ctx := context.Background()

	// Create request without city code
	req := &contentv1.SearchPostsRequest{
		Keyword:  "测试",
		CityCode: "", // Empty city code
		Page:     1,
		PageSize: 20,
	}

	// Create expected DTO
	expectedDTO := &dto.PostsListDTO{
		Posts:    []*dto.PostDTO{},
		Total:    0,
		Page:     1,
		PageSize: 20,
	}

	// Setup expectations - cityCode should be nil when empty
	mockSearch.On("Execute", ctx, search.SearchPostsQuery{
		Keyword:  req.Keyword,
		CityCode: nil,
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}).Return(expectedDTO, nil)

	// Execute
	resp, err := service.SearchPosts(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int32(0), resp.Total)

	// Verify mock was called
	mockSearch.AssertExpectations(t)
}

// TestContentService_ErrorConversion tests error conversion to gRPC status codes.
func TestContentService_ErrorConversion(t *testing.T) {
	testCases := []struct {
		name           string
		appError       error
		expectedCode   codes.Code
		expectedInMsg  string
		handler        func(*grpchandler.ContentService, context.Context) (interface{}, error)
	}{
		{
			name:          "validation error",
			appError:      apperrors.NewValidationError("validation failed"),
			expectedCode:  codes.InvalidArgument,
			expectedInMsg: "validation failed",
			handler: func(s *grpchandler.ContentService, ctx context.Context) (interface{}, error) {
				mockCreate := new(MockCreatePostUseCase)
				// Create context with peer for CreatePost
				createCtx := peer.NewContext(ctx, &peer.Peer{
					Addr: &net.TCPAddr{
						IP:   net.ParseIP("192.168.1.100"),
						Port: 12345,
					},
				})
				mockCreate.On("Execute", createCtx, mock.Anything).Return(nil, apperrors.NewValidationError("validation failed"))
				s = grpchandler.NewContentService(mockCreate, nil, nil, nil)
				return s.CreatePost(createCtx, &contentv1.CreatePostRequest{
					Company:  "test",
					CityCode: "beijing",
					CityName: "北京",
					Content:  "测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
				})
			},
		},
		{
			name:          "not found error",
			appError:      apperrors.NewNotFoundError("not found"),
			expectedCode:  codes.NotFound,
			expectedInMsg: "not found",
			handler: func(s *grpchandler.ContentService, ctx context.Context) (interface{}, error) {
				mockGet := new(MockGetPostUseCase)
				mockGet.On("Execute", ctx, "test-id").Return(nil, apperrors.NewNotFoundError("not found"))
				s = grpchandler.NewContentService(nil, nil, mockGet, nil)
				return s.GetPost(ctx, &contentv1.GetPostRequest{PostId: "test-id"})
			},
		},
		{
			name:          "rate limit error",
			appError:      apperrors.NewRateLimitError("rate limit exceeded"),
			expectedCode:  codes.ResourceExhausted,
			expectedInMsg: "rate limit exceeded",
			handler: func(s *grpchandler.ContentService, ctx context.Context) (interface{}, error) {
				mockCreate := new(MockCreatePostUseCase)
				// Create context with peer for CreatePost
				createCtx := peer.NewContext(ctx, &peer.Peer{
					Addr: &net.TCPAddr{
						IP:   net.ParseIP("192.168.1.100"),
						Port: 12345,
					},
				})
				mockCreate.On("Execute", createCtx, mock.Anything).Return(nil, apperrors.NewRateLimitError("rate limit exceeded"))
				s = grpchandler.NewContentService(mockCreate, nil, nil, nil)
				return s.CreatePost(createCtx, &contentv1.CreatePostRequest{
					Company:  "test",
					CityCode: "beijing",
					CityName: "北京",
					Content:  "测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
				})
			},
		},
		{
			name:          "database error",
			appError:      apperrors.NewDatabaseError("database error"),
			expectedCode:  codes.Internal,
			expectedInMsg: "internal server error",
			handler: func(s *grpchandler.ContentService, ctx context.Context) (interface{}, error) {
				mockGet := new(MockGetPostUseCase)
				mockGet.On("Execute", ctx, "test-id").Return(nil, apperrors.NewDatabaseError("database error"))
				s = grpchandler.NewContentService(nil, nil, mockGet, nil)
				return s.GetPost(ctx, &contentv1.GetPostRequest{PostId: "test-id"})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			service := grpchandler.NewContentService(nil, nil, nil, nil)

			_, err := tc.handler(service, ctx)

			require.Error(t, err)
			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tc.expectedCode, st.Code())
			assert.Contains(t, st.Message(), tc.expectedInMsg)
		})
	}
}

// TestContentService_ExtractClientIP tests client IP extraction from context.
func TestContentService_ExtractClientIP(t *testing.T) {
	testCases := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name: "from peer TCPAddr",
			ctx: peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.TCPAddr{
					IP:   net.ParseIP("192.168.1.100"),
					Port: 12345,
				},
			}),
			expected: "192.168.1.100",
		},
		{
			name: "from metadata X-Forwarded-For",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"x-forwarded-for": "10.0.0.1",
			})),
			expected: "10.0.0.1",
		},
		{
			name: "from metadata X-Real-IP",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"x-real-ip": "10.0.0.2",
			})),
			expected: "10.0.0.2",
		},
		{
			name:     "default fallback",
			ctx:      context.Background(),
			expected: "127.0.0.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request that will trigger client IP extraction
			mockCreate := new(MockCreatePostUseCase)
			mockList := new(MockListPostsUseCase)
			mockGet := new(MockGetPostUseCase)
			mockSearch := new(MockSearchPostsUseCase)

			service := grpchandler.NewContentService(mockCreate, mockList, mockGet, mockSearch)

			req := &contentv1.CreatePostRequest{
				Company:  "测试公司",
				CityCode: "beijing",
				CityName: "北京",
				Content:  "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
			}

			// Setup expectation to capture the ClientIP
			var capturedIP string
			mockCreate.On("Execute", tc.ctx, mock.MatchedBy(func(cmd content.CreatePostCommand) bool {
				capturedIP = cmd.ClientIP
				return true
			})).Return(&dto.PostDTO{
				ID:        "test-id",
				Company:   "测试公司",
				CityCode:  "beijing",
				CityName:  "北京",
				Content:   "测试内容",
				CreatedAt: time.Now(),
			}, nil)

			_, err := service.CreatePost(tc.ctx, req)
			require.NoError(t, err)

			// Verify client IP was extracted correctly
			assert.Equal(t, tc.expected, capturedIP)
		})
	}
}

