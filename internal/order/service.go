package order

import (
	"context"
	"log"
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gumla-hds/gumla-backend/internal/address"
	"github.com/gumla-hds/gumla-backend/internal/cart"
	"github.com/gumla-hds/gumla-backend/internal/product"
	"github.com/gumla-hds/gumla-backend/internal/types"
	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Service struct {
	repo     *Repository
	cartSvc  *cart.Service
	prodRepo *product.Repository
	addrRepo *address.Repository
}

func NewService(repo *Repository, cartSvc *cart.Service, prodRepo *product.Repository, addrRepo *address.Repository) *Service {
	return &Service{
		repo:     repo,
		cartSvc:  cartSvc,
		prodRepo: prodRepo,
		addrRepo: addrRepo,
	}
}

type OrderResponse struct {
	ID             int64              `json:"id"`
	UserID         int64              `json:"user_id"`
	AddressID      int64              `json:"address_id"`
	AddressText    string             `json:"address_text"`
	Status         string             `json:"status"`
	Subtotal       float64            `json:"subtotal"`
	DeliveryCharge float64            `json:"delivery_charge"`
	Total          float64            `json:"total"`
	Notes          string             `json:"notes"`
	Items          []OrderItemResponse `json:"items,omitempty"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
}

type OrderItemResponse struct {
	ID          int64   `json:"id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    float64 `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type CreateOrderRequest struct {
	AddressID int64  `json:"address_id" binding:"required"`
	Notes     string `json:"notes"`
}

type ListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64           `json:"total"`
	Limit      int32           `json:"limit"`
	Offset     int32           `json:"offset"`
	TotalPages int             `json:"total_pages"`
}

type StatusUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	v, _ := n.Float64Value()
	return v.Float64
}

func pgNumericFromFloat64(f float64) pgtype.Numeric {
	rat := new(big.Rat).SetFloat64(f)
	if rat == nil {
		return pgtype.Numeric{Valid: false}
	}
	val := rat.FloatString(2)
	var n pgtype.Numeric
	n.Scan(val)
	return n
}

func orderRowToResponse(o db.GetOrderByIDRow, items []db.OrderItem) OrderResponse {
	itemResp := make([]OrderItemResponse, len(items))
	for i, item := range items {
		itemResp[i] = OrderItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			UnitPrice:   numericToFloat64(item.UnitPrice),
			Quantity:    numericToFloat64(item.Quantity),
			TotalPrice:  numericToFloat64(item.TotalPrice),
		}
	}
	return OrderResponse{
		ID:             o.ID,
		UserID:         o.UserID,
		AddressID:      o.AddressID,
		AddressText:    o.AddressText,
		Status:         string(o.Status),
		Subtotal:       numericToFloat64(o.Subtotal),
		DeliveryCharge: numericToFloat64(o.DeliveryCharge),
		Total:          numericToFloat64(o.Total),
		Notes:          o.Notes,
		Items:          itemResp,
		CreatedAt:      o.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      o.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func orderToSimpleResponse(o db.GetOrdersByUserIDRow) OrderResponse {
	return OrderResponse{
		ID:             o.ID,
		UserID:         o.UserID,
		AddressID:      o.AddressID,
		AddressText:    o.AddressText,
		Status:         string(o.Status),
		Subtotal:       numericToFloat64(o.Subtotal),
		DeliveryCharge: numericToFloat64(o.DeliveryCharge),
		Total:          numericToFloat64(o.Total),
		Notes:          o.Notes,
		CreatedAt:      o.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      o.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateOrderRequest) (*OrderResponse, *apperrors.AppError) {
	cartResp, appErr := s.cartSvc.GetCart(ctx, userID)
	if appErr != nil {
		return nil, appErr
	}
	if len(cartResp.Items) == 0 {
		return nil, apperrors.BadRequest("cart is empty", "कार्ट खाली है")
	}

	_, err := s.addrRepo.GetAddressByID(ctx, req.AddressID, userID)
	if err != nil {
		return nil, apperrors.BadRequest("invalid address", "गलत पता")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to begin transaction", "ऑर्डर बनाने में विफल", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.Queries.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, db.CreateOrderParams{
		UserID:         userID,
		AddressID:      req.AddressID,
		Status:         types.OrderPending,
		Subtotal:       pgNumericFromFloat64(cartResp.Subtotal),
		DeliveryCharge: pgNumericFromFloat64(cartResp.DeliveryCharge),
		Total:          pgNumericFromFloat64(cartResp.Total),
		Notes:          req.Notes,
	})
	if err != nil {
		return nil, apperrors.InternalWrap("failed to create order", "ऑर्डर बनाने में विफल", err)
	}

	for _, item := range cartResp.Items {
		if err := qtx.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:     order.ID,
			ProductID:   item.ProductID,
			ProductName: item.NameEn,
			UnitPrice:   pgNumericFromFloat64(item.Price),
			Quantity:    pgNumericFromFloat64(item.Quantity),
			TotalPrice:  pgNumericFromFloat64(item.Subtotal),
		}); err != nil {
			return nil, apperrors.InternalWrap("failed to create order item", "ऑर्डर आइटम बनाने में विफल", err)
		}
	}

	if err := qtx.ClearCart(ctx, userID); err != nil {
		log.Printf("Warning: failed to clear cart: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperrors.InternalWrap("failed to commit order", "ऑर्डर बनाने में विफल", err)
	}

	fullOrder, err := s.repo.GetOrderByID(ctx, order.ID)
	if err != nil {
		return nil, apperrors.InternalWrap("order created but failed to fetch", "ऑर्डर बन गया है पर लोड नहीं हो सका", err)
	}

	items, err := s.repo.GetOrderItems(ctx, order.ID)
	if err != nil {
		items = []db.OrderItem{}
	}

	resp := orderRowToResponse(fullOrder, items)
	return &resp, nil
}

