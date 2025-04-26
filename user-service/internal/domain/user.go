package domain

import (
    "context"
    "time"
)

// Struct untuk user
type User struct {
    ID        string    `json:"id" db:"id"`
    Username  string    `json:"username" db:"username"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"`
    Fullname  string    `json:"fullname" db:"fullname"`
    Role      string    `json:"role" db:"role"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository interface untuk akses data user
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    GetByUsername(ctx context.Context, username string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, page, limit int) ([]*User, int, error)
    
    // Metode untuk token
    StoreToken(ctx context.Context, userID, token string, expiry time.Duration) error
    GetUserIDByToken(ctx context.Context, token string) (string, error)
    DeleteToken(ctx context.Context, token string) error
}

// UserUsecase interface untuk logika bisnis user
type UserUsecase interface {
    Register(ctx context.Context, username, email, password, fullname string) (*User, error)
    Login(ctx context.Context, username, password string) (*User, string, error)
    GetByID(ctx context.Context, id string) (*User, error)
    Update(ctx context.Context, id, username, email, fullname string) (*User, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, page, limit int) ([]*User, int, error)
    ValidateToken(ctx context.Context, token string) (bool, string, string, string, error)
    Logout(ctx context.Context, token string) error
}