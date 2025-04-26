package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/domain"
	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/pkg/grpc_client"
)

type bookUsecase struct {
	bookRepo       domain.BookRepository
	categoryClient grpc_client.CategoryClient
}

// NewBookUsecase menciptakan instance baru dari BookUsecase
func NewBookUsecase(bookRepo domain.BookRepository, categoryClient grpc_client.CategoryClient) domain.BookUsecase {
	return &bookUsecase{
		bookRepo:       bookRepo,
		categoryClient: categoryClient,
	}
}

// Create membuat buku baru setelah memvalidasi kategori buku
func (u *bookUsecase) Create(ctx context.Context, book *domain.Book) error {
	// Validasi kategori apakah ada
	if err := u.ValidateCategoryIDs(ctx, book.CategoryIDs); err != nil {
		return err
	}

	// Validasi data buku
	if book.Title == "" {
		return errors.New("title is required")
	}
	if book.Author == "" {
		return errors.New("author is required")
	}
	if book.ISBN == "" {
		return errors.New("ISBN is required")
	}
	if book.PublishedYear <= 0 {
		return errors.New("published year must be a positive number")
	}
	if book.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	return u.bookRepo.Create(ctx, book)
}

// GetByID mengambil buku berdasarkan ID
func (u *bookUsecase) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	return u.bookRepo.GetByID(ctx, id)
}

// Update memperbarui buku yang sudah ada
func (u *bookUsecase) Update(ctx context.Context, book *domain.Book) error {
	// Periksa apakah buku ada
	existingBook, err := u.bookRepo.GetByID(ctx, book.ID)
	if err != nil {
		return err
	}
	if existingBook == nil {
		return errors.New("book not found")
	}

	// Validasi kategori apakah ada
	if err := u.ValidateCategoryIDs(ctx, book.CategoryIDs); err != nil {
		return err
	}

	// Validasi data buku
	if book.Title == "" {
		return errors.New("title is required")
	}
	if book.Author == "" {
		return errors.New("author is required")
	}
	if book.ISBN == "" {
		return errors.New("ISBN is required")
	}
	if book.PublishedYear <= 0 {
		return errors.New("published year must be a positive number")
	}
	if book.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	return u.bookRepo.Update(ctx, book)
}

// Delete menghapus buku berdasarkan ID
func (u *bookUsecase) Delete(ctx context.Context, id string) error {
	// Periksa apakah buku ada
	existingBook, err := u.bookRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingBook == nil {
		return errors.New("book not found")
	}

	return u.bookRepo.Delete(ctx, id)
}

// List mengambil daftar buku dengan pagination
func (u *bookUsecase) List(ctx context.Context, page, limit int) ([]*domain.Book, int, error) {
	return u.bookRepo.List(ctx, page, limit)
}

// Search mencari buku berdasarkan query
func (u *bookUsecase) Search(ctx context.Context, query string, page, limit int) ([]*domain.Book, int, error) {
	return u.bookRepo.Search(ctx, query, page, limit)
}

// ValidateCategoryIDs memvalidasi apakah semua ID kategori ada di category service
func (u *bookUsecase) ValidateCategoryIDs(ctx context.Context, categoryIDs []string) error {
	if len(categoryIDs) == 0 {
		return errors.New("at least one category must be selected")
	}

	// Validasi setiap ID kategori
	for _, categoryID := range categoryIDs {
		exists, err := u.categoryClient.ValidateCategory(ctx, categoryID)
		if err != nil {
			return fmt.Errorf("failed to validate category: %w", err)
		}
		if !exists {
			return fmt.Errorf("category with ID %s does not exist", categoryID)
		}
	}

	return nil
}