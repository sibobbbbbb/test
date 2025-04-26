package domain

import (
	"context"
	"time"
)

// Struct Book
type Book struct {
	ID            string    `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Author        string    `json:"author" db:"author"`
	ISBN          string    `json:"isbn" db:"isbn"`
	PublishedYear int       `json:"published_year" db:"published_year"`
	CategoryIDs   []string  `json:"category_ids" db:"category_ids"`
	Stock         int       `json:"stock" db:"stock"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// BookRepository methods CRUD for Book
type BookRepository interface {
	Create(ctx context.Context, book *Book) error
	GetByID(ctx context.Context, id string) (*Book, error)
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Book, int, error)
	Search(ctx context.Context, query string, page, limit int) ([]*Book, int, error)
}

// BookUsecase methods for business logic
type BookUsecase interface {
	Create(ctx context.Context, book *Book) error
	GetByID(ctx context.Context, id string) (*Book, error)
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Book, int, error)
	Search(ctx context.Context, query string, page, limit int) ([]*Book, int, error)
	ValidateCategoryIDs(ctx context.Context, categoryIDs []string) error
}