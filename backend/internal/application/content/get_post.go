// Package content provides use cases for content management.
package content

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"fuck_boss/backend/internal/application/cache"
	"fuck_boss/backend/internal/application/dto"
	"fuck_boss/backend/internal/domain/content"
	apperrors "fuck_boss/backend/pkg/errors"
)

// GetPostUseCase handles getting a single post by ID with caching.
// It coordinates domain repository and cache repository.
type GetPostUseCase struct {
	// repo is the Post repository.
	repo content.PostRepository

	// cacheRepo is the cache repository for caching query results.
	cacheRepo cache.CacheRepository
}

// NewGetPostUseCase creates a new GetPostUseCase instance.
func NewGetPostUseCase(
	repo content.PostRepository,
	cacheRepo cache.CacheRepository,
) *GetPostUseCase {
	return &GetPostUseCase{
		repo:      repo,
		cacheRepo: cacheRepo,
	}
}

// Execute executes the get post query by ID.
// It checks cache first, then queries the repository if cache misses.
func (uc *GetPostUseCase) Execute(ctx context.Context, postID string) (*dto.PostDTO, error) {
	// Validate post ID
	if postID == "" {
		return nil, apperrors.NewValidationError("post ID is required")
	}

	// Parse post ID
	postIDVO, err := content.NewPostID(postID)
	if err != nil {
		return nil, apperrors.NewValidationErrorWithDetails("invalid post ID", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Build cache key
	cacheKey := uc.buildCacheKey(postID)

	// Try to get from cache
	cachedData, err := uc.cacheRepo.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		// Cache hit: deserialize and return
		var result dto.PostDTO
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			return &result, nil
		}
		// If deserialization fails, continue to query database
	}

	// Cache miss or error: query repository
	post, err := uc.repo.FindByID(ctx, postIDVO)
	if err != nil {
		// Check if it's a not found error
		if apperrors.IsNotFoundError(err) {
			return nil, err
		}
		return nil, apperrors.NewDatabaseErrorWithCause("failed to query post", err)
	}

	// Convert to DTO
	result := uc.toDTO(post)

	// Update cache (non-blocking, errors are ignored)
	uc.updateCache(ctx, cacheKey, result)

	return result, nil
}

// buildCacheKey builds the cache key for the given post ID.
// Format: "post:{postID}"
func (uc *GetPostUseCase) buildCacheKey(postID string) string {
	return fmt.Sprintf("post:%s", postID)
}

// getCacheTTL returns the cache TTL for post details.
// Post details are cached for 10 minutes.
func (uc *GetPostUseCase) getCacheTTL() time.Duration {
	return 10 * time.Minute
}

// updateCache updates the cache with the query result.
// Errors are ignored to ensure cache failures don't affect the main flow.
func (uc *GetPostUseCase) updateCache(ctx context.Context, key string, result *dto.PostDTO) {
	// Serialize to JSON
	data, err := json.Marshal(result)
	if err != nil {
		// Log error but don't fail
		return
	}

	// Get TTL
	ttl := uc.getCacheTTL()

	// Set cache (non-blocking)
	_ = uc.cacheRepo.Set(ctx, key, string(data), ttl)
}

// toDTO converts a Post entity to PostDTO.
func (uc *GetPostUseCase) toDTO(post *content.Post) *dto.PostDTO {
	return &dto.PostDTO{
		ID:         post.ID().String(),
		Company:    post.Company().String(),
		CityCode:   post.City().Code(),
		CityName:   post.City().Name(),
		Content:    post.Content().String(),
		OccurredAt: nil, // Not stored in Post entity
		CreatedAt:  post.CreatedAt(),
	}
}
