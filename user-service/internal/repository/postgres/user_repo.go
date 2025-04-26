package postgres

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/domain"
    "golang.org/x/crypto/bcrypt"
)

type userRepository struct {
    db *sqlx.DB
}

// NewUserRepository membuat instance baru user repository
func NewUserRepository(db *sqlx.DB) domain.UserRepository {
    return &userRepository{
        db: db,
    }
}

// Create menyimpan user baru ke database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, username, email, password, fullname, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

    if user.ID == "" {
        user.ID = uuid.New().String()
    }

    // Set role default jika kosong
    if user.Role == "" {
        user.Role = "user" // Default role
    }

    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }
    user.Password = string(hashedPassword)

    _, err = r.db.ExecContext(
        ctx,
        query,
        user.ID,
        user.Username,
        user.Email,
        user.Password,
        user.Fullname,
        user.Role,
        user.CreatedAt,
        user.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return nil
}

// GetByID mengambil user berdasarkan ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
    query := `
        SELECT id, username, email, password, fullname, role, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    var user domain.User
    err := r.db.GetContext(ctx, &user, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    return &user, nil
}

// GetByUsername mengambil user berdasarkan username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
    query := `
        SELECT id, username, email, password, fullname, role, created_at, updated_at
        FROM users
        WHERE username = $1
    `

    var user domain.User
    err := r.db.GetContext(ctx, &user, query, username)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    return &user, nil
}

// GetByEmail mengambil user berdasarkan email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT id, username, email, password, fullname, role, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    var user domain.User
    err := r.db.GetContext(ctx, &user, query, email)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    return &user, nil
}

// Update memperbarui data user yang sudah ada
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users
        SET username = $1, email = $2, fullname = $3, updated_at = $4
        WHERE id = $5
    `

    user.UpdatedAt = time.Now()

    _, err := r.db.ExecContext(
        ctx,
        query,
        user.Username,
        user.Email,
        user.Fullname,
        user.UpdatedAt,
        user.ID,
    )
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    return nil
}

// Delete menghapus user berdasarkan ID
func (r *userRepository) Delete(ctx context.Context, id string) error {
    query := "DELETE FROM users WHERE id = $1"

    _, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }

    return nil
}

// List mengambil daftar user dengan pagination
func (r *userRepository) List(ctx context.Context, page, limit int) ([]*domain.User, int, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }

    offset := (page - 1) * limit

    countQuery := "SELECT COUNT(*) FROM users"
    var total int
    err := r.db.GetContext(ctx, &total, countQuery)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count users: %w", err)
    }

    query := `
        SELECT id, username, email, password, fullname, role, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

    var users []*domain.User
    err = r.db.SelectContext(ctx, &users, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to list users: %w", err)
    }

    return users, total, nil
}

// StoreToken, GetUserIDByToken, dan DeleteToken
// akan diimplementasikan di repository Redis
func (r *userRepository) StoreToken(ctx context.Context, userID, token string, expiry time.Duration) error {
    // Implementasi dummy, sebenarnya akan menggunakan Redis
    return nil
}

func (r *userRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
    // Implementasi dummy, sebenarnya akan menggunakan Redis
    return "", errors.New("not implemented in postgres repository")
}

func (r *userRepository) DeleteToken(ctx context.Context, token string) error {
    // Implementasi dummy, sebenarnya akan menggunakan Redis
    return nil
}