// Package shared provides shared domain concepts that are used across multiple bounded contexts.
// It includes value objects and domain concepts that are not specific to a single context.
package shared

import (
	"fmt"
	"strings"
)

// City represents a city.
// It is a value object that encapsulates the business rules for cities.
type City struct {
	// code is the city code (e.g., "beijing", "shanghai").
	code string
	
	// name is the city name (e.g., "北京", "上海").
	name string
}

// NewCity creates a new City from code and name.
// It validates that both code and name are non-empty.
// Whitespace is automatically trimmed before validation.
// Returns an error if validation fails.
func NewCity(code, name string) (City, error) {
	// Trim whitespace
	trimmedCode := strings.TrimSpace(code)
	trimmedName := strings.TrimSpace(name)

	// Validate code
	if trimmedCode == "" {
		return City{}, fmt.Errorf("city code cannot be empty")
	}

	// Validate name
	if trimmedName == "" {
		return City{}, fmt.Errorf("city name cannot be empty")
	}

	return City{
		code: trimmedCode,
		name: trimmedName,
	}, nil
}

// Code returns the city code.
func (c City) Code() string {
	return c.code
}

// Name returns the city name.
func (c City) Name() string {
	return c.name
}

// String returns a string representation of the City.
// Format: "City{code: <code>, name: <name>}".
func (c City) String() string {
	return fmt.Sprintf("City{code: %s, name: %s}", c.code, c.name)
}

// IsZero returns true if the City is the zero value.
func (c City) IsZero() bool {
	return c.code == "" && c.name == ""
}

// Equals returns true if this City equals the other City.
func (c City) Equals(other City) bool {
	return c.code == other.code && c.name == other.name
}

