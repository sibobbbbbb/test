package usecase

import (
    "context"
    "errors"
    "time"

    "github.com/sibobbbbbb/backend-engineer-challenge/category-service/internal/domain"
)

type categoryUsecase struct {
    categoryRepo domain.CategoryRepository
    contextTimeout time.Duration
    cacheExp       time.Duration
}

func NewCategoryUsecase(cr domain.CategoryRepository, timeout, cacheExp time.Duration) domain.CategoryUsecase {
    return &categoryUsecase{
        categoryRepo:   cr,
        contextTimeout: timeout,
        cacheExp:       cacheExp,
    }
}

func (u *categoryUsecase) Create(ctx context.Context, cat *domain.Category) error {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    // Validasi
    if cat.Name == "" {
        return errors.New("name is required")
    }

    // Simpan ke database
    if err := u.categoryRepo.Create(ctx, cat); err != nil {
        return err
    }

    return nil
}

func (u *categoryUsecase) GetByID(ctx context.Context, id string) (*domain.Category, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    return u.categoryRepo.GetByID(ctx, id)
}

func (u *categoryUsecase) Update(ctx context.Context, cat *domain.Category) error {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    // Periksa apakah kategori ada
    existingCat, err := u.categoryRepo.GetByID(ctx, cat.ID)
    if err != nil {
        return err
    }
    if existingCat == nil {
        return errors.New("category not found")
    }

    // Validasi
    if cat.Name == "" {
        return errors.New("name is required")
    }

    // Update
    return u.categoryRepo.Update(ctx, cat)
}

func (u *categoryUsecase) Delete(ctx context.Context, id string) error {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    // Periksa apakah kategori ada
    existingCat, err := u.categoryRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if existingCat == nil {
        return errors.New("category not found")
    }

    // Delete
    return u.categoryRepo.Delete(ctx, id)
}

func (u *categoryUsecase) List(ctx context.Context, page, limit int) ([]*domain.Category, int, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    return u.categoryRepo.List(ctx, page, limit)
}

func (u *categoryUsecase) GetMultipleByIDs(ctx context.Context, ids []string) ([]*domain.Category, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    if len(ids) == 0 {
        return []*domain.Category{}, nil
    }

    return u.categoryRepo.GetMultipleByIDs(ctx, ids)
}

func (u *categoryUsecase) ValidateCategoryExists(ctx context.Context, id string) (bool, error) {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    return u.categoryRepo.IsExistsByID(ctx, id)
}