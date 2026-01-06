// Package content provides use cases for content management.
package content

import (
	"context"
	"fmt"
	"time"

	"fuck_boss/backend/internal/application/cache"
	"fuck_boss/backend/internal/application/dto"
	"fuck_boss/backend/internal/application/ratelimit"
	"fuck_boss/backend/internal/domain/content"
	"fuck_boss/backend/internal/domain/shared"
	apperrors "fuck_boss/backend/pkg/errors"
)

// CreatePostCommand represents the command to create a new post.
type CreatePostCommand struct {
	// Company is the company name (required, 1-100 characters).
	Company string

	// CityCode is the city code (required, e.g., "beijing").
	CityCode string

	// CityName is the city name (required, e.g., "北京").
	// This is needed to create the City value object.
	CityName string

	// Content is the post content (required, 10-5000 characters).
	Content string

	// OccurredAt is when the incident occurred (optional).
	OccurredAt *time.Time

	// ClientIP is the client IP address for rate limiting (required).
	ClientIP string
}

// CreatePostUseCase handles the creation of posts.
// It coordinates domain entities, repositories, caching, and rate limiting.
type CreatePostUseCase struct {
	// repo is the Post repository.
	repo content.PostRepository

	// cacheRepo is the cache repository for cache invalidation.
	cacheRepo cache.CacheRepository

	// rateLimiter is the rate limiter for preventing abuse.
	rateLimiter ratelimit.RateLimiter
}

// NewCreatePostUseCase creates a new CreatePostUseCase instance.
func NewCreatePostUseCase(
	repo content.PostRepository,
	cacheRepo cache.CacheRepository,
	rateLimiter ratelimit.RateLimiter,
) *CreatePostUseCase {
	return &CreatePostUseCase{
		repo:        repo,
		cacheRepo:   cacheRepo,
		rateLimiter: rateLimiter,
	}
}

// Execute executes the create post command.
// It performs validation, rate limiting, creates the post, saves it, and clears cache.
func (uc *CreatePostUseCase) Execute(ctx context.Context, cmd CreatePostCommand) (*dto.PostDTO, error) {
	// 1. Validate input
	if err := uc.validateCommand(cmd); err != nil {
		return nil, err
	}

	// 2. Check rate limit (3 requests per hour per IP)
	rateLimitKey := uc.buildRateLimitKey(cmd.ClientIP)
	allowed, err := uc.rateLimiter.Allow(ctx, rateLimitKey, 3, time.Hour)
	if err != nil {
		// If rate limiter fails, log but continue (graceful degradation)
		// In production, you might want to fail or use a different strategy
		return nil, apperrors.NewDatabaseErrorWithCause("rate limit check failed", err)
	}
	if !allowed {
		return nil, apperrors.NewRateLimitError("rate limit exceeded: maximum 3 posts per hour")
	}

	// 3. Create domain value objects
	company, err := content.NewCompanyName(cmd.Company)
	if err != nil {
		return nil, apperrors.NewValidationErrorWithDetails("invalid company name", map[string]interface{}{
			"error": err.Error(),
		})
	}

	city, err := shared.NewCity(cmd.CityCode, cmd.CityName)
	if err != nil {
		return nil, apperrors.NewValidationErrorWithDetails("invalid city", map[string]interface{}{
			"error": err.Error(),
		})
	}

	postContent, err := content.NewContent(cmd.Content)
	if err != nil {
		return nil, apperrors.NewValidationErrorWithDetails("invalid content", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 4. Create Post entity
	post, err := content.NewPost(company, city, postContent)
	if err != nil {
		return nil, apperrors.NewInternalErrorWithCause("failed to create post", err)
	}

	// 5. Save to repository
	err = uc.repo.Save(ctx, post)
	if err != nil {
		return nil, err
	}

	// 6. Clear related cache
	// Clear city list cache for the city
	cachePattern := fmt.Sprintf("posts:city:%s:*", city.Code())
	err = uc.cacheRepo.DeleteByPattern(ctx, cachePattern)
	if err != nil {
		// Log error but don't fail the operation
		// Cache invalidation failure should not prevent post creation
		// In production, you might want to log this error
	}

	// 7. Convert to DTO and return
	return uc.toDTO(post, cmd.OccurredAt), nil
}

// validateCommand validates the create post command.
func (uc *CreatePostUseCase) validateCommand(cmd CreatePostCommand) error {
	if cmd.Company == "" {
		return apperrors.NewValidationError("company name is required")
	}

	if cmd.CityCode == "" {
		return apperrors.NewValidationError("city code is required")
	}

	if cmd.CityName == "" {
		return apperrors.NewValidationError("city name is required")
	}

	if cmd.Content == "" {
		return apperrors.NewValidationError("content is required")
	}

	if cmd.ClientIP == "" {
		return apperrors.NewValidationError("client IP is required for rate limiting")
	}

	return nil
}

// buildRateLimitKey builds the rate limit key for the given IP.
// Format: "rate_limit:post:{ip}:{hour}"
func (uc *CreatePostUseCase) buildRateLimitKey(ip string) string {
	now := time.Now()
	hour := now.Format("2006-01-02-15") // Format: YYYY-MM-DD-HH
	return fmt.Sprintf("rate_limit:post:%s:%s", ip, hour)
}

// toDTO converts a Post entity to PostDTO.
func (uc *CreatePostUseCase) toDTO(post *content.Post, occurredAt *time.Time) *dto.PostDTO {
	return &dto.PostDTO{
		ID:         post.ID().String(),
		Company:    post.Company().String(),
		CityCode:   post.City().Code(),
		CityName:   post.City().Name(),
		Content:    post.Content().String(),
		OccurredAt: occurredAt,
		CreatedAt:  post.CreatedAt(),
	}
}
