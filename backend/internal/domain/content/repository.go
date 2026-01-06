// Package content provides domain models for content management.
// It includes value objects, entities, and repository interfaces.
package content

import (
	"context"

	"fuck_boss/backend/internal/domain/shared"
)

// PostRepository defines the interface for Post persistence.
// It follows the Dependency Inversion Principle by defining the interface
// in the Domain Layer, while implementations are in the Infrastructure Layer.
type PostRepository interface {
	// Save saves a Post to the repository.
	// If the Post already exists (same ID), it updates the existing record.
	// Returns an error if the operation fails.
	Save(ctx context.Context, post *Post) error

	// FindByID finds a Post by its ID.
	// Returns the Post if found, or an error if not found or operation fails.
	FindByID(ctx context.Context, id PostID) (*Post, error)

	// FindByCity finds Posts by city with pagination.
	// Returns a slice of Posts, total count, and an error.
	// The page parameter is 1-based (page 1 is the first page).
	// The pageSize parameter specifies the number of items per page.
	FindByCity(ctx context.Context, city shared.City, page, pageSize int) ([]*Post, int, error)

	// Search searches Posts by keyword with optional city filter and pagination.
	// If city is nil, searches across all cities.
	// Returns a slice of Posts, total count, and an error.
	// The page parameter is 1-based (page 1 is the first page).
	// The pageSize parameter specifies the number of items per page.
	Search(ctx context.Context, keyword string, city *shared.City, page, pageSize int) ([]*Post, int, error)
}
