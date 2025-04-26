package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/domain"
	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/usecase"
)

// Mock untuk repository
type mockBookRepository struct {
	mock.Mock
}

func (m *mockBookRepository) Create(ctx context.Context, book *domain.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *mockBookRepository) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *mockBookRepository) Update(ctx context.Context, book *domain.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *mockBookRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockBookRepository) List(ctx context.Context, page, limit int) ([]*domain.Book, int, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]*domain.Book), args.Int(1), args.Error(2)
}

func (m *mockBookRepository) Search(ctx context.Context, query string, page, limit int) ([]*domain.Book, int, error) {
	args := m.Called(ctx, query, page, limit)
	return args.Get(0).([]*domain.Book), args.Int(1), args.Error(2)
}

// Mock untuk category client
type mockCategoryClient struct {
	mock.Mock
}

func (m *mockCategoryClient) ValidateCategory(ctx context.Context, categoryID string) (bool, error) {
	args := m.Called(ctx, categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *mockCategoryClient) GetCategoryByID(ctx context.Context, categoryID string) (interface{}, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0), args.Error(1)
}

func (m *mockCategoryClient) GetMultipleCategories(ctx context.Context, categoryIDs []string) (interface{}, error) {
	args := m.Called(ctx, categoryIDs)
	return args.Get(0), args.Error(1)
}

