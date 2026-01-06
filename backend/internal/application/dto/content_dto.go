// Package dto provides data transfer objects for application layer.
// DTOs are used to transfer data between layers without exposing domain entities.
package dto

import (
	"time"
)

// PostDTO represents a Post data transfer object.
// It is used for API responses and does not contain business logic.
type PostDTO struct {
	// ID is the unique identifier of the post.
	ID string

	// Company is the company name.
	Company string

	// CityCode is the city code (e.g., "beijing").
	CityCode string

	// CityName is the city name (e.g., "北京").
	CityName string

	// Content is the post content.
	Content string

	// OccurredAt is when the incident occurred (optional).
	OccurredAt *time.Time

	// CreatedAt is when the post was created.
	CreatedAt time.Time
}

// PostsListDTO represents a list of posts with pagination information.
type PostsListDTO struct {
	// Posts is the list of posts.
	Posts []*PostDTO

	// Total is the total number of posts (across all pages).
	Total int

	// Page is the current page number (1-based).
	Page int

	// PageSize is the number of items per page.
	PageSize int
}
