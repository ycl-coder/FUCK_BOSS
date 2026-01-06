package content_test

import (
	"strings"
	"testing"

	"fuck_boss/backend/internal/domain/content"

	"github.com/google/uuid"
)

func TestNewPostID_ValidUUID(t *testing.T) {
	validUUID := uuid.New().String()

	id, err := content.NewPostID(validUUID)
	if err != nil {
		t.Fatalf("NewPostID() error = %v, want nil", err)
	}

	if id.String() != validUUID {
		t.Errorf("PostID.String() = %v, want %v", id.String(), validUUID)
	}
}

func TestNewPostID_InvalidUUID(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "empty string",
			value: "",
		},
		{
			name:  "invalid format",
			value: "not-a-uuid",
		},
		{
			name:  "incomplete UUID",
			value: "123e4567-e89b-12d3-a456",
		},
		{
			name:  "with whitespace around invalid UUID",
			value: " not-a-uuid ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := content.NewPostID(tt.value)
			if err == nil {
				t.Errorf("NewPostID() error = nil, want error")
			}
		})
	}
}

func TestNewPostID_WhitespaceTrimmed(t *testing.T) {
	validUUID := uuid.New().String()
	withWhitespace := "  " + validUUID + "  "

	id, err := content.NewPostID(withWhitespace)
	if err != nil {
		t.Fatalf("NewPostID() error = %v, want nil", err)
	}

	if id.String() != validUUID {
		t.Errorf("PostID.String() = %v, want %v", id.String(), validUUID)
	}
}

func TestNewPostIDFromUUID(t *testing.T) {
	u := uuid.New()
	id := content.NewPostIDFromUUID(u)

	if id.String() != u.String() {
		t.Errorf("PostID.String() = %v, want %v", id.String(), u.String())
	}
}

func TestGeneratePostID(t *testing.T) {
	id1 := content.GeneratePostID()
	id2 := content.GeneratePostID()

	if id1.String() == id2.String() {
		t.Error("GeneratePostID() generated duplicate IDs")
	}

	if id1.IsZero() {
		t.Error("GeneratePostID() returned zero value")
	}

	if id2.IsZero() {
		t.Error("GeneratePostID() returned zero value")
	}
}

func TestPostID_String(t *testing.T) {
	validUUID := uuid.New().String()
	id, _ := content.NewPostID(validUUID)

	if id.String() != validUUID {
		t.Errorf("PostID.String() = %v, want %v", id.String(), validUUID)
	}
}

func TestPostID_Value(t *testing.T) {
	validUUID := uuid.New().String()
	id, _ := content.NewPostID(validUUID)

	if id.Value() != validUUID {
		t.Errorf("PostID.Value() = %v, want %v", id.Value(), validUUID)
	}
}

func TestPostID_IsZero(t *testing.T) {
	tests := []struct {
		name string
		id   content.PostID
		want bool
	}{
		{
			name: "zero value",
			id:   content.PostID{},
			want: true,
		},
		{
			name: "non-zero value",
			id:   content.GeneratePostID(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsZero(); got != tt.want {
				t.Errorf("PostID.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostID_Equals(t *testing.T) {
	validUUID := uuid.New().String()
	id1, _ := content.NewPostID(validUUID)
	id2, _ := content.NewPostID(validUUID)
	id3 := content.GeneratePostID()

	if !id1.Equals(id2) {
		t.Error("PostID.Equals() = false, want true for same UUID")
	}

	if id1.Equals(id3) {
		t.Error("PostID.Equals() = true, want false for different UUIDs")
	}
}

// CompanyName tests

func TestNewCompanyName_Valid(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "min length",
			value: "A",
		},
		{
			name:  "max length",
			value: strings.Repeat("A", 100),
		},
		{
			name:  "normal name",
			value: "Example Company Ltd.",
		},
		{
			name:  "with whitespace",
			value: "  Example Company  ",
		},
		{
			name:  "Unicode characters",
			value: "示例公司",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cn, err := content.NewCompanyName(tt.value)
			if err != nil {
				t.Fatalf("NewCompanyName() error = %v, want nil", err)
			}

			// Check that whitespace is trimmed
			expected := strings.TrimSpace(tt.value)
			if cn.String() != expected {
				t.Errorf("CompanyName.String() = %v, want %v", cn.String(), expected)
			}
		})
	}
}

func TestNewCompanyName_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "empty string",
			value: "",
		},
		{
			name:  "only whitespace",
			value: "   ",
		},
		{
			name:  "too long",
			value: strings.Repeat("A", 101),
		},
		{
			name:  "zero length after trim",
			value: " \t\n ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := content.NewCompanyName(tt.value)
			if err == nil {
				t.Errorf("NewCompanyName() error = nil, want error")
			}
		})
	}
}

