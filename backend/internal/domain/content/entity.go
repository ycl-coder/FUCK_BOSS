// Package content provides domain models for content management.
// It includes value objects, entities, and repository interfaces.
package content

import (
	"time"

	"fuck_boss/backend/internal/domain/shared"
)

// Post represents a post (exposure content) aggregate root.
// It encapsulates business rules and invariants for posts.
type Post struct {
	// id is the unique identifier of the post.
	id PostID

	// company is the company name.
	company CompanyName

	// city is the city where the company is located.
	city shared.City

	// content is the post content.
	content Content

	// createdAt is the time when the post was created.
	createdAt time.Time
}

// NewPost creates a new Post aggregate root.
// It generates a UUID for the ID and sets createdAt to the current time.
// All value objects are validated through their factory methods.
// Returns an error if any validation fails.
func NewPost(company CompanyName, city shared.City, content Content) (*Post, error) {
	// Generate PostID
	id := GeneratePostID()

	// Set createdAt to current time
	createdAt := time.Now()

	// Create Post
	post := &Post{
		id:        id,
		company:   company,
		city:      city,
		content:   content,
		createdAt: createdAt,
	}

	return post, nil
}

// NewPostFromDB creates a Post aggregate root from database data.
// This is used by repositories to reconstruct Post entities from database rows.
// It accepts an existing ID and createdAt timestamp from the database.
// All value objects are validated through their factory methods.
// Returns an error if any validation fails.
func NewPostFromDB(id PostID, company CompanyName, city shared.City, content Content, createdAt time.Time) (*Post, error) {
	// Create Post with provided ID and createdAt
	post := &Post{
		id:        id,
		company:   company,
		city:      city,
		content:   content,
		createdAt: createdAt,
	}

	return post, nil
}

// Publish publishes the post.
// This is a business method that can add validation logic or trigger domain events.
// Currently, it serves as a placeholder for future business logic.
func (p *Post) Publish() error {
	// Future: Add validation logic, trigger domain events, etc.
	// For now, the post is considered published when created.
	return nil
}

// ID returns the Post ID.
func (p *Post) ID() PostID {
	return p.id
}

// Company returns the company name.
func (p *Post) Company() CompanyName {
	return p.company
}

// City returns the city.
func (p *Post) City() shared.City {
	return p.city
}

// Content returns the content.
func (p *Post) Content() Content {
	return p.content
}

// CreatedAt returns the creation time.
func (p *Post) CreatedAt() time.Time {
	return p.createdAt
}
