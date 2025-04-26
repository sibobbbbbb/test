package grpc_client

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
	ValidateCategory(ctx context.Context, categoryID string) (bool, error)
	GetCategoryByID(ctx context.Context, categoryID string) (*pb.Category, error)
	GetMultipleCategories(ctx context.Context, categoryIDs []string) ([]*pb.Category, error)
}

type categoryCli struct {
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
	return &categoryCli{client: client}, nil
}

// ValidateCategory memeriksa apakah kategori dengan ID tertentu ada
func (c *categoryCli) ValidateCategory(ctx context.Context, categoryID string) (bool, error) {
	_, err := c.client.GetCategory(ctx, &pb.GetCategoryRequest{Id: categoryID})
	if err != nil {
		return false, nil // Kategori tidak ditemukan
	}
	return true, nil
}

// GetCategoryByID mengambil kategori dari Category Service berdasarkan ID
func (c *categoryCli) GetCategoryByID(ctx context.Context, categoryID string) (*pb.Category, error) {
	resp, err := c.client.GetCategory(ctx, &pb.GetCategoryRequest{Id: categoryID})
	if err != nil {
		return nil, err
	}
	return resp.Category, nil
}

// GetMultipleCategories mengambil beberapa kategori sekaligus berdasarkan ID
func (c *categoryCli) GetMultipleCategories(ctx context.Context, categoryIDs []string) ([]*pb.Category, error) {
	categories := make([]*pb.Category, 0, len(categoryIDs))
	
	// Ambil setiap kategori satu per satu
	// Catatan: Idealnya, kita akan mengimplementasikan metode gRPC baru di Category Service
	// untuk mendapatkan beberapa kategori sekaligus, untuk optimasi.
	for _, id := range categoryIDs {
		cat, err := c.GetCategoryByID(ctx, id)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	
	return categories, nil
}