func (s *Service) GetByID(ctx context.Context, userID, orderID int64) (*OrderResponse, *apperrors.AppError) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFoundWrap("order not found", "ऑर्डर नहीं मिला", err)
	}
	if order.UserID != userID {
		return nil, apperrors.Forbidden("not your order", "यह आपका ऑर्डर नहीं है")
	}

	items, err := s.repo.GetOrderItems(ctx, orderID)
	if err != nil {
		items = []db.OrderItem{}
	}

	resp := orderRowToResponse(order, items)
	return &resp, nil
}

func (s *Service) ListUserOrders(ctx context.Context, userID int64, limit, offset int32) (*ListResponse, *apperrors.AppError) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	orders, err := s.repo.GetOrdersByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to list orders", "ऑर्डर लोड करने में विफल", err)
	}

	total, err := s.repo.CountOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to count orders", "ऑर्डर गिनती में विफल", err)
	}

	resp := make([]OrderResponse, len(orders))
	for i, o := range orders {
		resp[i] = orderToSimpleResponse(o)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &ListResponse{Orders: resp, Total: total, Limit: limit, Offset: offset, TotalPages: totalPages}, nil
}

func (s *Service) GetByIDAdmin(ctx context.Context, orderID int64) (*OrderResponse, *apperrors.AppError) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFoundWrap("order not found", "ऑर्डर नहीं मिला", err)
	}

	items, err := s.repo.GetOrderItems(ctx, orderID)
	if err != nil {
		items = []db.OrderItem{}
	}

	resp := orderRowToResponse(order, items)
	return &resp, nil
}

