package cart

import (
	"context"

	"github.com/gumla-hds/gumla-backend/internal/product"
	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Service struct {
	repo     *Repository
	prodRepo *product.Repository
}

func NewService(repo *Repository, prodRepo *product.Repository) *Service {
	return &Service{repo: repo, prodRepo: prodRepo}
}

type CartItemResponse struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	NameHi    string  `json:"name_hi"`
	NameEn    string  `json:"name_en"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Unit      string  `json:"unit"`
	ImageURL  string  `json:"image_url"`
	Subtotal  float64 `json:"subtotal"`
}

type CartResponse struct {
	Items          []CartItemResponse `json:"items"`
	Subtotal       float64            `json:"subtotal"`
	DeliveryCharge float64            `json:"delivery_charge"`
	Total          float64            `json:"total"`
	ItemCount      int                `json:"item_count"`
}

type AddToCartRequest struct {
	ProductID int64   `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

func (s *Service) GetCart(ctx context.Context, userID int64) (*CartResponse, *apperrors.AppError) {
	items, err := s.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to get cart", "कार्ट लोड करने में विफल", err)
	}

	productIDs := make([]int64, 0, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := s.prodRepo.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to get products", "उत्पाद लोड करने में विफल", err)
	}

	productMap := make(map[int64]db.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	respItems := make([]CartItemResponse, 0, len(items))
	var subtotal float64

	for _, item := range items {
		p, ok := productMap[item.ProductID]
		if !ok || !p.IsActive {
			continue
		}

		qty := numericToFloat64(item.Quantity)
		var price float64
		if p.Price.Valid {
			v, _ := p.Price.Float64Value()
			price = v.Float64
		}

		lineTotal := price * qty
		subtotal += lineTotal

		respItems = append(respItems, CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			NameHi:    p.NameHi,
			NameEn:    p.NameEn,
			Price:     price,
			Quantity:  qty,
			Unit:      p.Unit,
			ImageURL:  p.ImageUrl,
			Subtotal:  lineTotal,
		})
	}

	deliveryCharge := calcDeliveryCharge(subtotal, len(respItems))

	return &CartResponse{
		Items:          respItems,
		Subtotal:       subtotal,
		DeliveryCharge: deliveryCharge,
		Total:          subtotal + deliveryCharge,
		ItemCount:      len(respItems),
	}, nil
}

func calcDeliveryCharge(subtotal float64, itemCount int) float64 {
	if subtotal <= 1000 {
		return 50
	}
	return float64(itemCount) * 10
}

func (s *Service) AddItem(ctx context.Context, userID int64, req AddToCartRequest) *apperrors.AppError {
	if req.Quantity <= 0 {
		return apperrors.Validation("quantity must be positive", "मात्रा सकारात्मक होनी चाहिए")
	}

	if err := s.repo.UpsertCartItem(ctx, userID, req.ProductID, req.Quantity); err != nil {
		return apperrors.InternalWrap("failed to add item to cart", "कार्ट में आइटम जोड़ने में विफल", err)
	}
	return nil
}

func (s *Service) UpdateItem(ctx context.Context, userID, itemID int64, req UpdateCartItemRequest) *apperrors.AppError {
	if req.Quantity <= 0 {
		return apperrors.Validation("quantity must be positive", "मात्रा सकारात्मक होनी चाहिए")
	}

	item, err := s.repo.GetCartItem(ctx, itemID, userID)
	if err != nil {
		return apperrors.NotFoundWrap("cart item not found", "कार्ट आइटम नहीं मिला", err)
	}

	if err := s.repo.UpsertCartItem(ctx, userID, item.ProductID, req.Quantity); err != nil {
		return apperrors.InternalWrap("failed to update cart item", "कार्ट आइटम अपडेट करने में विफल", err)
	}
	return nil
}

func (s *Service) RemoveItem(ctx context.Context, userID, itemID int64) *apperrors.AppError {
	if err := s.repo.RemoveCartItem(ctx, itemID, userID); err != nil {
		return apperrors.NotFoundWrap("cart item not found", "कार्ट आइटम नहीं मिला", err)
	}
	return nil
}

func (s *Service) ClearCart(ctx context.Context, userID int64) *apperrors.AppError {
	if err := s.repo.ClearCart(ctx, userID); err != nil {
		return apperrors.InternalWrap("failed to clear cart", "कार्ट खाली करने में विफल", err)
	}
	return nil
}
