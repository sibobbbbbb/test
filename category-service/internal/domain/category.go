package domain

import (
	"context"
	"time"
)

// Category adalah entitas kategori buku
type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CategoryRepository mendefinisikan metode yang dibutuhkan oleh repository
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Category, int, error)
	GetMultipleByIDs(ctx context.Context, ids []string) ([]*Category, error)
	IsExistsByID(ctx context.Context, id string) (bool, error)
}

// CategoryUsecase mendefinisikan business logic yang dibutuhkan
type CategoryUsecase interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Category, int, error)
	GetMultipleByIDs(ctx context.Context, ids []string) ([]*Category, error)
	ValidateCategoryExists(ctx context.Context, id string) (bool, error)
}