func TestNewCompanyName_LengthValidation(t *testing.T) {
	// Test exactly at boundaries
	minValid := "A"
	maxValid := strings.Repeat("A", 100)
	tooShort := ""
	tooLong := strings.Repeat("A", 101)

	// Valid cases
	if _, err := content.NewCompanyName(minValid); err != nil {
		t.Errorf("NewCompanyName() with min length error = %v, want nil", err)
	}

	if _, err := content.NewCompanyName(maxValid); err != nil {
		t.Errorf("NewCompanyName() with max length error = %v, want nil", err)
	}

	// Invalid cases
	if _, err := content.NewCompanyName(tooShort); err == nil {
		t.Error("NewCompanyName() with empty string error = nil, want error")
	}

	if _, err := content.NewCompanyName(tooLong); err == nil {
		t.Error("NewCompanyName() with too long string error = nil, want error")
	}
}

func TestCompanyName_String(t *testing.T) {
	value := "Example Company"
	cn, _ := content.NewCompanyName(value)

	if cn.String() != value {
		t.Errorf("CompanyName.String() = %v, want %v", cn.String(), value)
	}
}

func TestCompanyName_Value(t *testing.T) {
	value := "Example Company"
	cn, _ := content.NewCompanyName(value)

	if cn.Value() != value {
		t.Errorf("CompanyName.Value() = %v, want %v", cn.Value(), value)
	}
}

func TestCompanyName_IsZero(t *testing.T) {
	tests := []struct {
		name string
		cn   content.CompanyName
		want bool
	}{
		{
			name: "zero value",
			cn:   content.CompanyName{},
			want: true,
		},
		{
			name: "non-zero value",
			cn:   func() content.CompanyName { cn, _ := content.NewCompanyName("Test"); return cn }(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cn.IsZero(); got != tt.want {
				t.Errorf("CompanyName.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyName_Equals(t *testing.T) {
	value := "Example Company"
	cn1, _ := content.NewCompanyName(value)
	cn2, _ := content.NewCompanyName(value)
	cn3, _ := content.NewCompanyName("Different Company")

	if !cn1.Equals(cn2) {
		t.Error("CompanyName.Equals() = false, want true for same value")
	}

	if cn1.Equals(cn3) {
		t.Error("CompanyName.Equals() = true, want false for different values")
	}
}

func TestCompanyName_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "leading whitespace",
			input:    "  Company",
			expected: "Company",
		},
		{
			name:     "trailing whitespace",
			input:    "Company  ",
			expected: "Company",
		},
		{
			name:     "both sides whitespace",
			input:    "  Company  ",
			expected: "Company",
		},
		{
			name:     "internal whitespace preserved",
			input:    "  Company Name  ",
			expected: "Company Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cn, err := content.NewCompanyName(tt.input)
			if err != nil {
				t.Fatalf("NewCompanyName() error = %v, want nil", err)
			}

			if cn.String() != tt.expected {
				t.Errorf("CompanyName.String() = %v, want %v", cn.String(), tt.expected)
			}
		})
	}
}

// Content tests

func TestNewContent_Valid(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "min length",
			value: strings.Repeat("A", 10),
		},
		{
			name:  "max length",
			value: strings.Repeat("A", 5000),
		},
		{
			name:  "normal content",
			value: "This is a normal content that describes the company's misconduct in detail.",
		},
		{
			name:  "with whitespace",
			value: "  This is content with whitespace  ",
		},
		{
			name:  "Unicode characters",
			value: "这是一段中文内容，用于测试Unicode字符的支持情况。",
		},
		{
			name:  "multiline content",
			value: "Line 1\nLine 2\nLine 3\nLine 4\nLine 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := content.NewContent(tt.value)
			if err != nil {
				t.Fatalf("NewContent() error = %v, want nil", err)
			}

			// Check that whitespace is trimmed
			expected := strings.TrimSpace(tt.value)
			if c.String() != expected {
				t.Errorf("Content.String() = %v, want %v", c.String(), expected)
			}
		})
	}
}

func TestNewContent_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "empty string",
			value: "",
		},
		{
			name:  "only whitespace",
			value: "   ",
		},
		{
			name:  "too short",
			value: "short",
		},
		{
			name:  "exactly 9 characters",
			value: strings.Repeat("A", 9),
		},
		{
			name:  "too long",
			value: strings.Repeat("A", 5001),
		},
		{
			name:  "zero length after trim",
			value: " \t\n ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := content.NewContent(tt.value)
			if err == nil {
				t.Errorf("NewContent() error = nil, want error")
			}
		})
	}
}

