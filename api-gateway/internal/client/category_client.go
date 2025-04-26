// api-gateway/internal/client/category_client.go
package client

import (
    "context"
    "log"
    "time"

    pb "github.com/sibobbbbbb/backend-engineer-challenge/proto/category"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// CategoryClient adalah interface untuk berinteraksi dengan Category Service
type CategoryClient interface {
    CreateCategory(ctx context.Context, name, description string) (*pb.Category, error)
    GetCategory(ctx context.Context, id string) (*pb.Category, error)
    ListCategories(ctx context.Context, page, limit int) ([]*pb.Category, int, error)
    UpdateCategory(ctx context.Context, id, name, description string) (*pb.Category, error)
    DeleteCategory(ctx context.Context, id string) (bool, error)
}

type categoryClient struct {
    client pb.CategoryServiceClient
}

// NewCategoryClient membuat instance baru CategoryClient
func NewCategoryClient(categoryServiceAddr string) (CategoryClient, error) {
    conn, err := grpc.Dial(
        categoryServiceAddr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
        grpc.WithTimeout(5*time.Second),
    )
    if err != nil {
        log.Printf("Failed to connect to category service: %v", err)
        return nil, err
    }

    client := pb.NewCategoryServiceClient(conn)
    return &categoryClient{client: client}, nil
}

// CreateCategory membuat kategori baru
func (c *categoryClient) CreateCategory(ctx context.Context, name, description string) (*pb.Category, error) {
    resp, err := c.client.CreateCategory(ctx, &pb.CreateCategoryRequest{
        Name:        name,
        Description: description,
    })
    if err != nil {
        return nil, err
    }
    return resp.Category, nil
}

// GetCategory mendapatkan detail kategori berdasarkan ID
func (c *categoryClient) GetCategory(ctx context.Context, id string) (*pb.Category, error) {
    resp, err := c.client.GetCategory(ctx, &pb.GetCategoryRequest{
        Id: id,
    })
    if err != nil {
        return nil, err
    }
    return resp.Category, nil
}

// ListCategories mendapatkan daftar kategori dengan pagination
func (c *categoryClient) ListCategories(ctx context.Context, page, limit int) ([]*pb.Category, int, error) {
    resp, err := c.client.ListCategories(ctx, &pb.ListCategoriesRequest{
        Page:  int32(page),
        Limit: int32(limit),
    })
    if err != nil {
        return nil, 0, err
    }
    return resp.Categories, int(resp.Total), nil
}

// UpdateCategory memperbarui kategori
func (c *categoryClient) UpdateCategory(ctx context.Context, id, name, description string) (*pb.Category, error) {
    resp, err := c.client.UpdateCategory(ctx, &pb.UpdateCategoryRequest{
        Id:          id,
        Name:        name,
        Description: description,
    })
    if err != nil {
        return nil, err
    }
    return resp.Category, nil
}

// DeleteCategory menghapus kategori
func (c *categoryClient) DeleteCategory(ctx context.Context, id string) (bool, error) {
    resp, err := c.client.DeleteCategory(ctx, &pb.DeleteCategoryRequest{
        Id: id,
    })
    if err != nil {
        return false, err
    }
    return resp.Success, nil
}