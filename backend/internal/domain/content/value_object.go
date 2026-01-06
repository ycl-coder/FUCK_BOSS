// Package content provides domain models for content management.
// It includes value objects, entities, and repository interfaces.
package content

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// PostID represents a unique identifier for a Post.
// It is a value object that encapsulates the business rules for Post IDs.
type PostID struct {
	// value is the UUID string representation of the Post ID.
	value string
}

// NewPostID creates a new PostID from a UUID string.
// It validates that the string is a valid UUID format.
// Returns an error if the UUID format is invalid.
func NewPostID(value string) (PostID, error) {
	// Trim whitespace
	value = strings.TrimSpace(value)

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return PostID{}, fmt.Errorf("invalid PostID format: %w", err)
	}

	return PostID{value: value}, nil
}

// NewPostIDFromUUID creates a new PostID from a uuid.UUID.
func NewPostIDFromUUID(id uuid.UUID) PostID {
	return PostID{value: id.String()}
}

// GeneratePostID generates a new PostID with a random UUID.
func GeneratePostID() PostID {
	return PostID{value: uuid.New().String()}
}

// String returns the string representation of the PostID.
func (id PostID) String() string {
	return id.value
}

// Value returns the underlying UUID string value.
// This method is provided for cases where the raw value is needed.
func (id PostID) Value() string {
	return id.value
}

// IsZero returns true if the PostID is the zero value.
func (id PostID) IsZero() bool {
	return id.value == ""
}

// Equals returns true if this PostID equals the other PostID.
func (id PostID) Equals(other PostID) bool {
	return id.value == other.value
}

// CompanyName represents a company name.
// It is a value object that encapsulates the business rules for company names.
type CompanyName struct {
	// value is the company name string.
	value string
}

const (
	// MinCompanyNameLength is the minimum length for a company name.
	MinCompanyNameLength = 1
	// MaxCompanyNameLength is the maximum length for a company name.
	MaxCompanyNameLength = 100
)

// NewCompanyName creates a new CompanyName from a string.
// It validates that the string is non-empty and has length between 1 and 100 characters.
// Whitespace is automatically trimmed before validation.
// Returns an error if validation fails.
func NewCompanyName(value string) (CompanyName, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(value)

	// Validate non-empty
	if trimmed == "" {
		return CompanyName{}, fmt.Errorf("company name cannot be empty")
	}

	// Validate length
	length := len([]rune(trimmed)) // Use rune count for proper Unicode support
	if length < MinCompanyNameLength {
		return CompanyName{}, fmt.Errorf("company name must be at least %d character(s)", MinCompanyNameLength)
	}
	if length > MaxCompanyNameLength {
		return CompanyName{}, fmt.Errorf("company name must be at most %d characters", MaxCompanyNameLength)
	}

	return CompanyName{value: trimmed}, nil
}

// String returns the string representation of the CompanyName.
func (cn CompanyName) String() string {
	return cn.value
}

// Value returns the underlying string value.
// This method is provided for cases where the raw value is needed.
func (cn CompanyName) Value() string {
	return cn.value
}

// IsZero returns true if the CompanyName is the zero value.
func (cn CompanyName) IsZero() bool {
	return cn.value == ""
}

// Equals returns true if this CompanyName equals the other CompanyName.
func (cn CompanyName) Equals(other CompanyName) bool {
	return cn.value == other.value
}

// Content represents the content of a Post.
// It is a value object that encapsulates the business rules for content.
type Content struct {
	// value is the content string.
	value string
}

const (
	// MinContentLength is the minimum length for content.
	MinContentLength = 10
	// MaxContentLength is the maximum length for content.
	MaxContentLength = 5000
	// SummaryLength is the maximum length for content summary.
	SummaryLength = 200
)

// NewContent creates a new Content from a string.
// It validates that the string is non-empty and has length between 10 and 5000 characters.
// Whitespace is automatically trimmed before validation.
// Returns an error if validation fails.
func NewContent(value string) (Content, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(value)

	// Validate non-empty
	if trimmed == "" {
		return Content{}, fmt.Errorf("content cannot be empty")
	}

	// Validate length
	length := len([]rune(trimmed)) // Use rune count for proper Unicode support
	if length < MinContentLength {
		return Content{}, fmt.Errorf("content must be at least %d characters", MinContentLength)
	}
	if length > MaxContentLength {
		return Content{}, fmt.Errorf("content must be at most %d characters", MaxContentLength)
	}

	return Content{value: trimmed}, nil
}

// String returns the string representation of the Content.
func (c Content) String() string {
	return c.value
}

// Value returns the underlying string value.
// This method is provided for cases where the raw value is needed.
func (c Content) Value() string {
	return c.value
}

// Summary returns a summary of the content (first 200 characters).
// If the content is longer than 200 characters, it appends "..." to indicate truncation.
func (c Content) Summary() string {
	runes := []rune(c.value)
	if len(runes) <= SummaryLength {
		return c.value
	}
	return string(runes[:SummaryLength]) + "..."
}

// IsZero returns true if the Content is the zero value.
func (c Content) IsZero() bool {
	return c.value == ""
}

// Equals returns true if this Content equals the other Content.
func (c Content) Equals(other Content) bool {
	return c.value == other.value
}
