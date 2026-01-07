// Package search provides use cases for search functionality.
package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"fuck_boss/backend/internal/application/cache"
	"fuck_boss/backend/internal/application/dto"
	"fuck_boss/backend/internal/domain/content"
	"fuck_boss/backend/internal/domain/shared"
	apperrors "fuck_boss/backend/pkg/errors"
)

// SearchPostsQuery represents the query parameters for searching posts.
type SearchPostsQuery struct {
	// Keyword is the search keyword (required, minimum 2 characters).
	Keyword string

	// CityCode is the city code to filter by (optional).
	// If nil or empty, searches across all cities.
	CityCode *string

	// Page is the page number (1-based, default: 1).
	Page int

	// PageSize is the number of items per page (default: 20).
	PageSize int
}

// SearchPostsUseCase handles searching posts with caching.
// It coordinates domain repository and cache repository.
type SearchPostsUseCase struct {
	// repo is the Post repository.
	repo content.PostRepository

	// cacheRepo is the cache repository for caching query results.
	cacheRepo cache.CacheRepository
}

// NewSearchPostsUseCase creates a new SearchPostsUseCase instance.
func NewSearchPostsUseCase(
	repo content.PostRepository,
	cacheRepo cache.CacheRepository,
) *SearchPostsUseCase {
	return &SearchPostsUseCase{
		repo:      repo,
		cacheRepo: cacheRepo,
	}
}

// Execute executes the search posts query.
// It checks cache first, then queries the repository if cache misses.
func (uc *SearchPostsUseCase) Execute(ctx context.Context, query SearchPostsQuery) (*dto.PostsListDTO, error) {
	// Validate and set defaults
	if err := uc.validateQuery(query); err != nil {
		return nil, err
	}

	page := query.Page
	if page < 1 {
		page = 1
	}

	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	// Build cache key
	cacheKey := uc.buildCacheKey(query.Keyword, query.CityCode, page)

	// Try to get from cache
	cachedData, err := uc.cacheRepo.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		// Cache hit: deserialize and return
		var result dto.PostsListDTO
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			return &result, nil
		}
		// If deserialization fails, continue to query database
	}

	// Cache miss or error: query repository
	var city *shared.City
	if query.CityCode != nil && *query.CityCode != "" {
		cityName := uc.getCityName(*query.CityCode)
		c, err := shared.NewCity(*query.CityCode, cityName)
		if err != nil {
			return nil, apperrors.NewValidationErrorWithDetails("invalid city code", map[string]interface{}{
				"error": err.Error(),
			})
		}
		city = &c
	}

	posts, total, err := uc.repo.Search(ctx, query.Keyword, city, page, pageSize)
	if err != nil {
		return nil, apperrors.NewDatabaseErrorWithCause("failed to search posts", err)
	}

	// Convert to DTO
	result := &dto.PostsListDTO{
		Posts:    uc.toDTOs(posts),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	// Update cache (non-blocking, errors are ignored)
	uc.updateCache(ctx, cacheKey, result)

	return result, nil
}

// validateQuery validates the search posts query.
func (uc *SearchPostsUseCase) validateQuery(query SearchPostsQuery) error {
	// Trim keyword
	keyword := strings.TrimSpace(query.Keyword)

	// Validate keyword
	if keyword == "" {
		return apperrors.NewValidationError("keyword is required")
	}

	// Validate minimum length (2 characters)
	if len(keyword) < 2 {
		return apperrors.NewValidationError("keyword must be at least 2 characters")
	}

	return nil
}

// buildCacheKey builds the cache key for the given search parameters.
// Format: "search:{keyword}:city:{cityCode}:page:{page}" or "search:{keyword}:page:{page}" if no city
func (uc *SearchPostsUseCase) buildCacheKey(keyword string, cityCode *string, page int) string {
	// Normalize keyword (lowercase, trim)
	normalizedKeyword := strings.ToLower(strings.TrimSpace(keyword))

	if cityCode != nil && *cityCode != "" {
		return fmt.Sprintf("search:%s:city:%s:page:%d", normalizedKeyword, *cityCode, page)
	}
	return fmt.Sprintf("search:%s:page:%d", normalizedKeyword, page)
}

// getCityName returns the city name for the given city code.
// This is a temporary implementation. In production, you should use a city repository or configuration.
func (uc *SearchPostsUseCase) getCityName(cityCode string) string {
	// Simple mapping for common cities
	cityMap := map[string]string{
		"beijing":   "北京",
		"shanghai":  "上海",
		"guangzhou": "广州",
		"shenzhen":  "深圳",
		"hangzhou":  "杭州",
		"nanjing":   "南京",
		"chengdu":   "成都",
		"wuhan":     "武汉",
		"xian":      "西安",
		"tianjin":   "天津",
	}

	if name, ok := cityMap[cityCode]; ok {
		return name
	}

	// If not found, return the code as fallback
	// In production, you should query from a city repository or configuration
	return cityCode
}

// getCacheTTL returns the cache TTL for search results.
// Search results are cached for 5 minutes.
func (uc *SearchPostsUseCase) getCacheTTL() time.Duration {
	return 5 * time.Minute
}

// updateCache updates the cache with the query result.
// Errors are ignored to ensure cache failures don't affect the main flow.
func (uc *SearchPostsUseCase) updateCache(ctx context.Context, key string, result *dto.PostsListDTO) {
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

// toDTOs converts a slice of Post entities to PostDTOs.
func (uc *SearchPostsUseCase) toDTOs(posts []*content.Post) []*dto.PostDTO {
	dtos := make([]*dto.PostDTO, 0, len(posts))
	for _, post := range posts {
		dtos = append(dtos, uc.toDTO(post))
	}
	return dtos
}

// toDTO converts a Post entity to PostDTO.
func (uc *SearchPostsUseCase) toDTO(post *content.Post) *dto.PostDTO {
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
