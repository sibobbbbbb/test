package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/domain"
)

type bookRepository struct {
	db *sqlx.DB
}

// NewBookRepository returns a new instance of BookRepository
func NewBookRepository(db *sqlx.DB) domain.BookRepository {
	return &bookRepository{
		db: db,
	}
}

// Method Create menambahkan buku baru
func (r *bookRepository) Create(ctx context.Context, book *domain.Book) error {
	query := `
		INSERT INTO books (id, title, author, isbn, published_year, category_ids, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if book.ID == "" {
		book.ID = uuid.New().String()
	}

	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now

	_, err := r.db.ExecContext(
		ctx,
		query,
		book.ID,
		book.Title,
		book.Author,
		book.ISBN,
		book.PublishedYear,
		pq.Array(book.CategoryIDs),
		book.Stock,
		book.CreatedAt,
		book.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

// GetByID mengambil buku berdasarkan ID
func (r *bookRepository) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	query := `
		SELECT id, title, author, isbn, published_year, category_ids, stock, created_at, updated_at
		FROM books
		WHERE id = $1
	`

	var book domain.Book
	err := r.db.GetContext(ctx, &book, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("book not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return &book, nil
}

// Update memperbarui buku yang sudah ada
func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, isbn = $3, published_year = $4, category_ids = $5, 
		    stock = $6, updated_at = $7
		WHERE id = $8
	`

	book.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		book.Title,
		book.Author,
		book.ISBN,
		book.PublishedYear,
		pq.Array(book.CategoryIDs),
		book.Stock,
		book.UpdatedAt,
		book.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	return nil
}

// Delete menghapus buku berdasarkan ID
func (r *bookRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM books WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

// List mengambil daftar buku dengan pagination
func (r *bookRepository) List(ctx context.Context, page, limit int) ([]*domain.Book, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	countQuery := "SELECT COUNT(*) FROM books"
	var total int
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	query := `
		SELECT id, title, author, isbn, published_year, category_ids, stock, created_at, updated_at
		FROM books
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	books := []*domain.Book{}
	err = r.db.SelectContext(ctx, &books, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list books: %w", err)
	}

	return books, total, nil
}

// Search mencari buku berdasarkan query
func (r *bookRepository) Search(ctx context.Context, query string, page, limit int) ([]*domain.Book, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	searchQuery := "%" + query + "%"

	countQuery := `
		SELECT COUNT(*) FROM books 
		WHERE title ILIKE $1 OR author ILIKE $1 OR isbn ILIKE $1
	`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	selectQuery := `
		SELECT id, title, author, isbn, published_year, category_ids, stock, created_at, updated_at
		FROM books
		WHERE title ILIKE $1 OR author ILIKE $1 OR isbn ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	books := []*domain.Book{}
	err = r.db.SelectContext(ctx, &books, selectQuery, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
	}

	return books, total, nil
}