package usecase_test

import (
    "context"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/sibobbbbbb/backend-engineer-challenge/category-service/internal/domain"
    "github.com/sibobbbbbb/backend-engineer-challenge/category-service/internal/usecase"
)

// Mock untuk repository
type mockCategoryRepository struct {
    mock.Mock
}

func (m *mockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
    args := m.Called(ctx, category)
    return args.Error(0)
}

func (m *mockCategoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
    args := m.Called(ctx, category)
    return args.Error(0)
}

func (m *mockCategoryRepository) Delete(ctx context.Context, id string) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *mockCategoryRepository) List(ctx context.Context, page, limit int) ([]*domain.Category, int, error) {
    args := m.Called(ctx, page, limit)
    return args.Get(0).([]*domain.Category), args.Int(1), args.Error(2)
}

func (m *mockCategoryRepository) GetMultipleByIDs(ctx context.Context, ids []string) ([]*domain.Category, error) {
    args := m.Called(ctx, ids)
    return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *mockCategoryRepository) IsExistsByID(ctx context.Context, id string) (bool, error) {
    args := m.Called(ctx, id)
    return args.Bool(0), args.Error(1)
}

func TestCreateCategory(t *testing.T) {
    mockRepo := new(mockCategoryRepository)
    timeout := 2 * time.Second
    cacheExp := 30 * time.Minute
    categoryUsecase := usecase.NewCategoryUsecase(mockRepo, timeout, cacheExp)

    now := time.Now()
    category := &domain.Category{
        Name:        "Test Category",
        Description: "Test Description",
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    mockRepo.On("Create", mock.Anything, category).Return(nil)

    err := categoryUsecase.Create(context.Background(), category)

    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestGetCategoryByID(t *testing.T) {
    mockRepo := new(mockCategoryRepository)
    timeout := 2 * time.Second
    cacheExp := 30 * time.Minute
    categoryUsecase := usecase.NewCategoryUsecase(mockRepo, timeout, cacheExp)

    categoryID := uuid.New().String()
    now := time.Now()
    expectedCategory := &domain.Category{
        ID:          categoryID,
        Name:        "Test Category",
        Description: "Test Description",
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    mockRepo.On("GetByID", mock.Anything, categoryID).Return(expectedCategory, nil)

    category, err := categoryUsecase.GetByID(context.Background(), categoryID)

    assert.NoError(t, err)
    assert.Equal(t, expectedCategory.ID, category.ID)
    assert.Equal(t, expectedCategory.Name, category.Name)
    mockRepo.AssertExpectations(t)
}

func TestUpdateCategory(t *testing.T) {
    mockRepo := new(mockCategoryRepository)
    timeout := 2 * time.Second
    cacheExp := 30 * time.Minute
    categoryUsecase := usecase.NewCategoryUsecase(mockRepo, timeout, cacheExp)

    categoryID := uuid.New().String()
    now := time.Now()
    existingCategory := &domain.Category{
        ID:          categoryID,
        Name:        "Original Category",
        Description: "Original Description",
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    updatedCategory := &domain.Category{
        ID:          categoryID,
        Name:        "Updated Category",
        Description: "Updated Description",
    }

    mockRepo.On("GetByID", mock.Anything, categoryID).Return(existingCategory, nil)
    mockRepo.On("Update", mock.Anything, updatedCategory).Return(nil)

    err := categoryUsecase.Update(context.Background(), updatedCategory)

    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestDeleteCategory(t *testing.T) {
    mockRepo := new(mockCategoryRepository)
    timeout := 2 * time.Second
    cacheExp := 30 * time.Minute
    categoryUsecase := usecase.NewCategoryUsecase(mockRepo, timeout, cacheExp)

    categoryID := uuid.New().String()
    now := time.Now()
    existingCategory := &domain.Category{
        ID:          categoryID,
        Name:        "Test Category",
        Description: "Test Description",
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    mockRepo.On("GetByID", mock.Anything, categoryID).Return(existingCategory, nil)
    mockRepo.On("Delete", mock.Anything, categoryID).Return(nil)

    err := categoryUsecase.Delete(context.Background(), categoryID)

    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestListCategories(t *testing.T) {
    mockRepo := new(mockCategoryRepository)
    timeout := 2 * time.Second
    cacheExp := 30 * time.Minute
    categoryUsecase := usecase.NewCategoryUsecase(mockRepo, timeout, cacheExp)

    now := time.Now()
    expectedCategories := []*domain.Category{
        {
            ID:          uuid.New().String(),
            Name:        "Category 1",
            Description: "Description 1",
            CreatedAt:   now,
            UpdatedAt:   now,
        },
        {
            ID:          uuid.New().String(),
            Name:        "Category 2",
            Description: "Description 2",
            CreatedAt:   now,
            UpdatedAt:   now,
        },
    }

    mockRepo.On("List", mock.Anything, 1, 10).Return(expectedCategories, 2, nil)

    categories, total, err := categoryUsecase.List(context.Background(), 1, 10)

    assert.NoError(t, err)
    assert.Equal(t, 2, total)
    assert.Len(t, categories, 2)
    assert.Equal(t, expectedCategories[0].Name, categories[0].Name)
    assert.Equal(t, expectedCategories[1].Name, categories[1].Name)
    mockRepo.AssertExpectations(t)
}

func TestGetMultipleCategories(t *testing.T) {
    mockRepo := new(mockCategoryRepository)
    timeout := 2 * time.Second
    cacheExp := 30 * time.Minute
    categoryUsecase := usecase.NewCategoryUsecase(mockRepo, timeout, cacheExp)

    now := time.Now()
    ids := []string{uuid.New().String(), uuid.New().String()}
    expectedCategories := []*domain.Category{
        {
            ID:          ids[0],
            Name:        "Category 1",
            Description: "Description 1",
            CreatedAt:   now,
            UpdatedAt:   now,
        },
        {
            ID:          ids[1],
            Name:        "Category 2",
            Description: "Description 2",
            CreatedAt:   now,
            UpdatedAt:   now,
        },
    }

    mockRepo.On("GetMultipleByIDs", mock.Anything, ids).Return(expectedCategories, nil)

    categories, err := categoryUsecase.GetMultipleByIDs(context.Background(), ids)

    assert.NoError(t, err)
    assert.Len(t, categories, 2)
    assert.Equal(t, expectedCategories[0].ID, categories[0].ID)
    assert.Equal(t, expectedCategories[1].ID, categories[1].ID)
    mockRepo.AssertExpectations(t)
}

func TestValidateCategoryExists(t *testing.T) {
    mockRepo := new(mockCategoryRepository)
    timeout := 2 * time.Second
    cacheExp := 30 * time.Minute
    categoryUsecase := usecase.NewCategoryUsecase(mockRepo, timeout, cacheExp)

    categoryID := uuid.New().String()

    mockRepo.On("IsExistsByID", mock.Anything, categoryID).Return(true, nil)

    exists, err := categoryUsecase.ValidateCategoryExists(context.Background(), categoryID)

    assert.NoError(t, err)
    assert.True(t, exists)
    mockRepo.AssertExpectations(t)
}