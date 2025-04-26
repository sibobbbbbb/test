package client

import (
    "context"
    "log"
    "time"

    pb "github.com/sibobbbbbb/backend-engineer-challenge/proto/book"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// BookClient adalah interface untuk berinteraksi dengan Book Service
type BookClient interface {
    CreateBook(ctx context.Context, title, author, isbn string, publishedYear int, categoryIDs []string, stock int) (*pb.Book, error)
    GetBook(ctx context.Context, id string) (*pb.Book, error)
    ListBooks(ctx context.Context, page, limit int) ([]*pb.Book, int, error)
    UpdateBook(ctx context.Context, id, title, author, isbn string, publishedYear int, categoryIDs []string, stock int) (*pb.Book, error)
    DeleteBook(ctx context.Context, id string) (bool, error)
    SearchBooks(ctx context.Context, query string, page, limit int) ([]*pb.Book, int, error)
}

type bookClient struct {
    client pb.BookServiceClient
}

// NewBookClient membuat instance baru BookClient
func NewBookClient(bookServiceAddr string) (BookClient, error) {
    conn, err := grpc.Dial(
        bookServiceAddr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
        grpc.WithTimeout(5*time.Second),
    )
    if err != nil {
        log.Printf("Failed to connect to book service: %v", err)
        return nil, err
    }

    client := pb.NewBookServiceClient(conn)
    return &bookClient{client: client}, nil
}

// CreateBook membuat buku baru
func (c *bookClient) CreateBook(ctx context.Context, title, author, isbn string, publishedYear int, categoryIDs []string, stock int) (*pb.Book, error) {
    resp, err := c.client.CreateBook(ctx, &pb.CreateBookRequest{
        Title:         title,
        Author:        author,
        Isbn:          isbn,
        PublishedYear: int32(publishedYear),
        CategoryIds:   categoryIDs,
        Stock:         int32(stock),
    })
    if err != nil {
        return nil, err
    }
    return resp.Book, nil
}

// GetBook mendapatkan detail buku berdasarkan ID
func (c *bookClient) GetBook(ctx context.Context, id string) (*pb.Book, error) {
    resp, err := c.client.GetBook(ctx, &pb.GetBookRequest{
        Id: id,
    })
    if err != nil {
        return nil, err
    }
    return resp.Book, nil
}

// ListBooks mendapatkan daftar buku dengan pagination
func (c *bookClient) ListBooks(ctx context.Context, page, limit int) ([]*pb.Book, int, error) {
    resp, err := c.client.ListBooks(ctx, &pb.ListBooksRequest{
        Page:  int32(page),
        Limit: int32(limit),
    })
    if err != nil {
        return nil, 0, err
    }
    return resp.Books, int(resp.Total), nil
}

// UpdateBook memperbarui buku
func (c *bookClient) UpdateBook(ctx context.Context, id, title, author, isbn string, publishedYear int, categoryIDs []string, stock int) (*pb.Book, error) {
    resp, err := c.client.UpdateBook(ctx, &pb.UpdateBookRequest{
        Id:            id,
        Title:         title,
        Author:        author,
        Isbn:          isbn,
        PublishedYear: int32(publishedYear),
        CategoryIds:   categoryIDs,
        Stock:         int32(stock),
    })
    if err != nil {
        return nil, err
    }
    return resp.Book, nil
}

// DeleteBook menghapus buku
func (c *bookClient) DeleteBook(ctx context.Context, id string) (bool, error) {
    resp, err := c.client.DeleteBook(ctx, &pb.DeleteBookRequest{
        Id: id,
    })
    if err != nil {
        return false, err
    }
    return resp.Success, nil
}

// SearchBooks mencari buku berdasarkan query
func (c *bookClient) SearchBooks(ctx context.Context, query string, page, limit int) ([]*pb.Book, int, error) {
    resp, err := c.client.SearchBooks(ctx, &pb.SearchBooksRequest{
        Query: query,
        Page:  int32(page),
        Limit: int32(limit),
    })
    if err != nil {
        return nil, 0, err
    }
    return resp.Books, int(resp.Total), nil
}