package usecase_test

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/domain"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/usecase"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/pkg/token"
)

// Mock untuk repository
type mockUserRepository struct {
    mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
    args := m.Called(ctx, username)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *mockUserRepository) List(ctx context.Context, page, limit int) ([]*domain.User, int, error) {
    args := m.Called(ctx, page, limit)
    return args.Get(0).([]*domain.User), args.Int(1), args.Error(2)
}

func (m *mockUserRepository) StoreToken(ctx context.Context, userID, token string, expiry time.Duration) error {
    args := m.Called(ctx, userID, token, expiry)
    return args.Error(0)
}

func (m *mockUserRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
    args := m.Called(ctx, token)
    return args.String(0), args.Error(1)
}

func (m *mockUserRepository) DeleteToken(ctx context.Context, token string) error {
    args := m.Called(ctx, token)
    return args.Error(0)
}

func TestRegister(t *testing.T) {
    mockRepo := new(mockUserRepository)
    mockTokenRepo := new(mockUserRepository)
    
    tokenManager := token.NewJWTManager("test-secret", 1*time.Hour)
    timeout := 2 * time.Second
    
    userUsecase := usecase.NewUserUsecase(
        mockRepo,
        mockTokenRepo,
        tokenManager,
        timeout,
    )

    t.Run("Success", func(t *testing.T) {
        mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(nil, errors.New("user not found"))
        mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("user not found"))
        mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
            return user.Username == "testuser" && user.Email == "test@example.com"
        })).Return(nil)

        user, err := userUsecase.Register(context.Background(), "testuser", "test@example.com", "password", "Test User")

        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, "testuser", user.Username)
        assert.Equal(t, "test@example.com", user.Email)
        assert.Equal(t, "Test User", user.Fullname)
        assert.Equal(t, "user", user.Role)
        assert.Empty(t, user.Password) // Password should be stripped from response
        
        mockRepo.AssertExpectations(t)
    })

    t.Run("Username already exists", func(t *testing.T) {
        existingUser := &domain.User{
            ID:       uuid.New().String(),
            Username: "testuser",
        }
        
        mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(existingUser, nil)

        user, err := userUsecase.Register(context.Background(), "testuser", "new@example.com", "password", "Test User")

        assert.Error(t, err)
        assert.Nil(t, user)
        assert.Contains(t, err.Error(), "username already exists")
        
        mockRepo.AssertExpectations(t)
    })

    t.Run("Email already exists", func(t *testing.T) {
        mockRepo.On("GetByUsername", mock.Anything, "newuser").Return(nil, errors.New("user not found"))
        
        existingUser := &domain.User{
            ID:    uuid.New().String(),
            Email: "test@example.com",
        }
        
        mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

        user, err := userUsecase.Register(context.Background(), "newuser", "test@example.com", "password", "Test User")

        assert.Error(t, err)
        assert.Nil(t, user)
        assert.Contains(t, err.Error(), "email already exists")
        
        mockRepo.AssertExpectations(t)
    })
}

func TestLogin(t *testing.T) {
    mockRepo := new(mockUserRepository)
    mockTokenRepo := new(mockUserRepository)
    
    tokenManager := token.NewJWTManager("test-secret", 1*time.Hour)
    timeout := 2 * time.Second
    
    userUsecase := usecase.NewUserUsecase(
        mockRepo,
        mockTokenRepo,
        tokenManager,
        timeout,
    )

    t.Run("Success", func(t *testing.T) {
        // Create a user with bcrypt hashed password for "password"
        hashedPassword := "$2a$10$xVtXMsEH/0rYJ0wAgkxoYe.JsrR0Rw38ScUB4ehAm0/QA9CKgHxxa" // "password"
        
        existingUser := &domain.User{
            ID:       uuid.New().String(),
            Username: "testuser",
            Email:    "test@example.com",
            Password: hashedPassword,
            Fullname: "Test User",
            Role:     "user",
        }
        
        mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(existingUser, nil)
        mockTokenRepo.On("StoreToken", mock.Anything, existingUser.ID, mock.AnythingOfType("string"), tokenManager.GetTokenDuration()).Return(nil)

        user, token, err := userUsecase.Login(context.Background(), "testuser", "password")

        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.NotEmpty(t, token)
        assert.Equal(t, existingUser.ID, user.ID)
        assert.Equal(t, "testuser", user.Username)
        assert.Empty(t, user.Password) // Password should be stripped from response
        
        mockRepo.AssertExpectations(t)
        mockTokenRepo.AssertExpectations(t)
    })

    t.Run("User not found", func(t *testing.T) {
        mockRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, errors.New("user not found"))

        user, token, err := userUsecase.Login(context.Background(), "nonexistent", "password")

        assert.Error(t, err)
        assert.Nil(t, user)
        assert.Empty(t, token)
        assert.Contains(t, err.Error(), "invalid username or password")
        
        mockRepo.AssertExpectations(t)
    })

    t.Run("Wrong password", func(t *testing.T) {
        // Create a user with bcrypt hashed password for "password"
        hashedPassword := "$2a$10$xVtXMsEH/0rYJ0wAgkxoYe.JsrR0Rw38ScUB4ehAm0/QA9CKgHxxa" // "password"
        
        existingUser := &domain.User{
            ID:       uuid.New().String(),
            Username: "testuser",
            Email:    "test@example.com",
            Password: hashedPassword,
            Fullname: "Test User",
            Role:     "user",
        }
        
        mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(existingUser, nil)

        user, token, err := userUsecase.Login(context.Background(), "testuser", "wrongpassword")

        assert.Error(t, err)
        assert.Nil(t, user)
        assert.Empty(t, token)
        assert.Contains(t, err.Error(), "invalid username or password")
        
        mockRepo.AssertExpectations(t)
    })
}