func TestCreateBook(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	now := time.Now()
	book := &domain.Book{
		Title:         "Test Book",
		Author:        "Test Author",
		ISBN:          "1234567890",
		PublishedYear: 2022,
		CategoryIDs:   []string{"cat1", "cat2"},
		Stock:         10,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Setup expectations
	mockCatClient.On("ValidateCategory", mock.Anything, "cat1").Return(true, nil)
	mockCatClient.On("ValidateCategory", mock.Anything, "cat2").Return(true, nil)
	mockRepo.On("Create", mock.Anything, book).Return(nil)

	// Execute
	err := bookUsecase.Create(context.Background(), book)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCatClient.AssertExpectations(t)
}

func TestCreateBookWithInvalidCategory(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	now := time.Now()
	book := &domain.Book{
		Title:         "Test Book",
		Author:        "Test Author",
		ISBN:          "1234567890",
		PublishedYear: 2022,
		CategoryIDs:   []string{"cat1", "invalidCat"},
		Stock:         10,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Setup expectations
	mockCatClient.On("ValidateCategory", mock.Anything, "cat1").Return(true, nil)
	mockCatClient.On("ValidateCategory", mock.Anything, "invalidCat").Return(false, nil)

	// Execute
	err := bookUsecase.Create(context.Background(), book)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
	mockRepo.AssertNotCalled(t, "Create")
	mockCatClient.AssertExpectations(t)
}

func TestGetBook(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	bookID := uuid.New().String()
	now := time.Now()
	expectedBook := &domain.Book{
		ID:            bookID,
		Title:         "Test Book",
		Author:        "Test Author",
		ISBN:          "1234567890",
		PublishedYear: 2022,
		CategoryIDs:   []string{"cat1", "cat2"},
		Stock:         10,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Setup expectations
	mockRepo.On("GetByID", mock.Anything, bookID).Return(expectedBook, nil)

	// Execute
	book, err := bookUsecase.GetByID(context.Background(), bookID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedBook.ID, book.ID)
	assert.Equal(t, expectedBook.Title, book.Title)
	mockRepo.AssertExpectations(t)
}

func TestGetBookNotFound(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	bookID := uuid.New().String()
	expectedErr := errors.New("book not found")

	// Setup expectations
	mockRepo.On("GetByID", mock.Anything, bookID).Return(nil, expectedErr)

	// Execute
	book, err := bookUsecase.GetByID(context.Background(), bookID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBook(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	bookID := uuid.New().String()
	now := time.Now()
	existingBook := &domain.Book{
		ID:            bookID,
		Title:         "Test Book",
		Author:        "Test Author",
		ISBN:          "1234567890",
		PublishedYear: 2022,
		CategoryIDs:   []string{"cat1", "cat2"},
		Stock:         10,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	updatedBook := &domain.Book{
		ID:            bookID,
		Title:         "Updated Book",
		Author:        "Updated Author",
		ISBN:          "0987654321",
		PublishedYear: 2023,
		CategoryIDs:   []string{"cat3", "cat4"},
		Stock:         20,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Setup expectations
	mockRepo.On("GetByID", mock.Anything, bookID).Return(existingBook, nil)
	mockCatClient.On("ValidateCategory", mock.Anything, "cat3").Return(true, nil)
	mockCatClient.On("ValidateCategory", mock.Anything, "cat4").Return(true, nil)
	mockRepo.On("Update", mock.Anything, updatedBook).Return(nil)

	// Execute
	err := bookUsecase.Update(context.Background(), updatedBook)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCatClient.AssertExpectations(t)
}

func TestUpdateBookNotFound(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	bookID := uuid.New().String()
	expectedErr := errors.New("book not found")
	updatedBook := &domain.Book{
		ID:            bookID,
		Title:         "Updated Book",
		Author:        "Updated Author",
		ISBN:          "0987654321",
		PublishedYear: 2023,
		CategoryIDs:   []string{"cat3", "cat4"},
		Stock:         20,
	}

	// Setup expectations
	mockRepo.On("GetByID", mock.Anything, bookID).Return(nil, expectedErr)

	// Execute
	err := bookUsecase.Update(context.Background(), updatedBook)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertNotCalled(t, "Update")
	mockCatClient.AssertNotCalled(t, "ValidateCategory")
}

func TestDeleteBook(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	bookID := uuid.New().String()
	now := time.Now()
	existingBook := &domain.Book{
		ID:            bookID,
		Title:         "Test Book",
		Author:        "Test Author",
		ISBN:          "1234567890",
		PublishedYear: 2022,
		CategoryIDs:   []string{"cat1", "cat2"},
		Stock:         10,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Setup expectations
	mockRepo.On("GetByID", mock.Anything, bookID).Return(existingBook, nil)
	mockRepo.On("Delete", mock.Anything, bookID).Return(nil)

	// Execute
	err := bookUsecase.Delete(context.Background(), bookID)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBookNotFound(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	bookID := uuid.New().String()
	expectedErr := errors.New("book not found")

	// Setup expectations
	mockRepo.On("GetByID", mock.Anything, bookID).Return(nil, expectedErr)

	// Execute
	err := bookUsecase.Delete(context.Background(), bookID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestListBooks(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	now := time.Now()
	expectedBooks := []*domain.Book{
		{
			ID:            uuid.New().String(),
			Title:         "Book 1",
			Author:        "Author 1",
			ISBN:          "1111111111",
			PublishedYear: 2020,
			CategoryIDs:   []string{"cat1"},
			Stock:         5,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            uuid.New().String(),
			Title:         "Book 2",
			Author:        "Author 2",
			ISBN:          "2222222222",
			PublishedYear: 2021,
			CategoryIDs:   []string{"cat2"},
			Stock:         10,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	// Setup expectations
	mockRepo.On("List", mock.Anything, 1, 10).Return(expectedBooks, 2, nil)

	// Execute
	books, total, err := bookUsecase.List(context.Background(), 1, 10)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, books, 2)
	assert.Equal(t, expectedBooks[0].Title, books[0].Title)
	assert.Equal(t, expectedBooks[1].Title, books[1].Title)
	mockRepo.AssertExpectations(t)
}

func TestSearchBooks(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	now := time.Now()
	query := "test"
	expectedBooks := []*domain.Book{
		{
			ID:            uuid.New().String(),
			Title:         "Test Book 1",
			Author:        "Test Author",
			ISBN:          "1234567890",
			PublishedYear: 2020,
			CategoryIDs:   []string{"cat1"},
			Stock:         5,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	// Setup expectations
	mockRepo.On("Search", mock.Anything, query, 1, 10).Return(expectedBooks, 1, nil)

	// Execute
	books, total, err := bookUsecase.Search(context.Background(), query, 1, 10)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, books, 1)
	assert.Equal(t, expectedBooks[0].Title, books[0].Title)
	mockRepo.AssertExpectations(t)
}

func TestValidateCategoryIDs(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	categoryIDs := []string{"cat1", "cat2"}

	// Setup expectations
	mockCatClient.On("ValidateCategory", mock.Anything, "cat1").Return(true, nil)
	mockCatClient.On("ValidateCategory", mock.Anything, "cat2").Return(true, nil)

	// Execute
	err := bookUsecase.ValidateCategoryIDs(context.Background(), categoryIDs)

	// Assertions
	assert.NoError(t, err)
	mockCatClient.AssertExpectations(t)
}

func TestValidateCategoryIDsWithInvalidCategory(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	categoryIDs := []string{"cat1", "invalidCat"}

	// Setup expectations
	mockCatClient.On("ValidateCategory", mock.Anything, "cat1").Return(true, nil)
	mockCatClient.On("ValidateCategory", mock.Anything, "invalidCat").Return(false, nil)

	// Execute
	err := bookUsecase.ValidateCategoryIDs(context.Background(), categoryIDs)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalidCat does not exist")
	mockCatClient.AssertExpectations(t)
}

func TestValidateCategoryIDsEmpty(t *testing.T) {
	// Setup mocks
	mockRepo := new(mockBookRepository)
	mockCatClient := new(mockCategoryClient)

	// Create usecase dengan mocks
	bookUsecase := usecase.NewBookUsecase(mockRepo, mockCatClient)

	// Test data
	var categoryIDs []string

	// Execute
	err := bookUsecase.ValidateCategoryIDs(context.Background(), categoryIDs)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one category must be selected")
	mockCatClient.AssertNotCalled(t, "ValidateCategory")
}