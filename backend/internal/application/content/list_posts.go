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
	"fuck_boss/backend/internal/domain/shared"
	apperrors "fuck_boss/backend/pkg/errors"
)

// ListPostsQuery represents the query parameters for listing posts.
type ListPostsQuery struct {
	// CityCode is the city code to filter by (required).
	CityCode string

	// Page is the page number (1-based, default: 1).
	Page int

	// PageSize is the number of items per page (default: 20).
	PageSize int
}

// ListPostsUseCase handles listing posts by city with caching.
// It coordinates domain repository and cache repository.
type ListPostsUseCase struct {
	// repo is the Post repository.
	repo content.PostRepository

	// cacheRepo is the cache repository for caching query results.
	cacheRepo cache.CacheRepository
}

// NewListPostsUseCase creates a new ListPostsUseCase instance.
func NewListPostsUseCase(
	repo content.PostRepository,
	cacheRepo cache.CacheRepository,
) *ListPostsUseCase {
	return &ListPostsUseCase{
		repo:      repo,
		cacheRepo: cacheRepo,
	}
}

// Execute executes the list posts query.
// It checks cache first, then queries the repository if cache misses.
func (uc *ListPostsUseCase) Execute(ctx context.Context, query ListPostsQuery) (*dto.PostsListDTO, error) {
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
	var cacheKey string
	var city *shared.City
	if query.CityCode != "" {
		// Create city value object
		c, err := shared.NewCity(query.CityCode, uc.getCityName(query.CityCode))
		if err != nil {
			return nil, apperrors.NewValidationErrorWithDetails("invalid city code", map[string]interface{}{
				"error": err.Error(),
			})
		}
		city = &c
		cacheKey = uc.buildCacheKey(city.Code(), page)
	} else {
		// All cities
		cacheKey = uc.buildCacheKey("all", page)
	}

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
	var posts []*content.Post
	var total int
	if city != nil {
		// Query by city
		posts, total, err = uc.repo.FindByCity(ctx, *city, page, pageSize)
	} else {
		// Query all cities
		posts, total, err = uc.repo.FindAll(ctx, page, pageSize)
	}
	if err != nil {
		return nil, apperrors.NewDatabaseErrorWithCause("failed to query posts", err)
	}

	// Convert to DTO
	result := &dto.PostsListDTO{
		Posts:    uc.toDTOs(posts),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	// Update cache (non-blocking, errors are ignored)
	cityCode := "all"
	if city != nil {
		cityCode = city.Code()
	}
	uc.updateCache(ctx, cacheKey, result, cityCode)

	return result, nil
}

// validateQuery validates the list posts query.
// CityCode is optional - if empty, returns posts from all cities.
func (uc *ListPostsUseCase) validateQuery(query ListPostsQuery) error {
	// CityCode is optional, no validation needed
	return nil
}

// buildCacheKey builds the cache key for the given city and page.
// Format: "posts:city:{cityCode}:page:{page}"
func (uc *ListPostsUseCase) buildCacheKey(cityCode string, page int) string {
	return fmt.Sprintf("posts:city:%s:page:%d", cityCode, page)
}

// getCityName returns the city name for the given city code.
// This is a temporary implementation. In production, you should use a city repository or configuration.
func (uc *ListPostsUseCase) getCityName(cityCode string) string {
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

// getCacheTTL returns the cache TTL based on city popularity.
// Popular cities (beijing, shanghai) get shorter TTL (5 minutes).
// Other cities get longer TTL (10 minutes).
func (uc *ListPostsUseCase) getCacheTTL(cityCode string) time.Duration {
	popularCities := map[string]bool{
		"beijing":  true,
		"shanghai": true,
		"guangzhou": true,
		"shenzhen": true,
	}

	if popularCities[cityCode] {
		return 5 * time.Minute
	}
	return 10 * time.Minute
}

// updateCache updates the cache with the query result.
// Errors are ignored to ensure cache failures don't affect the main flow.
func (uc *ListPostsUseCase) updateCache(ctx context.Context, key string, result *dto.PostsListDTO, cityCode string) {
	// Serialize to JSON
	data, err := json.Marshal(result)
	if err != nil {
		// Log error but don't fail
		return
	}

	// Get TTL based on city popularity
	ttl := uc.getCacheTTL(cityCode)

	// Set cache (non-blocking)
	_ = uc.cacheRepo.Set(ctx, key, string(data), ttl)
}

// toDTOs converts a slice of Post entities to PostDTOs.
func (uc *ListPostsUseCase) toDTOs(posts []*content.Post) []*dto.PostDTO {
	dtos := make([]*dto.PostDTO, 0, len(posts))
	for _, post := range posts {
		dtos = append(dtos, uc.toDTO(post))
	}
	return dtos
}

// toDTO converts a Post entity to PostDTO.
func (uc *ListPostsUseCase) toDTO(post *content.Post) *dto.PostDTO {
	return &dto.PostDTO{
		ID:        post.ID().String(),
		Company:   post.Company().String(),
		CityCode:  post.City().Code(),
		CityName:  post.City().Name(),
		Content:   post.Content().String(),
		OccurredAt: nil, // Not stored in Post entity
		CreatedAt: post.CreatedAt(),
	}
}

