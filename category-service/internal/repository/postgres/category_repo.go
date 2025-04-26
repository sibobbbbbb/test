package postgres

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/sibobbbbbb/backend-engineer-challenge/category-service/internal/domain"
)

type categoryRepository struct {
    db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) domain.CategoryRepository {
    return &categoryRepository{
        db: db,
    }
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
    query := `
        INSERT INTO categories (id, name, description, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `

    if category.ID == "" {
        category.ID = uuid.New().String()
    }

    now := time.Now()
    category.CreatedAt = now
    category.UpdatedAt = now

    _, err := r.db.ExecContext(
        ctx,
        query,
        category.ID,
        category.Name,
        category.Description,
        category.CreatedAt,
        category.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to create category: %w", err)
    }

    return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
    query := `
        SELECT id, name, description, created_at, updated_at
        FROM categories
        WHERE id = $1
    `

    var category domain.Category
    err := r.db.GetContext(ctx, &category, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("category not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get category: %w", err)
    }

    return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
    query := `
        UPDATE categories
        SET name = $1, description = $2, updated_at = $3
        WHERE id = $4
    `

    category.UpdatedAt = time.Now()

    _, err := r.db.ExecContext(
        ctx,
        query,
        category.Name,
        category.Description,
        category.UpdatedAt,
        category.ID,
    )
    if err != nil {
        return fmt.Errorf("failed to update category: %w", err)
    }

    return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
    query := "DELETE FROM categories WHERE id = $1"

    _, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete category: %w", err)
    }

    return nil
}

func (r *categoryRepository) List(ctx context.Context, page, limit int) ([]*domain.Category, int, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }

    offset := (page - 1) * limit

    countQuery := "SELECT COUNT(*) FROM categories"
    var total int
    err := r.db.GetContext(ctx, &total, countQuery)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count categories: %w", err)
    }

    query := `
        SELECT id, name, description, created_at, updated_at
        FROM categories
        ORDER BY name ASC
        LIMIT $1 OFFSET $2
    `

    categories := []*domain.Category{}
    err = r.db.SelectContext(ctx, &categories, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to list categories: %w", err)
    }

    return categories, total, nil
}

func (r *categoryRepository) GetMultipleByIDs(ctx context.Context, ids []string) ([]*domain.Category, error) {
    if len(ids) == 0 {
        return []*domain.Category{}, nil
    }

    query, args, err := sqlx.In(`
        SELECT id, name, description, created_at, updated_at
        FROM categories
        WHERE id IN (?)
    `, ids)
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }

    query = r.db.Rebind(query)
    categories := []*domain.Category{}
    err = r.db.SelectContext(ctx, &categories, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to get categories by IDs: %w", err)
    }

    return categories, nil
}

func (r *categoryRepository) IsExistsByID(ctx context.Context, id string) (bool, error) {
    query := "SELECT COUNT(*) FROM categories WHERE id = $1"
    
    var count int
    err := r.db.GetContext(ctx, &count, query, id)
    if err != nil {
        return false, fmt.Errorf("failed to check category existence: %w", err)
    }
    
    return count > 0, nil
}