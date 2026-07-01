package category

import (
	"context"

	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CategoryResponse struct {
	ID        int64  `json:"id"`
	NameHi    string `json:"name_hi"`
	NameEn    string `json:"name_en"`
	ImageURL  string `json:"image_url"`
	SortOrder int32  `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

type CreateCategoryRequest struct {
	NameHi    string `json:"name_hi" binding:"required"`
	NameEn    string `json:"name_en" binding:"required"`
	ImageURL  string `json:"image_url"`
	SortOrder int32  `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	NameHi    string `json:"name_hi" binding:"required"`
	NameEn    string `json:"name_en" binding:"required"`
	ImageURL  string `json:"image_url"`
	SortOrder int32  `json:"sort_order"`
}

func toCategoryResponse(c db.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		NameHi:    c.NameHi,
		NameEn:    c.NameEn,
		ImageURL:  c.ImageUrl,
		SortOrder: c.SortOrder,
		IsActive:  c.IsActive,
	}
}

func (s *Service) List(ctx context.Context) ([]CategoryResponse, *apperrors.AppError) {
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to list categories", "श्रेणियाँ लोड करने में विफल", err)
	}
	resp := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		resp[i] = toCategoryResponse(c)
	}
	return resp, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*CategoryResponse, *apperrors.AppError) {
	c, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, apperrors.NotFoundWrap("category not found", "श्रेणी नहीं मिली", err)
	}
	resp := toCategoryResponse(c)
	return &resp, nil
}

func (s *Service) Create(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, *apperrors.AppError) {
	if req.NameHi == "" || req.NameEn == "" {
		return nil, apperrors.Validation("name is required", "नाम आवश्यक है")
	}

	c, err := s.repo.CreateCategory(ctx, db.CreateCategoryParams{
		NameHi:    req.NameHi,
		NameEn:    req.NameEn,
		ImageUrl:  req.ImageURL,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, apperrors.InternalWrap("failed to create category", "श्रेणी बनाने में विफल", err)
	}
	resp := toCategoryResponse(c)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateCategoryRequest) (*CategoryResponse, *apperrors.AppError) {
	c, err := s.repo.UpdateCategory(ctx, db.UpdateCategoryParams{
		ID:        id,
		NameHi:    req.NameHi,
		NameEn:    req.NameEn,
		ImageUrl:  req.ImageURL,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, apperrors.NotFoundWrap("category not found", "श्रेणी नहीं मिली", err)
	}
	resp := toCategoryResponse(c)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, id int64) *apperrors.AppError {
	if err := s.repo.DeleteCategory(ctx, id); err != nil {
		return apperrors.NotFoundWrap("category not found", "श्रेणी नहीं मिली", err)
	}
	return nil
}
