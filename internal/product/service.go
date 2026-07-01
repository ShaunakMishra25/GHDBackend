package product

import (
	"context"
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"

	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type ProductResponse struct {
	ID            int64   `json:"id"`
	CategoryID    int64   `json:"category_id"`
	NameHi        string  `json:"name_hi"`
	NameEn        string  `json:"name_en"`
	DescriptionHi string  `json:"description_hi"`
	DescriptionEn string  `json:"description_en"`
	Price         float64 `json:"price"`
	Unit          string  `json:"unit"`
	ImageURL      string  `json:"image_url"`
	IsActive      bool    `json:"is_active"`
	StockQty      float64 `json:"stock_qty"`
}

type CreateProductRequest struct {
	CategoryID    int64   `json:"category_id" binding:"required"`
	NameHi        string  `json:"name_hi" binding:"required"`
	NameEn        string  `json:"name_en" binding:"required"`
	DescriptionHi string  `json:"description_hi"`
	DescriptionEn string  `json:"description_en"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	Unit          string  `json:"unit" binding:"required"`
	ImageURL      string  `json:"image_url"`
	StockQty      float64 `json:"stock_qty"`
}

type UpdateProductRequest struct {
	CategoryID    int64   `json:"category_id" binding:"required"`
	NameHi        string  `json:"name_hi" binding:"required"`
	NameEn        string  `json:"name_en" binding:"required"`
	DescriptionHi string  `json:"description_hi"`
	DescriptionEn string  `json:"description_en"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	Unit          string  `json:"unit" binding:"required"`
	ImageURL      string  `json:"image_url"`
	StockQty      float64 `json:"stock_qty"`
}

type ListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Limit      int32             `json:"limit"`
	Offset     int32             `json:"offset"`
	TotalPages int               `json:"total_pages"`
}

func productToResponse(p db.Product) ProductResponse {
	var price float64
	if p.Price.Valid {
		v, _ := p.Price.Float64Value()
		price = v.Float64
	}

	var stockQty float64
	if p.StockQty.Valid {
		v, _ := p.StockQty.Float64Value()
		stockQty = v.Float64
	}

	return ProductResponse{
		ID:            p.ID,
		CategoryID:    p.CategoryID,
		NameHi:        p.NameHi,
		NameEn:        p.NameEn,
		DescriptionHi: p.DescriptionHi,
		DescriptionEn: p.DescriptionEn,
		Price:         price,
		Unit:          p.Unit,
		ImageURL:      p.ImageUrl,
		IsActive:      p.IsActive,
		StockQty:      stockQty,
	}
}

func float64ToNumeric(f float64) pgtype.Numeric {
	if f == 0 {
		return pgtype.Numeric{Valid: true, Int: big.NewInt(0), Exp: 0}
	}
	rat := new(big.Rat).SetFloat64(f)
	if rat == nil {
		return pgtype.Numeric{Valid: false}
	}
	var n pgtype.Numeric
	if err := n.Scan(rat.FloatString(2)); err != nil {
		return pgtype.Numeric{Valid: false}
	}
	return n
}

func (s *Service) List(ctx context.Context, categoryID *int64, search *string, limit, offset int32) (*ListResponse, *apperrors.AppError) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	products, err := s.repo.ListProducts(ctx, categoryID, search, limit, offset)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to list products", "उत्पाद लोड करने में विफल", err)
	}

	total, err := s.repo.CountProducts(ctx, categoryID, search)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to count products", "उत्पाद गिनती में विफल", err)
	}

	resp := make([]ProductResponse, len(products))
	for i, p := range products {
		resp[i] = productToResponse(p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &ListResponse{
		Products:   resp,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*ProductResponse, *apperrors.AppError) {
	p, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, apperrors.NotFoundWrap("product not found", "उत्पाद नहीं मिला", err)
	}
	resp := productToResponse(p)
	return &resp, nil
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*ProductResponse, *apperrors.AppError) {
	if req.Price <= 0 {
		return nil, apperrors.Validation("price must be greater than 0", "कीमत 0 से अधिक होनी चाहिए")
	}
	if req.Unit == "" {
		return nil, apperrors.Validation("unit is required", "यूनिट आवश्यक है")
	}

	price := float64ToNumeric(req.Price)
	stock := float64ToNumeric(req.StockQty)

	p, err := s.repo.CreateProduct(ctx, db.CreateProductParams{
		CategoryID:    req.CategoryID,
		NameHi:        req.NameHi,
		NameEn:        req.NameEn,
		DescriptionHi: req.DescriptionHi,
		DescriptionEn: req.DescriptionEn,
		Price:         price,
		Unit:          req.Unit,
		ImageUrl:      req.ImageURL,
		StockQty:      stock,
	})
	if err != nil {
		return nil, apperrors.InternalWrap("failed to create product", "उत्पाद बनाने में विफल", err)
	}
	resp := productToResponse(p)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateProductRequest) (*ProductResponse, *apperrors.AppError) {
	p, err := s.repo.UpdateProduct(ctx, db.UpdateProductParams{
		ID:            id,
		CategoryID:    req.CategoryID,
		NameHi:        req.NameHi,
		NameEn:        req.NameEn,
		DescriptionHi: req.DescriptionHi,
		DescriptionEn: req.DescriptionEn,
		Price:         float64ToNumeric(req.Price),
		Unit:          req.Unit,
		ImageUrl:      req.ImageURL,
		StockQty:      float64ToNumeric(req.StockQty),
	})
	if err != nil {
		return nil, apperrors.NotFoundWrap("product not found", "उत्पाद नहीं मिला", err)
	}
	resp := productToResponse(p)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, id int64) *apperrors.AppError {
	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return apperrors.NotFoundWrap("product not found", "उत्पाद नहीं मिला", err)
	}
	return nil
}
