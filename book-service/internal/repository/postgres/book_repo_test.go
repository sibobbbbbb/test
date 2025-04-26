package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/domain"
	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/repository/postgres"
)

func TestCreateBook(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Buat repository dengan mock db
	repo := postgres.NewBookRepository(sqlxDB)

	// Test data
	bookID := uuid.New().String()
	now := time.Now()
	book := &domain.Book{
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
	mock.ExpectExec("INSERT INTO books").
		WithArgs(
			book.ID,
			book.Title,
			book.Author,
			book.ISBN,
			book.PublishedYear,
			pq.Array(book.CategoryIDs),
			book.Stock,
			book.CreatedAt,
			book.UpdatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.Create(context.Background(), book)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookByID(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Buat repository dengan mock db
	repo := postgres.NewBookRepository(sqlxDB)

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
	rows := sqlmock.NewRows([]string{"id", "title", "author", "isbn", "published_year", "category_ids", "stock", "created_at", "updated_at"}).
		AddRow(expectedBook.ID, expectedBook.Title, expectedBook.Author, expectedBook.ISBN, expectedBook.PublishedYear, pq.Array(expectedBook.CategoryIDs), expectedBook.Stock, expectedBook.CreatedAt, expectedBook.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM books WHERE id = (.+)").
		WithArgs(bookID).
		WillReturnRows(rows)

	// Execute
	book, err := repo.GetByID(context.Background(), bookID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedBook.ID, book.ID)
	assert.Equal(t, expectedBook.Title, book.Title)
	assert.Equal(t, expectedBook.Author, book.Author)
	assert.Equal(t, expectedBook.ISBN, book.ISBN)
	assert.Equal(t, expectedBook.PublishedYear, book.PublishedYear)
	assert.Equal(t, expectedBook.Stock, book.Stock)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBook(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Buat repository dengan mock db
	repo := postgres.NewBookRepository(sqlxDB)

	// Test data
	bookID := uuid.New().String()
	now := time.Now()
	book := &domain.Book{
		ID:            bookID,
		Title:         "Updated Book",
		Author:        "Updated Author",
		ISBN:          "0987654321",
		PublishedYear: 2023,
		CategoryIDs:   []string{"cat3", "cat4"},
		Stock:         20,
		UpdatedAt:     now,
	}

	// Setup expectations
	mock.ExpectExec("UPDATE books SET").
		WithArgs(
			book.Title,
			book.Author,
			book.ISBN,
			book.PublishedYear,
			pq.Array(book.CategoryIDs),
			book.Stock,
			book.UpdatedAt,
			book.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.Update(context.Background(), book)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBook(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Buat repository dengan mock db
	repo := postgres.NewBookRepository(sqlxDB)

	// Test data
	bookID := uuid.New().String()

	// Setup expectations
	mock.ExpectExec("DELETE FROM books WHERE id = (.+)").
		WithArgs(bookID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.Delete(context.Background(), bookID)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListBooks(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Buat repository dengan mock db
	repo := postgres.NewBookRepository(sqlxDB)

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

	// Setup expectations for count
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM books").
		WillReturnRows(countRows)

	// Setup expectations for list
	rows := sqlmock.NewRows([]string{"id", "title", "author", "isbn", "published_year", "category_ids", "stock", "created_at", "updated_at"})
	for _, book := range expectedBooks {
		rows.AddRow(book.ID, book.Title, book.Author, book.ISBN, book.PublishedYear, pq.Array(book.CategoryIDs), book.Stock, book.CreatedAt, book.UpdatedAt)
	}

	mock.ExpectQuery("SELECT (.+) FROM books ORDER BY created_at DESC LIMIT (.+) OFFSET (.+)").
		WithArgs(10, 0).
		WillReturnRows(rows)

	// Execute
	books, total, err := repo.List(context.Background(), 1, 10)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, books, 2)
	assert.Equal(t, expectedBooks[0].Title, books[0].Title)
	assert.Equal(t, expectedBooks[1].Title, books[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchBooks(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Buat repository dengan mock db
	repo := postgres.NewBookRepository(sqlxDB)

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

	// Setup expectations for count
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM books WHERE title ILIKE (.+) OR author ILIKE (.+) OR isbn ILIKE (.+)").
		WithArgs("%"+query+"%", "%"+query+"%", "%"+query+"%").
		WillReturnRows(countRows)

	// Setup expectations for search
	rows := sqlmock.NewRows([]string{"id", "title", "author", "isbn", "published_year", "category_ids", "stock", "created_at", "updated_at"})
	for _, book := range expectedBooks {
		rows.AddRow(book.ID, book.Title, book.Author, book.ISBN, book.PublishedYear, pq.Array(book.CategoryIDs), book.Stock, book.CreatedAt, book.UpdatedAt)
	}

	mock.ExpectQuery("SELECT (.+) FROM books WHERE title ILIKE (.+) OR author ILIKE (.+) OR isbn ILIKE (.+) ORDER BY created_at DESC LIMIT (.+) OFFSET (.+)").
		WithArgs("%"+query+"%", "%"+query+"%", "%"+query+"%", 10, 0).
		WillReturnRows(rows)

	// Execute
	books, total, err := repo.Search(context.Background(), query, 1, 10)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, books, 1)
	assert.Equal(t, expectedBooks[0].Title, books[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}