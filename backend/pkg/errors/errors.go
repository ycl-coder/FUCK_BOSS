// Package errors provides unified error handling mechanism for the application.
// It defines error codes, error types, and error wrapping utilities.
package errors

import (
	"fmt"
)

// ErrorCode represents the type of error.
type ErrorCode string

const (
	// ErrCodeValidation indicates a validation error.
	ErrCodeValidation ErrorCode = "VALIDATION_ERROR"

	// ErrCodeNotFound indicates a resource not found error.
	ErrCodeNotFound ErrorCode = "NOT_FOUND"

	// ErrCodeRateLimit indicates a rate limit exceeded error.
	ErrCodeRateLimit ErrorCode = "RATE_LIMIT_EXCEEDED"

	// ErrCodeInternal indicates an internal server error.
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"

	// ErrCodeDatabase indicates a database error.
	ErrCodeDatabase ErrorCode = "DATABASE_ERROR"
)

// AppError represents an application error with structured information.
type AppError struct {
	// Code is the error code.
	Code ErrorCode

	// Message is the error message (without punctuation).
	Message string

	// Details contains additional error details (optional).
	Details map[string]interface{}

	// Cause is the underlying error (for error wrapping).
	Cause error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error (for error wrapping support).
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a new validation error.
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
	}
}

// NewValidationErrorWithDetails creates a new validation error with details.
func NewValidationErrorWithDetails(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
		Details: details,
	}
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// NewRateLimitError creates a new rate limit error.
func NewRateLimitError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeRateLimit,
		Message: message,
	}
}

// NewRateLimitErrorWithWindow creates a new rate limit error with window information.
func NewRateLimitErrorWithWindow(message string, retryAfter int) *AppError {
	return &AppError{
		Code:    ErrCodeRateLimit,
		Message: message,
		Details: map[string]interface{}{
			"retry_after": retryAfter,
		},
	}
}

// NewInternalError creates a new internal error.
func NewInternalError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeInternal,
		Message: message,
	}
}

// NewInternalErrorWithCause creates a new internal error with a cause.
func NewInternalErrorWithCause(message string, cause error) *AppError {
	return &AppError{
		Code:    ErrCodeInternal,
		Message: message,
		Cause:   cause,
	}
}

// NewDatabaseError creates a new database error.
func NewDatabaseError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeDatabase,
		Message: message,
	}
}

// NewDatabaseErrorWithCause creates a new database error with a cause.
func NewDatabaseErrorWithCause(message string, cause error) *AppError {
	return &AppError{
		Code:    ErrCodeDatabase,
		Message: message,
		Cause:   cause,
	}
}

// Wrap wraps an error with additional context.
// This follows Go 1.13+ error wrapping convention using fmt.Errorf with %w.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapWithCode wraps an error with a specific error code.
func WrapWithCode(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// IsValidationError checks if the error is a validation error.
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code == ErrCodeValidation
	}

	return false
}

// IsNotFoundError checks if the error is a not found error.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code == ErrCodeNotFound
	}

	return false
}

// IsRateLimitError checks if the error is a rate limit error.
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code == ErrCodeRateLimit
	}

	return false
}

// IsInternalError checks if the error is an internal error.
func IsInternalError(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code == ErrCodeInternal
	}

	return false
}

// IsDatabaseError checks if the error is a database error.
func IsDatabaseError(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code == ErrCodeDatabase
	}

	return false
}

// As checks if the error can be converted to AppError.
// This is a helper function similar to errors.As from standard library.
func As(err error, target **AppError) bool {
	if err == nil {
		return false
	}

	if appErr, ok := err.(*AppError); ok {
		*target = appErr
		return true
	}

	// Try to unwrap and check recursively
	unwrapped := err
	for {
		if appErr, ok := unwrapped.(*AppError); ok {
			*target = appErr
			return true
		}

		// Check if error has Unwrap method
		if unwrapper, ok := unwrapped.(interface{ Unwrap() error }); ok {
			unwrapped = unwrapper.Unwrap()
			if unwrapped == nil {
				break
			}
			continue
		}

		break
	}

	return false
}

// GetCode extracts the error code from an error.
// Returns empty string if the error is not an AppError.
func GetCode(err error) ErrorCode {
	if err == nil {
		return ""
	}

	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code
	}

	return ""
}

// GetDetails extracts the details from an error.
// Returns nil if the error is not an AppError or has no details.
func GetDetails(err error) map[string]interface{} {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Details
	}

	return nil
}
