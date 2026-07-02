package payment

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gumla-hds/gumla-backend/internal/order"
	"github.com/gumla-hds/gumla-backend/internal/types"
	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	"github.com/gumla-hds/gumla-backend/pkg/razorpay"
	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

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

type Service struct {
	repo         *Repository
	rp           *razorpay.Client
	orderService *order.Service
	orderRepo    *order.Repository
}

func NewService(repo *Repository, rp *razorpay.Client, os *order.Service, or *order.Repository) *Service {
	return &Service{
		repo:         repo,
		rp:           rp,
		orderService: os,
		orderRepo:    or,
	}
}

type InitiatePaymentRequest struct {
	OrderID int64 `json:"order_id" binding:"required"`
}

type InitiatePaymentResponse struct {
	RazorpayOrderID string  `json:"razorpay_order_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Key             string  `json:"key"`
}

type VerifyPaymentRequest struct {
	RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
	RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
	RazorpaySignature string `json:"razorpay_signature" binding:"required"`
}

func (s *Service) Initiate(ctx context.Context, userID int64, req InitiatePaymentRequest) (*InitiatePaymentResponse, *apperrors.AppError) {
	orderResp, appErr := s.orderService.GetByID(ctx, userID, req.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	if orderResp.Status != string(types.OrderPending) {
		return nil, apperrors.BadRequest("order is not pending", "ऑर्डर पेंडिंग नहीं है")
	}

	existing, err := s.repo.GetPaymentByOrderID(ctx, req.OrderID)
	if err == nil && existing.ID > 0 {
		return &InitiatePaymentResponse{
			RazorpayOrderID: existing.RazorpayOrderID,
		Amount:          orderResp.Total,
			Currency:        "INR",
			Key:             s.rp.KeyID,
		}, nil
	}

	amountPaise := int64(orderResp.Total * 100)
	rpOrder, err := s.rp.CreateOrder(razorpay.OrderRequest{
		Amount:   amountPaise,
		Currency: "INR",
		Receipt:  fmt.Sprintf("order_%d", req.OrderID),
	})
	if err != nil {
		return nil, apperrors.InternalWrap("failed to create razorpay order", "भुगतान आरंभ करने में विफल", err)
	}

	_, err = s.repo.CreatePayment(ctx, db.CreatePaymentParams{
		OrderID:         req.OrderID,
		RazorpayOrderID: rpOrder.ID,
		Amount:          pgNumericFromFloat64(orderResp.Total),
		Currency:        "INR",
		Status:          "created",
	})
	if err != nil {
		return nil, apperrors.InternalWrap("failed to save payment", "भुगतान सहेजने में विफल", err)
	}

	return &InitiatePaymentResponse{
		RazorpayOrderID: rpOrder.ID,
		Amount:          orderResp.Total,
		Currency:        "INR",
		Key:             s.rp.KeyID,
	}, nil
}

func (s *Service) Verify(ctx context.Context, userID int64, req VerifyPaymentRequest) *apperrors.AppError {
	payment, err := s.repo.GetPaymentByRazorpayOrderID(ctx, req.RazorpayOrderID)
	if err != nil {
		return apperrors.NotFoundWrap("payment not found", "भुगतान नहीं मिला", err)
	}

	if payment.Status == "captured" {
		return nil
	}

	valid := s.rp.VerifyPaymentSignature(req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature)
	if !valid {
		s.repo.UpdatePaymentFailed(ctx, req.RazorpayOrderID)
		return apperrors.Validation("invalid payment signature", "गलत भुगतान हस्ताक्षर")
	}

	_, err = s.repo.UpdatePaymentSuccess(ctx, req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature)
	if err != nil {
		return apperrors.InternalWrap("failed to update payment", "भुगतान अपडेट करने में विफल", err)
	}

	_, dbErr := s.orderRepo.UpdateOrderStatus(ctx, payment.OrderID, types.OrderConfirmed)
	if dbErr != nil {
		log.Printf("Warning: failed to update order status after payment: %v", dbErr)
	}

	return nil
}

func (s *Service) VerifyWebhookSignature(body []byte, signature string) bool {
	return s.rp.VerifyWebhookSignature(body, signature)
}

func (s *Service) HandleWebhook(ctx context.Context, razorpayOrderID string) *apperrors.AppError {
	payment, err := s.repo.GetPaymentByRazorpayOrderID(ctx, razorpayOrderID)
	if err != nil {
		return apperrors.NotFoundWrap("payment not found for webhook", "भुगतान नहीं मिला", err)
	}

	if payment.Status == "captured" {
		return nil
	}

	if payment.Status == "failed" {
		return nil
	}

	_, err = s.repo.UpdatePaymentSuccess(ctx, razorpayOrderID, "", "")
	if err != nil {
		return apperrors.InternalWrap("failed to update payment via webhook", "भुगतान अपडेट करने में विफल", err)
	}

	_, dbErr := s.orderRepo.UpdateOrderStatus(ctx, payment.OrderID, types.OrderConfirmed)
	if dbErr != nil {
		log.Printf("Warning: failed to update order status via webhook: %v", dbErr)
	}

	return nil
}