func TestNewContent_LengthValidation(t *testing.T) {
	// Test exactly at boundaries
	minValid := strings.Repeat("A", 10)
	maxValid := strings.Repeat("A", 5000)
	tooShort := strings.Repeat("A", 9)
	tooLong := strings.Repeat("A", 5001)

	// Valid cases
	if _, err := content.NewContent(minValid); err != nil {
		t.Errorf("NewContent() with min length error = %v, want nil", err)
	}

	if _, err := content.NewContent(maxValid); err != nil {
		t.Errorf("NewContent() with max length error = %v, want nil", err)
	}

	// Invalid cases
	if _, err := content.NewContent(tooShort); err == nil {
		t.Error("NewContent() with too short string error = nil, want error")
	}

	if _, err := content.NewContent(tooLong); err == nil {
		t.Error("NewContent() with too long string error = nil, want error")
	}
}

func TestContent_String(t *testing.T) {
	value := "This is a test content that is long enough to pass validation."
	c, _ := content.NewContent(value)

	if c.String() != value {
		t.Errorf("Content.String() = %v, want %v", c.String(), value)
	}
}

func TestContent_Value(t *testing.T) {
	value := "This is a test content that is long enough to pass validation."
	c, _ := content.NewContent(value)

	if c.Value() != value {
		t.Errorf("Content.Value() = %v, want %v", c.Value(), value)
	}
}

func TestContent_Summary(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "short content",
			content: strings.Repeat("A", 50),
			want:    strings.Repeat("A", 50),
		},
		{
			name:    "exactly 200 characters",
			content: strings.Repeat("A", 200),
			want:    strings.Repeat("A", 200),
		},
		{
			name:    "longer than 200 characters",
			content: strings.Repeat("A", 300),
			want:    strings.Repeat("A", 200) + "...",
		},
		{
			name:    "very long content",
			content: strings.Repeat("A", 1000),
			want:    strings.Repeat("A", 200) + "...",
		},
		{
			name:    "Unicode characters",
			content: strings.Repeat("中", 100),
			want:    strings.Repeat("中", 100),
		},
		{
			name:    "Unicode longer than 200",
			content: strings.Repeat("中", 300),
			want:    strings.Repeat("中", 200) + "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create content with valid length
			c, err := content.NewContent(tt.content)
			if err != nil {
				t.Fatalf("NewContent() error = %v, want nil", err)
			}

			summary := c.Summary()
			if summary != tt.want {
				t.Errorf("Content.Summary() = %v, want %v", summary, tt.want)
			}
		})
	}
}

func TestContent_IsZero(t *testing.T) {
	tests := []struct {
		name string
		c    content.Content
		want bool
	}{
		{
			name: "zero value",
			c:    content.Content{},
			want: true,
		},
		{
			name: "non-zero value",
			c:    func() content.Content { c, _ := content.NewContent(strings.Repeat("A", 10)); return c }(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsZero(); got != tt.want {
				t.Errorf("Content.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContent_Equals(t *testing.T) {
	value := strings.Repeat("A", 50)
	c1, _ := content.NewContent(value)
	c2, _ := content.NewContent(value)
	c3, _ := content.NewContent(strings.Repeat("B", 50))

	if !c1.Equals(c2) {
		t.Error("Content.Equals() = false, want true for same value")
	}

	if c1.Equals(c3) {
		t.Error("Content.Equals() = true, want false for different values")
	}
}

func TestContent_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "leading whitespace",
			input:    "  " + strings.Repeat("A", 10),
			expected: strings.Repeat("A", 10),
		},
		{
			name:     "trailing whitespace",
			input:    strings.Repeat("A", 10) + "  ",
			expected: strings.Repeat("A", 10),
		},
		{
			name:     "both sides whitespace",
			input:    "  " + strings.Repeat("A", 10) + "  ",
			expected: strings.Repeat("A", 10),
		},
		{
			name:     "internal whitespace preserved",
			input:    "  " + strings.Repeat("A", 5) + " " + strings.Repeat("B", 5) + "  ",
			expected: strings.Repeat("A", 5) + " " + strings.Repeat("B", 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := content.NewContent(tt.input)
			if err != nil {
				t.Fatalf("NewContent() error = %v, want nil", err)
			}

			if c.String() != tt.expected {
				t.Errorf("Content.String() = %v, want %v", c.String(), tt.expected)
			}
		})
	}
}
