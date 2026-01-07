// Package postgres provides PostgreSQL implementation of domain repositories.
// It implements the PostRepository interface defined in the Domain Layer.
package postgres

import (
	"context"
	"database/sql"
	"time"

	"fuck_boss/backend/internal/domain/content"
	"fuck_boss/backend/internal/domain/shared"
	apperrors "fuck_boss/backend/pkg/errors"
)

// PostRepository is the PostgreSQL implementation of content.PostRepository.
type PostRepository struct {
	// db is the database connection.
	db *sql.DB
}

// NewPostRepository creates a new PostRepository instance.
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

// Save saves a Post to the database.
// If the Post already exists (same ID), it updates the existing record.
// Returns an error if the operation fails.
func (r *PostRepository) Save(ctx context.Context, post *content.Post) error {
	query := `
		INSERT INTO posts (id, company_name, city_code, city_name, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			company_name = EXCLUDED.company_name,
			city_code = EXCLUDED.city_code,
			city_name = EXCLUDED.city_name,
			content = EXCLUDED.content,
			updated_at = EXCLUDED.updated_at
	`

	id := post.ID().String()
	companyName := post.Company().String()
	cityCode := post.City().Code()
	cityName := post.City().Name()
	postContent := post.Content().String()
	createdAt := post.CreatedAt()
	updatedAt := time.Now()

	_, err := r.db.ExecContext(ctx, query,
		id, companyName, cityCode, cityName, postContent, createdAt, updatedAt,
	)
	if err != nil {
		return apperrors.NewDatabaseErrorWithCause("failed to save post", err)
	}

	return nil
}

// FindByID finds a Post by its ID.
// Returns the Post if found, or an error if not found or operation fails.
func (r *PostRepository) FindByID(ctx context.Context, id content.PostID) (*content.Post, error) {
	query := `
		SELECT id, company_name, city_code, city_name, content, created_at
		FROM posts
		WHERE id = $1
	`

	var (
		dbID        string
		companyName string
		cityCode    string
		cityName    string
		postContent string
		createdAt   time.Time
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&dbID, &companyName, &cityCode, &cityName, &postContent, &createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewNotFoundError("post")
		}
		return nil, apperrors.NewDatabaseErrorWithCause("failed to find post by id", err)
	}

	return r.scanPost(dbID, companyName, cityCode, cityName, postContent, createdAt)
}

// FindByCity finds Posts by city with pagination.
// Returns a slice of Posts, total count, and an error.
// The page parameter is 1-based (page 1 is the first page).
// The pageSize parameter specifies the number of items per page.
func (r *PostRepository) FindByCity(ctx context.Context, city shared.City, page, pageSize int) ([]*content.Post, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Query for posts
	query := `
		SELECT id, company_name, city_code, city_name, content, created_at
		FROM posts
		WHERE city_code = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, city.Code(), pageSize, offset)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to find posts by city", err)
	}
	defer rows.Close()

	var posts []*content.Post
	for rows.Next() {
		var (
			dbID        string
			companyName string
			cityCode    string
			cityName    string
			postContent string
			createdAt   time.Time
		)

		if err := rows.Scan(&dbID, &companyName, &cityCode, &cityName, &postContent, &createdAt); err != nil {
			return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to scan post", err)
		}

		post, err := r.scanPost(dbID, companyName, cityCode, cityName, postContent, createdAt)
		if err != nil {
			return nil, 0, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to iterate posts", err)
	}

	// Query for total count
	countQuery := `SELECT COUNT(*) FROM posts WHERE city_code = $1`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery, city.Code()).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to count posts", err)
	}

	return posts, total, nil
}

// FindAll finds all Posts with pagination (across all cities).
// Returns a slice of Posts, total count, and an error.
// The page parameter is 1-based (page 1 is the first page).
// The pageSize parameter specifies the number of items per page.
func (r *PostRepository) FindAll(ctx context.Context, page, pageSize int) ([]*content.Post, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Query for all posts
	query := `
		SELECT id, company_name, city_code, city_name, content, created_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to find all posts", err)
	}
	defer rows.Close()

	var posts []*content.Post
	for rows.Next() {
		var (
			dbID        string
			companyName string
			cityCode    string
			cityName    string
			postContent string
			createdAt   time.Time
		)

		if err := rows.Scan(&dbID, &companyName, &cityCode, &cityName, &postContent, &createdAt); err != nil {
			return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to scan post", err)
		}

		post, err := r.scanPost(dbID, companyName, cityCode, cityName, postContent, createdAt)
		if err != nil {
			return nil, 0, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to iterate posts", err)
	}

	// Query for total count
	countQuery := `SELECT COUNT(*) FROM posts`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to count posts", err)
	}

	return posts, total, nil
}