func (s *Service) ListAll(ctx context.Context, status string, limit, offset int32) (*ListResponse, *apperrors.AppError) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var (
		orders []OrderResponse
		total  int64
		err    error
	)

	if status != "" {
		orderStatus := types.OrderStatus(status)
		if !orderStatus.IsValid() {
			return nil, apperrors.Validation("invalid status filter", "गलत स्थिति फ़िल्टर")
		}
		dbOrders, dbErr := s.repo.GetOrdersByStatus(ctx, orderStatus, limit, offset)
		if dbErr != nil {
			return nil, apperrors.InternalWrap("failed to list orders", "ऑर्डर लोड करने में विफल", dbErr)
		}
		total, err = s.repo.CountOrdersByStatus(ctx, orderStatus)
		err = dbErr
		orders = make([]OrderResponse, len(dbOrders))
		for i, o := range dbOrders {
			orders[i] = OrderResponse{
				ID:             o.ID,
				UserID:         o.UserID,
				AddressID:      o.AddressID,
				AddressText:    o.AddressText,
				Status:         string(o.Status),
				Subtotal:       numericToFloat64(o.Subtotal),
				DeliveryCharge: numericToFloat64(o.DeliveryCharge),
				Total:          numericToFloat64(o.Total),
				Notes:          o.Notes,
				CreatedAt:      o.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:      o.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
	} else {
		dbOrders, dbErr := s.repo.GetAllOrders(ctx, limit, offset)
		if dbErr != nil {
			return nil, apperrors.InternalWrap("failed to list orders", "ऑर्डर लोड करने में विफल", dbErr)
		}
		total, err = s.repo.CountAllOrders(ctx)
		_ = err
		orders = make([]OrderResponse, len(dbOrders))
		for i, o := range dbOrders {
			orders[i] = OrderResponse{
				ID:             o.ID,
				UserID:         o.UserID,
				AddressID:      o.AddressID,
				AddressText:    o.AddressText,
				Status:         string(o.Status),
				Subtotal:       numericToFloat64(o.Subtotal),
				DeliveryCharge: numericToFloat64(o.DeliveryCharge),
				Total:          numericToFloat64(o.Total),
				Notes:          o.Notes,
				CreatedAt:      o.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:      o.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
	}

	if err != nil {
		return nil, apperrors.InternalWrap("failed to count orders", "ऑर्डर गिनती में विफल", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &ListResponse{Orders: orders, Total: total, Limit: limit, Offset: offset, TotalPages: totalPages}, nil
}

func (s *Service) UpdateStatus(ctx context.Context, orderID int64, req StatusUpdateRequest, actorRole string) (*OrderResponse, *apperrors.AppError) {
	newStatus := types.OrderStatus(req.Status)
	if !newStatus.IsValid() {
		return nil, apperrors.Validation("invalid status", "गलत स्थिति")
	}

	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFoundWrap("order not found", "ऑर्डर नहीं मिला", err)
	}

	currentStatus := order.Status

	if actorRole != "admin" {
		if currentStatus != types.OrderPending {
			return nil, apperrors.Forbidden("cannot cancel this order", "इस ऑर्डर को रद्द नहीं कर सकते")
		}
		if newStatus != types.OrderCancelled {
			return nil, apperrors.Forbidden("only admin can update order status", "केवल एडमिन ऑर्डर की स्थिति बदल सकता है")
		}
	}

	if !currentStatus.CanTransitionTo(newStatus) {
		return nil, apperrors.Validation(
			"cannot transition from "+string(currentStatus)+" to "+string(newStatus),
			"ऑर्डर की स्थिति '"+string(currentStatus)+"' से '"+string(newStatus)+"' में नहीं बदली जा सकती",
		)
	}

	updated, err := s.repo.UpdateOrderStatus(ctx, orderID, newStatus)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to update order status", "ऑर्डर की स्थिति बदलने में विफल", err)
	}

	items, err := s.repo.GetOrderItems(ctx, orderID)
	if err != nil {
		items = []db.OrderItem{}
	}

	itemResp := make([]OrderItemResponse, len(items))
	for i, item := range items {
		itemResp[i] = OrderItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			UnitPrice:   numericToFloat64(item.UnitPrice),
			Quantity:    numericToFloat64(item.Quantity),
			TotalPrice:  numericToFloat64(item.TotalPrice),
		}
	}

	resp := OrderResponse{
		ID:             updated.ID,
		UserID:         updated.UserID,
		AddressID:      updated.AddressID,
		AddressText:    order.AddressText,
		Status:         string(updated.Status),
		Subtotal:       numericToFloat64(updated.Subtotal),
		DeliveryCharge: numericToFloat64(updated.DeliveryCharge),
		Total:          numericToFloat64(updated.Total),
		Notes:          updated.Notes,
		Items:          itemResp,
		CreatedAt:      updated.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      updated.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
	return &resp, nil
}
