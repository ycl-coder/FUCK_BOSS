// Package grpc provides gRPC handlers for content management.
package grpc

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	contentv1 "fuck_boss/backend/api/proto/content/v1"
	"fuck_boss/backend/internal/application/content"
	"fuck_boss/backend/internal/application/dto"
	"fuck_boss/backend/internal/application/search"
	apperrors "fuck_boss/backend/pkg/errors"
)

// CreatePostUseCaseInterface defines the interface for creating posts.
type CreatePostUseCaseInterface interface {
	Execute(ctx context.Context, cmd content.CreatePostCommand) (*dto.PostDTO, error)
}

// ListPostsUseCaseInterface defines the interface for listing posts.
type ListPostsUseCaseInterface interface {
	Execute(ctx context.Context, query content.ListPostsQuery) (*dto.PostsListDTO, error)
}

// GetPostUseCaseInterface defines the interface for getting a post.
type GetPostUseCaseInterface interface {
	Execute(ctx context.Context, postID string) (*dto.PostDTO, error)
}

// SearchPostsUseCaseInterface defines the interface for searching posts.
type SearchPostsUseCaseInterface interface {
	Execute(ctx context.Context, query search.SearchPostsQuery) (*dto.PostsListDTO, error)
}

// ContentService implements the ContentService gRPC service.
type ContentService struct {
	contentv1.UnimplementedContentServiceServer

	// createUseCase handles post creation.
	createUseCase CreatePostUseCaseInterface

	// listUseCase handles post listing.
	listUseCase ListPostsUseCaseInterface

	// getUseCase handles post retrieval.
	getUseCase GetPostUseCaseInterface

	// searchUseCase handles post searching.
	searchUseCase SearchPostsUseCaseInterface
}

// NewContentService creates a new ContentService instance.
func NewContentService(
	createUseCase CreatePostUseCaseInterface,
	listUseCase ListPostsUseCaseInterface,
	getUseCase GetPostUseCaseInterface,
	searchUseCase SearchPostsUseCaseInterface,
) *ContentService {
	return &ContentService{
		createUseCase: createUseCase,
		listUseCase:   listUseCase,
		getUseCase:    getUseCase,
		searchUseCase: searchUseCase,
	}
}

// CreatePost handles the CreatePost gRPC request.
func (s *ContentService) CreatePost(ctx context.Context, req *contentv1.CreatePostRequest) (*contentv1.CreatePostResponse, error) {
	// Extract client IP from context
	clientIP := extractClientIP(ctx)

	// Convert occurred_at from Unix timestamp to time.Time
	var occurredAt *time.Time
	if req.OccurredAt > 0 {
		t := time.Unix(req.OccurredAt, 0)
		occurredAt = &t
	}

	// Create command
	cmd := content.CreatePostCommand{
		Company:    req.Company,
		CityCode:   req.CityCode,
		CityName:   req.CityName,
		Content:    req.Content,
		OccurredAt: occurredAt,
		ClientIP:   clientIP,
	}

	// Execute use case
	postDTO, err := s.createUseCase.Execute(ctx, cmd)
	if err != nil {
		return nil, convertError(err)
	}

	// Convert to response
	return &contentv1.CreatePostResponse{
		PostId:    postDTO.ID,
		CreatedAt: postDTO.CreatedAt.Unix(),
	}, nil
}

// ListPosts handles the ListPosts gRPC request.
func (s *ContentService) ListPosts(ctx context.Context, req *contentv1.ListPostsRequest) (*contentv1.ListPostsResponse, error) {
	// Create query
	query := content.ListPostsQuery{
		CityCode: req.CityCode,
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	// Execute use case
	result, err := s.listUseCase.Execute(ctx, query)
	if err != nil {
		return nil, convertError(err)
	}

	// Convert to response
	return &contentv1.ListPostsResponse{
		Posts:    convertPostsToProto(result.Posts),
		Total:    int32(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// GetPost handles the GetPost gRPC request.
func (s *ContentService) GetPost(ctx context.Context, req *contentv1.GetPostRequest) (*contentv1.GetPostResponse, error) {
	// Execute use case
	postDTO, err := s.getUseCase.Execute(ctx, req.PostId)
	if err != nil {
		return nil, convertError(err)
	}

	// Convert to response
	return &contentv1.GetPostResponse{
		Post: convertPostToProto(postDTO),
	}, nil
}

// SearchPosts handles the SearchPosts gRPC request.
func (s *ContentService) SearchPosts(ctx context.Context, req *contentv1.SearchPostsRequest) (*contentv1.SearchPostsResponse, error) {
	// Create query
	var cityCode *string
	if req.CityCode != "" {
		cityCode = &req.CityCode
	}

	query := search.SearchPostsQuery{
		Keyword:  req.Keyword,
		CityCode: cityCode,
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	// Execute use case
	result, err := s.searchUseCase.Execute(ctx, query)
	if err != nil {
		return nil, convertError(err)
	}

	// Convert to response
	return &contentv1.SearchPostsResponse{
		Posts:    convertPostsToProto(result.Posts),
		Total:    int32(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// extractClientIP extracts the client IP address from the gRPC context.
// It tries to get the IP from peer information first, then from metadata.
func extractClientIP(ctx context.Context) string {
	// Try to get IP from peer
	if p, ok := peer.FromContext(ctx); ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			return addr.IP.String()
		}
		if addr, ok := p.Addr.(*net.UDPAddr); ok {
			return addr.IP.String()
		}
		// Fallback to string representation
		return p.Addr.String()
	}

	// Try to get IP from metadata (X-Forwarded-For or X-Real-IP)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if forwardedFor := md.Get("x-forwarded-for"); len(forwardedFor) > 0 {
			return forwardedFor[0]
		}
		if realIP := md.Get("x-real-ip"); len(realIP) > 0 {
			return realIP[0]
		}
	}

	// Default fallback
	return "127.0.0.1"
}

// convertError converts application errors to gRPC status errors.
func convertError(err error) error {
	if err == nil {
		return nil
	}

	// Check error type
	switch {
	case apperrors.IsValidationError(err):
		return status.Error(codes.InvalidArgument, err.Error())
	case apperrors.IsNotFoundError(err):
		return status.Error(codes.NotFound, err.Error())
	case apperrors.IsRateLimitError(err):
		return status.Error(codes.ResourceExhausted, err.Error())
	case apperrors.IsDatabaseError(err):
		return status.Error(codes.Internal, "internal server error")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

// convertPostToProto converts a PostDTO to a protobuf Post message.
func convertPostToProto(postDTO *dto.PostDTO) *contentv1.Post {
	if postDTO == nil {
		return nil
	}

	var occurredAt int64
	if postDTO.OccurredAt != nil {
		occurredAt = postDTO.OccurredAt.Unix()
	}

	return &contentv1.Post{
		Id:         postDTO.ID,
		Company:    postDTO.Company,
		CityCode:   postDTO.CityCode,
		CityName:   postDTO.CityName,
		Content:    postDTO.Content,
		OccurredAt: occurredAt,
		CreatedAt:  postDTO.CreatedAt.Unix(),
	}
}

// convertPostsToProto converts a slice of PostDTOs to protobuf Post messages.
func convertPostsToProto(posts []*dto.PostDTO) []*contentv1.Post {
	if posts == nil {
		return nil
	}

	result := make([]*contentv1.Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, convertPostToProto(post))
	}
	return result
}
