// Package rest provides REST API handlers that convert JSON requests to gRPC calls.
package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fuck_boss/backend/internal/application/content"
	"fuck_boss/backend/internal/application/dto"
	"fuck_boss/backend/internal/application/search"
	apperrors "fuck_boss/backend/pkg/errors"
)

// ContentHandler handles REST API requests for content operations.
type ContentHandler struct {
	createUseCase CreatePostUseCaseInterface
	listUseCase   ListPostsUseCaseInterface
	getUseCase    GetPostUseCaseInterface
	searchUseCase SearchPostsUseCaseInterface
	logger        Logger
}

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

// Logger interface for logging.
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}

// NewContentHandler creates a new ContentHandler.
func NewContentHandler(
	createUseCase CreatePostUseCaseInterface,
	listUseCase ListPostsUseCaseInterface,
	getUseCase GetPostUseCaseInterface,
	searchUseCase SearchPostsUseCaseInterface,
	logger Logger,
) *ContentHandler {
	return &ContentHandler{
		createUseCase: createUseCase,
		listUseCase:   listUseCase,
		getUseCase:    getUseCase,
		searchUseCase: searchUseCase,
		logger:        logger,
	}
}

// CreatePostRequest is the JSON request for creating a post.
type CreatePostRequest struct {
	Company    string `json:"company"`
	CityCode   string `json:"cityCode"`
	CityName   string `json:"cityName"`
	Content    string `json:"content"`
	OccurredAt *int64 `json:"occurredAt,omitempty"`
}

// CreatePostResponse is the JSON response for creating a post.
type CreatePostResponse struct {
	PostID    string `json:"postId"`
	CreatedAt int64  `json:"createdAt"`
}

// ListPostsRequest is the JSON request for listing posts.
type ListPostsRequest struct {
	CityCode string `json:"cityCode"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

// PostResponse is the JSON response for a post.
type PostResponse struct {
	ID         string  `json:"id"`
	Company    string  `json:"company"`
	CityCode   string  `json:"cityCode"`
	CityName   string  `json:"cityName"`
	Content    string  `json:"content"`
	OccurredAt *int64  `json:"occurredAt,omitempty"`
	CreatedAt  int64   `json:"createdAt"`
}

// ListPostsResponse is the JSON response for listing posts.
type ListPostsResponse struct {
	Posts    []*PostResponse `json:"posts"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}

