package content_test

import (
	"strings"
	"testing"
	"time"

	"fuck_boss/backend/internal/domain/content"
	"fuck_boss/backend/internal/domain/shared"
)

func TestNewPost_Valid(t *testing.T) {
	company, _ := content.NewCompanyName("Example Company")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent(strings.Repeat("A", 50))

	post, err := content.NewPost(company, city, postContent)
	if err != nil {
		t.Fatalf("NewPost() error = %v, want nil", err)
	}

	if post == nil {
		t.Fatal("NewPost() returned nil")
	}

	// Verify ID is generated
	if post.ID().IsZero() {
		t.Error("Post.ID() is zero, want non-zero")
	}

	// Verify company
	if !post.Company().Equals(company) {
		t.Error("Post.Company() does not match input")
	}

	// Verify city
	if !post.City().Equals(city) {
		t.Error("Post.City() does not match input")
	}

	// Verify content
	if !post.Content().Equals(postContent) {
		t.Error("Post.Content() does not match input")
	}

	// Verify createdAt is set
	if post.CreatedAt().IsZero() {
		t.Error("Post.CreatedAt() is zero, want non-zero")
	}

	// Verify createdAt is recent (within last second)
	now := time.Now()
	if post.CreatedAt().After(now) {
		t.Error("Post.CreatedAt() is in the future")
	}
	if now.Sub(post.CreatedAt()) > time.Second {
		t.Error("Post.CreatedAt() is too old")
	}
}

func TestNewPost_GeneratesUniqueIDs(t *testing.T) {
	company, _ := content.NewCompanyName("Example Company")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent(strings.Repeat("A", 50))

	post1, _ := content.NewPost(company, city, postContent)
	post2, _ := content.NewPost(company, city, postContent)

	if post1.ID().Equals(post2.ID()) {
		t.Error("NewPost() generated duplicate IDs")
	}
}

func TestPost_ID(t *testing.T) {
	company, _ := content.NewCompanyName("Example Company")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent(strings.Repeat("A", 50))

	post, _ := content.NewPost(company, city, postContent)

	id := post.ID()
	if id.IsZero() {
		t.Error("Post.ID() is zero")
	}
}

func TestPost_Company(t *testing.T) {
	company, _ := content.NewCompanyName("Example Company")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent(strings.Repeat("A", 50))

	post, _ := content.NewPost(company, city, postContent)

	if !post.Company().Equals(company) {
		t.Error("Post.Company() does not match input")
	}
}

func TestPost_City(t *testing.T) {
	company, _ := content.NewCompanyName("Example Company")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent(strings.Repeat("A", 50))

	post, _ := content.NewPost(company, city, postContent)

	if !post.City().Equals(city) {
		t.Error("Post.City() does not match input")
	}
}

func TestPost_Content(t *testing.T) {
	company, _ := content.NewCompanyName("Example Company")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent(strings.Repeat("A", 50))

	post, _ := content.NewPost(company, city, postContent)

	if !post.Content().Equals(postContent) {
		t.Error("Post.Content() does not match input")
	}
}

func TestPost_CreatedAt(t *testing.T) {
	company, _ := content.NewCompanyName("Example Company")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent(strings.Repeat("A", 50))

	post, _ := content.NewPost(company, city, postContent)

	createdAt := post.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Post.CreatedAt() is zero")
	}

	// Verify it's recent
	now := time.Now()
	if createdAt.After(now) {
		t.Error("Post.CreatedAt() is in the future")
	}
	if now.Sub(createdAt) > time.Second {
		t.Error("Post.CreatedAt() is too old")
	}
}

func TestPost_Publish(t *testing.T) {
	company, _ := content.NewCompanyName("Example Company")
	city, _ := shared.NewCity("beijing", "北京")
	postContent, _ := content.NewContent(strings.Repeat("A", 50))

	post, _ := content.NewPost(company, city, postContent)

	err := post.Publish()
	if err != nil {
		t.Errorf("Post.Publish() error = %v, want nil", err)
	}
}

func TestNewPost_WithDifferentValues(t *testing.T) {
	tests := []struct {
		name    string
		company string
		cityCode string
		cityName string
		content string
	}{
		{
			name:     "normal post",
			company:  "Example Company",
			cityCode: "beijing",
			cityName: "北京",
			content:  strings.Repeat("A", 50),
		},
		{
			name:     "long company name",
			company:  strings.Repeat("A", 100),
			cityCode: "shanghai",
			cityName: "上海",
			content:  strings.Repeat("B", 100),
		},
		{
			name:     "long content",
			company:  "Test Company",
			cityCode: "guangzhou",
			cityName: "广州",
			content:  strings.Repeat("C", 1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			company, _ := content.NewCompanyName(tt.company)
			city, _ := shared.NewCity(tt.cityCode, tt.cityName)
			postContent, _ := content.NewContent(tt.content)

			post, err := content.NewPost(company, city, postContent)
			if err != nil {
				t.Fatalf("NewPost() error = %v, want nil", err)
			}

			if post.Company().String() != tt.company {
				t.Errorf("Post.Company() = %v, want %v", post.Company().String(), tt.company)
			}

			if post.City().Code() != tt.cityCode {
				t.Errorf("Post.City().Code() = %v, want %v", post.City().Code(), tt.cityCode)
			}

			if post.City().Name() != tt.cityName {
				t.Errorf("Post.City().Name() = %v, want %v", post.City().Name(), tt.cityName)
			}
		})
	}
}

