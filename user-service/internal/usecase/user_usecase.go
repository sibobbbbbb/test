package usecase

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/domain"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/pkg/token"
    "golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
    userRepo      domain.UserRepository
    tokenManager  *token.JWTManager
    contextTimeout time.Duration
}

// NewUserUsecase menciptakan instance baru dari UserUsecase
func NewUserUsecase(ur domain.UserRepository, tm *token.JWTManager, timeout time.Duration) domain.UserUsecase {
    return &userUsecase{
        userRepo:      ur,
        tokenManager:  tm,
        contextTimeout: timeout,
    }
}

// Register mendaftarkan user baru
func (u *userUsecase) Register(ctx context.Context, username, email, password, fullname string) (*domain.User, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    // Validasi data
    if username == "" {
        return nil, errors.New("username is required")
    }
    if email == "" {
        return nil, errors.New("email is required")
    }
    if password == "" {
        return nil, errors.New("password is required")
    }
    if fullname == "" {
        return nil, errors.New("fullname is required")
    }

    // Cek apakah username sudah digunakan
    existingUser, err := u.userRepo.GetByUsername(ctx, username)
    if err == nil && existingUser != nil {
        return nil, errors.New("username already exists")
    }

    // Cek apakah email sudah digunakan
    existingUser, err = u.userRepo.GetByEmail(ctx, email)
    if err == nil && existingUser != nil {
        return nil, errors.New("email already exists")
    }

    // Buat user baru
    user := &domain.User{
        Username: username,
        Email:    email,
        Password: password, // Password akan di-hash di repository
        Fullname: fullname,
        Role:     "user", // Default role
    }

    // Simpan ke database
    if err := u.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // Hapus password dari respons
    user.Password = ""

    return user, nil
}

// Login melakukan login user
func (u *userUsecase) Login(ctx context.Context, username, password string) (*domain.User, string, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    // Validasi
    if username == "" {
        return nil, "", errors.New("username is required")
    }
    if password == "" {
        return nil, "", errors.New("password is required")
    }

    // Cari user berdasarkan username
    user, err := u.userRepo.GetByUsername(ctx, username)
    if err != nil {
        return nil, "", fmt.Errorf("invalid username or password")
    }

    // Verifikasi password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, "", fmt.Errorf("invalid username or password")
    }

    // Generate token
    token, err := u.tokenManager.Generate(user.ID, user.Username, user.Role)
    if err != nil {
        return nil, "", fmt.Errorf("failed to generate token: %w", err)
    }

    // Simpan token di Redis
    err = u.userRepo.StoreToken(ctx, user.ID, token, u.tokenManager.GetTokenDuration())
    if err != nil {
        return nil, "", fmt.Errorf("failed to store token: %w", err)
    }

    // Hapus password dari respons
    user.Password = ""

    return user, token, nil
}

// GetByID mendapatkan user berdasarkan ID
func (u *userUsecase) GetByID(ctx context.Context, id string) (*domain.User, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    user, err := u.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Hapus password dari respons
    user.Password = ""

    return user, nil
}

// Update memperbarui data user
func (u *userUsecase) Update(ctx context.Context, id, username, email, fullname string) (*domain.User, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    // Validasi
    if username == "" {
        return nil, errors.New("username is required")
    }
    if email == "" {
        return nil, errors.New("email is required")
    }
    if fullname == "" {
        return nil, errors.New("fullname is required")
    }

    // Cek apakah user ada
    existingUser, err := u.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Cek apakah username sudah digunakan oleh user lain
    if username != existingUser.Username {
        userWithSameUsername, err := u.userRepo.GetByUsername(ctx, username)
        if err == nil && userWithSameUsername != nil && userWithSameUsername.ID != id {
            return nil, errors.New("username already exists")
        }
    }

    // Cek apakah email sudah digunakan oleh user lain
    if email != existingUser.Email {
        userWithSameEmail, err := u.userRepo.GetByEmail(ctx, email)
        if err == nil && userWithSameEmail != nil && userWithSameEmail.ID != id {
            return nil, errors.New("email already exists")
        }
    }

    // Update user
    existingUser.Username = username
    existingUser.Email = email
    existingUser.Fullname = fullname

    if err := u.userRepo.Update(ctx, existingUser); err != nil {
        return nil, err
    }

    // Hapus password dari respons
    existingUser.Password = ""

    return existingUser, nil
}

// Delete menghapus user
func (u *userUsecase) Delete(ctx context.Context, id string) error {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    // Cek apakah user ada
    _, err := u.userRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    return u.userRepo.Delete(ctx, id)
}

// List mendapatkan daftar user
func (u *userUsecase) List(ctx context.Context, page, limit int) ([]*domain.User, int, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    users, total, err := u.userRepo.List(ctx, page, limit)
    if err != nil {
        return nil, 0, err
    }

    // Hapus password dari respons
    for _, user := range users {
        user.Password = ""
    }

    return users, total, nil
}

// ValidateToken memvalidasi token JWT
func (u *userUsecase) ValidateToken(ctx context.Context, tokenStr string) (bool, string, string, string, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    // Validasi token
    claims, err := u.tokenManager.Validate(tokenStr)
    if err != nil {
        return false, "", "", "", err
    }

    // Cek apakah token ada di Redis
    userID, err := u.userRepo.GetUserIDByToken(ctx, tokenStr)
    if err != nil {
        return false, "", "", "", fmt.Errorf("token is invalid or expired: %w", err)
    }

    // Pastikan userID di token sama dengan yang di Redis
    if claims.UserID != userID {
        return false, "", "", "", errors.New("token is invalid")
    }

    return true, claims.UserID, claims.Username, claims.Role, nil
}

// Logout menghapus token dari penyimpanan
func (u *userUsecase) Logout(ctx context.Context, token string) error {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    return u.userRepo.DeleteToken(ctx, token)
}