// SearchPostsRequest is the JSON request for searching posts.
type SearchPostsRequest struct {
	Keyword  string  `json:"keyword"`
	CityCode *string `json:"cityCode,omitempty"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}

// CreatePost handles POST /api/posts
func (h *ContentHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Extract client IP from request
	clientIP := extractClientIP(r)

	// Convert to use case command
	var occurredAt *time.Time
	if req.OccurredAt != nil && *req.OccurredAt > 0 {
		t := time.Unix(*req.OccurredAt, 0)
		occurredAt = &t
	}
	cmd := content.CreatePostCommand{
		Company:    req.Company,
		CityCode:   req.CityCode,
		CityName:   req.CityName,
		Content:    req.Content,
		OccurredAt: occurredAt,
		ClientIP:   clientIP,
	}

	// Execute use case
	ctx := r.Context()
	dto, err := h.createUseCase.Execute(ctx, cmd)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Convert to response
	resp := CreatePostResponse{
		PostID:    dto.ID,
		CreatedAt: dto.CreatedAt.Unix(),
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// ListPosts handles GET /api/posts
func (h *ContentHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse query parameters
	cityCode := r.URL.Query().Get("cityCode")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}

	// Convert to use case query
	query := content.ListPostsQuery{
		CityCode: cityCode,
		Page:     page,
		PageSize: pageSize,
	}

	// Execute use case
	ctx := r.Context()
	dto, err := h.listUseCase.Execute(ctx, query)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Convert to response
	resp := ListPostsResponse{
		Posts:    convertPostsToResponse(dto.Posts),
		Total:    dto.Total,
		Page:     dto.Page,
		PageSize: dto.PageSize,
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// GetPost handles GET /api/posts/:id
func (h *ContentHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract post ID from URL path
	postID := r.URL.Path[len("/api/posts/"):]
	if postID == "" {
		h.writeError(w, http.StatusBadRequest, "Post ID is required")
		return
	}

	// Execute use case
	ctx := r.Context()
	dto, err := h.getUseCase.Execute(ctx, postID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Convert to response
	resp := convertPostToResponse(dto)
	h.writeJSON(w, http.StatusOK, resp)
}

// SearchPosts handles GET /api/posts/search
func (h *ContentHandler) SearchPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req SearchPostsRequest

	if r.Method == http.MethodPost {
		// Parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}
	} else {
		// Parse query parameters
		req.Keyword = r.URL.Query().Get("keyword")
		if cityCode := r.URL.Query().Get("cityCode"); cityCode != "" {
			req.CityCode = &cityCode
		}
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		req.Page = page
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if pageSize < 1 {
			pageSize = 20
		}
		req.PageSize = pageSize
	}

	// Convert to use case query
	query := search.SearchPostsQuery{
		Keyword:  req.Keyword,
		CityCode: nil,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	if req.CityCode != nil && *req.CityCode != "" {
		cityCode := *req.CityCode
		query.CityCode = &cityCode
	}

	// Execute use case
	ctx := r.Context()
	dto, err := h.searchUseCase.Execute(ctx, query)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Convert to response
	resp := ListPostsResponse{
		Posts:    convertPostsToResponse(dto.Posts),
		Total:    dto.Total,
		Page:     dto.Page,
		PageSize: dto.PageSize,
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// convertPostToResponse converts a DTO to a JSON response.
func convertPostToResponse(dto *dto.PostDTO) *PostResponse {
	resp := &PostResponse{
		ID:        dto.ID,
		Company:   dto.Company,
		CityCode:  dto.CityCode,
		CityName:  dto.CityName,
		Content:   dto.Content,
		CreatedAt: dto.CreatedAt.Unix(),
	}
	if dto.OccurredAt != nil {
		ts := dto.OccurredAt.Unix()
		resp.OccurredAt = &ts
	}
	return resp
}

// convertPostsToResponse converts a slice of DTOs to JSON responses.
func convertPostsToResponse(dtos []*dto.PostDTO) []*PostResponse {
	posts := make([]*PostResponse, len(dtos))
	for i, dto := range dtos {
		posts[i] = convertPostToResponse(dto)
	}
	return posts
}

// handleError converts application errors to HTTP responses.
func (h *ContentHandler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// Check for application errors
	var appErr *apperrors.AppError
	if apperrors.As(err, &appErr) {
		switch appErr.Code {
		case apperrors.ErrCodeValidation:
			h.writeError(w, http.StatusBadRequest, appErr.Message)
		case apperrors.ErrCodeNotFound:
			h.writeError(w, http.StatusNotFound, appErr.Message)
		case apperrors.ErrCodeRateLimit:
			h.writeError(w, http.StatusTooManyRequests, appErr.Message)
		default:
			h.logger.Error("Internal error", zap.Error(err))
			h.writeError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	// Check for gRPC status errors
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			h.writeError(w, http.StatusBadRequest, st.Message())
		case codes.NotFound:
			h.writeError(w, http.StatusNotFound, st.Message())
		case codes.ResourceExhausted:
			h.writeError(w, http.StatusTooManyRequests, st.Message())
		default:
			h.logger.Error("gRPC error", zap.Error(err))
			h.writeError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	// Unknown error
	h.logger.Error("Unknown error", zap.Error(err))
	h.writeError(w, http.StatusInternalServerError, "Internal server error")
}

// writeJSON writes a JSON response.
func (h *ContentHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}

// writeError writes an error response.
func (h *ContentHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// extractClientIP extracts the client IP address from the request.
func extractClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		return forwardedFor
	}
	// Check X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	// Fallback to RemoteAddr
	return r.RemoteAddr
}

