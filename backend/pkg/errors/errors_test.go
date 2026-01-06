package errors

import (
	"errors"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		expected string
	}{
		{
			name: "error without cause",
			err: &AppError{
				Code:    ErrCodeValidation,
				Message: "invalid input",
			},
			expected: "VALIDATION_ERROR: invalid input",
		},
		{
			name: "error with cause",
			err: &AppError{
				Code:    ErrCodeValidation,
				Message: "invalid input",
				Cause:   errors.New("underlying error"),
			},
			expected: "VALIDATION_ERROR: invalid input: underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &AppError{
		Code:    ErrCodeInternal,
		Message: "something went wrong",
		Cause:   cause,
	}

	if err.Unwrap() != cause {
		t.Errorf("AppError.Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("field is required")

	if err == nil {
		t.Fatal("NewValidationError() returned nil")
	}
	if err.Code != ErrCodeValidation {
		t.Errorf("NewValidationError() Code = %v, want %v", err.Code, ErrCodeValidation)
	}
	if err.Message != "field is required" {
		t.Errorf("NewValidationError() Message = %v, want %v", err.Message, "field is required")
	}
	if err.Cause != nil {
		t.Errorf("NewValidationError() Cause = %v, want nil", err.Cause)
	}
}

func TestNewValidationErrorWithDetails(t *testing.T) {
	details := map[string]interface{}{
		"field":  "company",
		"reason": "empty value",
	}
	err := NewValidationErrorWithDetails("validation failed", details)

	if err == nil {
		t.Fatal("NewValidationErrorWithDetails() returned nil")
	}
	if err.Code != ErrCodeValidation {
		t.Errorf("NewValidationErrorWithDetails() Code = %v, want %v", err.Code, ErrCodeValidation)
	}
	if err.Details["field"] != details["field"] || err.Details["reason"] != details["reason"] {
		t.Errorf("NewValidationErrorWithDetails() Details = %v, want %v", err.Details, details)
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("post")

	if err == nil {
		t.Fatal("NewNotFoundError() returned nil")
	}
	if err.Code != ErrCodeNotFound {
		t.Errorf("NewNotFoundError() Code = %v, want %v", err.Code, ErrCodeNotFound)
	}
	if err.Message != "post not found" {
		t.Errorf("NewNotFoundError() Message = %v, want %v", err.Message, "post not found")
	}
}

func TestNewRateLimitError(t *testing.T) {
	err := NewRateLimitError("too many requests")

	if err == nil {
		t.Fatal("NewRateLimitError() returned nil")
	}
	if err.Code != ErrCodeRateLimit {
		t.Errorf("NewRateLimitError() Code = %v, want %v", err.Code, ErrCodeRateLimit)
	}
	if err.Message != "too many requests" {
		t.Errorf("NewRateLimitError() Message = %v, want %v", err.Message, "too many requests")
	}
}

func TestNewRateLimitErrorWithWindow(t *testing.T) {
	err := NewRateLimitErrorWithWindow("rate limit exceeded", 3600)

	if err == nil {
		t.Fatal("NewRateLimitErrorWithWindow() returned nil")
	}
	if err.Code != ErrCodeRateLimit {
		t.Errorf("NewRateLimitErrorWithWindow() Code = %v, want %v", err.Code, ErrCodeRateLimit)
	}
	if err.Details == nil {
		t.Fatal("NewRateLimitErrorWithWindow() Details is nil")
	}
	if err.Details["retry_after"] != 3600 {
		t.Errorf("NewRateLimitErrorWithWindow() Details[\"retry_after\"] = %v, want %v", err.Details["retry_after"], 3600)
	}
}

func TestNewInternalError(t *testing.T) {
	err := NewInternalError("internal server error")

	if err == nil {
		t.Fatal("NewInternalError() returned nil")
	}
	if err.Code != ErrCodeInternal {
		t.Errorf("NewInternalError() Code = %v, want %v", err.Code, ErrCodeInternal)
	}
	if err.Message != "internal server error" {
		t.Errorf("NewInternalError() Message = %v, want %v", err.Message, "internal server error")
	}
}

func TestNewInternalErrorWithCause(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewInternalErrorWithCause("failed to process", cause)

	if err == nil {
		t.Fatal("NewInternalErrorWithCause() returned nil")
	}
	if err.Code != ErrCodeInternal {
		t.Errorf("NewInternalErrorWithCause() Code = %v, want %v", err.Code, ErrCodeInternal)
	}
	if err.Cause != cause {
		t.Errorf("NewInternalErrorWithCause() Cause = %v, want %v", err.Cause, cause)
	}
}

func TestNewDatabaseError(t *testing.T) {
	err := NewDatabaseError("database connection failed")

	if err == nil {
		t.Fatal("NewDatabaseError() returned nil")
	}
	if err.Code != ErrCodeDatabase {
		t.Errorf("NewDatabaseError() Code = %v, want %v", err.Code, ErrCodeDatabase)
	}
	if err.Message != "database connection failed" {
		t.Errorf("NewDatabaseError() Message = %v, want %v", err.Message, "database connection failed")
	}
}

func TestNewDatabaseErrorWithCause(t *testing.T) {
	cause := errors.New("connection timeout")
	err := NewDatabaseErrorWithCause("database error", cause)

	if err == nil {
		t.Fatal("NewDatabaseErrorWithCause() returned nil")
	}
	if err.Code != ErrCodeDatabase {
		t.Errorf("NewDatabaseErrorWithCause() Code = %v, want %v", err.Code, ErrCodeDatabase)
	}
	if err.Cause != cause {
		t.Errorf("NewDatabaseErrorWithCause() Cause = %v, want %v", err.Cause, cause)
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, "additional context")

	if wrappedErr == nil {
		t.Fatal("Wrap() returned nil")
	}
	if wrappedErr.Error() == "" {
		t.Error("Wrap() returned empty error")
	}

	// Verify error wrapping using errors.Is
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrap() did not properly wrap the error")
	}
}

func TestWrap_NilError(t *testing.T) {
	wrappedErr := Wrap(nil, "context")
	if wrappedErr != nil {
		t.Errorf("Wrap(nil) = %v, want nil", wrappedErr)
	}
}

func TestWrapWithCode(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapWithCode(originalErr, ErrCodeValidation, "validation failed")

	if wrappedErr == nil {
		t.Fatal("WrapWithCode() returned nil")
	}
	if wrappedErr.Code != ErrCodeValidation {
		t.Errorf("WrapWithCode() Code = %v, want %v", wrappedErr.Code, ErrCodeValidation)
	}
	if wrappedErr.Message != "validation failed" {
		t.Errorf("WrapWithCode() Message = %v, want %v", wrappedErr.Message, "validation failed")
	}
	if wrappedErr.Cause != originalErr {
		t.Errorf("WrapWithCode() Cause = %v, want %v", wrappedErr.Cause, originalErr)
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "validation error",
			err:      NewValidationError("invalid"),
			expected: true,
		},
		{
			name:     "not found error",
			err:      NewNotFoundError("post"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "standard error",
			err:      errors.New("standard error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidationError(tt.err); got != tt.expected {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	if !IsNotFoundError(NewNotFoundError("post")) {
		t.Error("IsNotFoundError() = false, want true for not found error")
	}
	if IsNotFoundError(NewValidationError("invalid")) {
		t.Error("IsNotFoundError() = true, want false for validation error")
	}
	if IsNotFoundError(nil) {
		t.Error("IsNotFoundError() = true, want false for nil error")
	}
}

func TestIsRateLimitError(t *testing.T) {
	if !IsRateLimitError(NewRateLimitError("rate limit")) {
		t.Error("IsRateLimitError() = false, want true for rate limit error")
	}
	if IsRateLimitError(NewValidationError("invalid")) {
		t.Error("IsRateLimitError() = true, want false for validation error")
	}
	if IsRateLimitError(nil) {
		t.Error("IsRateLimitError() = true, want false for nil error")
	}
}

func TestIsInternalError(t *testing.T) {
	if !IsInternalError(NewInternalError("internal")) {
		t.Error("IsInternalError() = false, want true for internal error")
	}
	if IsInternalError(NewValidationError("invalid")) {
		t.Error("IsInternalError() = true, want false for validation error")
	}
	if IsInternalError(nil) {
		t.Error("IsInternalError() = true, want false for nil error")
	}
}

func TestIsDatabaseError(t *testing.T) {
	if !IsDatabaseError(NewDatabaseError("database")) {
		t.Error("IsDatabaseError() = false, want true for database error")
	}
	if IsDatabaseError(NewValidationError("invalid")) {
		t.Error("IsDatabaseError() = true, want false for validation error")
	}
	if IsDatabaseError(nil) {
		t.Error("IsDatabaseError() = true, want false for nil error")
	}
}

func TestAs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "AppError",
			err:      NewValidationError("test"),
			expected: true,
		},
		{
			name:     "wrapped AppError",
			err:      WrapWithCode(errors.New("cause"), ErrCodeValidation, "test"),
			expected: true,
		},
		{
			name:     "standard error",
			err:      errors.New("standard"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var appErr *AppError
			result := As(tt.err, &appErr)
			if result != tt.expected {
				t.Errorf("As() = %v, want %v", result, tt.expected)
			}
			if tt.expected && appErr == nil {
				t.Error("As() returned true but appErr is nil")
			}
		})
	}
}

func TestGetCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCode
	}{
		{
			name:     "validation error",
			err:      NewValidationError("test"),
			expected: ErrCodeValidation,
		},
		{
			name:     "not found error",
			err:      NewNotFoundError("post"),
			expected: ErrCodeNotFound,
		},
		{
			name:     "standard error",
			err:      errors.New("standard"),
			expected: "",
		},
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCode(tt.err); got != tt.expected {
				t.Errorf("GetCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetDetails(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected map[string]interface{}
	}{
		{
			name:     "error with details",
			err:      NewValidationErrorWithDetails("test", map[string]interface{}{"field": "name"}),
			expected: map[string]interface{}{"field": "name"},
		},
		{
			name:     "error without details",
			err:      NewValidationError("test"),
			expected: nil,
		},
		{
			name:     "standard error",
			err:      errors.New("standard"),
			expected: nil,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDetails(tt.err)
			if tt.expected == nil && got != nil {
				t.Errorf("GetDetails() = %v, want nil", got)
			}
			if tt.expected != nil {
				if got == nil {
					t.Errorf("GetDetails() = nil, want %v", tt.expected)
				} else if got["field"] != tt.expected["field"] {
					t.Errorf("GetDetails() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}