// Search searches Posts by keyword with optional city filter and pagination.
// If city is nil, searches across all cities.
// Returns a slice of Posts, total count, and an error.
// The page parameter is 1-based (page 1 is the first page).
// The pageSize parameter specifies the number of items per page.
func (r *PostRepository) Search(ctx context.Context, keyword string, city *shared.City, page, pageSize int) ([]*content.Post, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Prepare keyword for full-text search
	// Convert keyword to tsquery format: split by spaces and join with &
	// Escape special characters and wrap each word
	// For simple implementation, use plainto_tsquery which handles this automatically
	// plainto_tsquery converts the input to a tsquery by normalizing and ANDing words

	var query string
	var countQuery string
	var args []interface{}

	if city != nil {
		// Search with city filter
		query = `
			SELECT id, company_name, city_code, city_name, content, created_at
			FROM posts
			WHERE city_code = $1
				AND to_tsvector('simple', company_name || ' ' || content) @@ plainto_tsquery('simple', $2)
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`

		countQuery = `
			SELECT COUNT(*)
			FROM posts
			WHERE city_code = $1
				AND to_tsvector('simple', company_name || ' ' || content) @@ plainto_tsquery('simple', $2)
		`

		args = []interface{}{city.Code(), keyword, pageSize, offset}
	} else {
		// Search across all cities
		query = `
			SELECT id, company_name, city_code, city_name, content, created_at
			FROM posts
			WHERE to_tsvector('simple', company_name || ' ' || content) @@ plainto_tsquery('simple', $1)
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`

		countQuery = `
			SELECT COUNT(*)
			FROM posts
			WHERE to_tsvector('simple', company_name || ' ' || content) @@ plainto_tsquery('simple', $1)
		`

		args = []interface{}{keyword, pageSize, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to search posts", err)
	}
	defer rows.Close()

	var posts []*content.Post
	for rows.Next() {
		var (
			dbID        string
			companyName string
			cityCode    string
			cityName    string
			postContent string
			createdAt   time.Time
		)

		if err := rows.Scan(&dbID, &companyName, &cityCode, &cityName, &postContent, &createdAt); err != nil {
			return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to scan post", err)
		}

		post, err := r.scanPost(dbID, companyName, cityCode, cityName, postContent, createdAt)
		if err != nil {
			return nil, 0, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to iterate posts", err)
	}

	// Query for total count
	var countArgs []interface{}
	if city != nil {
		countArgs = []interface{}{city.Code(), keyword}
	} else {
		countArgs = []interface{}{keyword}
	}

	var total int
	err = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseErrorWithCause("failed to count search results", err)
	}

	return posts, total, nil
}

// scanPost reconstructs a Post entity from database row data.
func (r *PostRepository) scanPost(dbID, companyName, cityCode, cityName, postContent string, createdAt time.Time) (*content.Post, error) {
	// Reconstruct value objects
	postID, err := content.NewPostID(dbID)
	if err != nil {
		return nil, apperrors.NewDatabaseErrorWithCause("invalid post id in database", err)
	}

	company, err := content.NewCompanyName(companyName)
	if err != nil {
		return nil, apperrors.NewDatabaseErrorWithCause("invalid company name in database", err)
	}

	city, err := shared.NewCity(cityCode, cityName)
	if err != nil {
		return nil, apperrors.NewDatabaseErrorWithCause("invalid city in database", err)
	}

	contentVO, err := content.NewContent(postContent)
	if err != nil {
		return nil, apperrors.NewDatabaseErrorWithCause("invalid content in database", err)
	}

	// Create Post from database data using NewPostFromDB
	post, err := content.NewPostFromDB(postID, company, city, contentVO, createdAt)
	if err != nil {
		return nil, apperrors.NewInternalErrorWithCause("failed to reconstruct post", err)
	}

	return post, nil
}
