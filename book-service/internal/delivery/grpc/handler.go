package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/domain"
	pb "github.com/sibobbbbbb/backend-engineer-challenge/proto/book"
)

type BookHandler struct {
	pb.UnimplementedBookServiceServer
	bookUsecase domain.BookUsecase
}

// NewBookHandler membuat instance baru BookHandler
func NewBookHandler(bookUsecase domain.BookUsecase) *BookHandler {
	return &BookHandler{
		bookUsecase: bookUsecase,
	}
}

// CreateBook menangani permintaan pembuatan buku baru
func (h *BookHandler) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.BookResponse, error) {
	book := &domain.Book{
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.Isbn,
		PublishedYear: int(req.PublishedYear),
		CategoryIDs:   req.CategoryIds,
		Stock:         int(req.Stock),
	}

	err := h.bookUsecase.Create(ctx, book)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create book: %v", err)
	}

	return &pb.BookResponse{
		Book: convertDomainToProto(book),
	}, nil
}

// GetBook menangani permintaan untuk mendapatkan buku berdasarkan ID
func (h *BookHandler) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.BookResponse, error) {
	book, err := h.bookUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "book not found: %v", err)
	}

	return &pb.BookResponse{
		Book: convertDomainToProto(book),
	}, nil
}

// ListBooks menangani permintaan untuk daftar buku dengan pagination
func (h *BookHandler) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	books, total, err := h.bookUsecase.List(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list books: %v", err)
	}

	pbBooks := make([]*pb.Book, len(books))
	for i, book := range books {
		pbBooks[i] = convertDomainToProto(book)
	}

	return &pb.ListBooksResponse{
		Books: pbBooks,
		Total: int32(total),
	}, nil
}

// UpdateBook menangani permintaan pembaruan buku
func (h *BookHandler) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.BookResponse, error) {
	book := &domain.Book{
		ID:            req.Id,
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.Isbn,
		PublishedYear: int(req.PublishedYear),
		CategoryIDs:   req.CategoryIds,
		Stock:         int(req.Stock),
	}

	err := h.bookUsecase.Update(ctx, book)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update book: %v", err)
	}

	// Dapatkan buku yang diperbarui
	updatedBook, err := h.bookUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "book updated but failed to retrieve: %v", err)
	}

	return &pb.BookResponse{
		Book: convertDomainToProto(updatedBook),
	}, nil
}

// DeleteBook menangani permintaan penghapusan buku
func (h *BookHandler) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	err := h.bookUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete book: %v", err)
	}

	return &pb.DeleteBookResponse{
		Success: true,
	}, nil
}

// SearchBooks menangani permintaan pencarian buku
func (h *BookHandler) SearchBooks(ctx context.Context, req *pb.SearchBooksRequest) (*pb.ListBooksResponse, error) {
	books, total, err := h.bookUsecase.Search(ctx, req.Query, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search books: %v", err)
	}

	pbBooks := make([]*pb.Book, len(books))
	for i, book := range books {
		pbBooks[i] = convertDomainToProto(book)
	}

	return &pb.ListBooksResponse{
		Books: pbBooks,
		Total: int32(total),
	}, nil
}

// convertDomainToProto mengkonversi domain.Book ke pb.Book
func convertDomainToProto(book *domain.Book) *pb.Book {
	return &pb.Book{
		Id:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		Isbn:          book.ISBN,
		PublishedYear: int32(book.PublishedYear),
		CategoryIds:   book.CategoryIDs,
		Stock:         int32(book.Stock),
		CreatedAt:     book.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     book.UpdatedAt.Format(time.RFC3339),
	}
}