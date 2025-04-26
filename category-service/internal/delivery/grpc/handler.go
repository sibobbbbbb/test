package grpc

import (
    "context"
    "time"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    "github.com/sibobbbbbb/backend-engineer-challenge/category-service/internal/domain"
    pb "github.com/sibobbbbbb/backend-engineer-challenge/proto/category"
)

type CategoryHandler struct {
    pb.UnimplementedCategoryServiceServer
    categoryUsecase domain.CategoryUsecase
}

func NewCategoryHandler(cu domain.CategoryUsecase) *CategoryHandler {
    return &CategoryHandler{
        categoryUsecase: cu,
    }
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
    category := &domain.Category{
        Name:        req.Name,
        Description: req.Description,
    }

    err := h.categoryUsecase.Create(ctx, category)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
    }

    return &pb.CategoryResponse{
        Category: convertDomainToProto(category),
    }, nil
}

func (h *CategoryHandler) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.CategoryResponse, error) {
    category, err := h.categoryUsecase.GetByID(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "category not found: %v", err)
    }

    return &pb.CategoryResponse{
        Category: convertDomainToProto(category),
    }, nil
}

func (h *CategoryHandler) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
    categories, total, err := h.categoryUsecase.List(ctx, int(req.Page), int(req.Limit))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to list categories: %v", err)
    }

    pbCategories := make([]*pb.Category, len(categories))
    for i, category := range categories {
        pbCategories[i] = convertDomainToProto(category)
    }

    return &pb.ListCategoriesResponse{
        Categories: pbCategories,
        Total:      int32(total),
    }, nil
}

func (h *CategoryHandler) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error) {
    category := &domain.Category{
        ID:          req.Id,
        Name:        req.Name,
        Description: req.Description,
    }

    err := h.categoryUsecase.Update(ctx, category)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to update category: %v", err)
    }

    // Get updated category
    updatedCategory, err := h.categoryUsecase.GetByID(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "category updated but failed to retrieve: %v", err)
    }

    return &pb.CategoryResponse{
        Category: convertDomainToProto(updatedCategory),
    }, nil
}

func (h *CategoryHandler) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.DeleteCategoryResponse, error) {
    err := h.categoryUsecase.Delete(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to delete category: %v", err)
    }

    return &pb.DeleteCategoryResponse{
        Success: true,
    }, nil
}

func convertDomainToProto(category *domain.Category) *pb.Category {
    return &pb.Category{
        Id:          category.ID,
        Name:        category.Name,
        Description: category.Description,
        CreatedAt:   category.CreatedAt.Format(time.RFC3339),
        UpdatedAt:   category.UpdatedAt.Format(time.RFC3339),
    }
}