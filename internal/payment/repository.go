package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Repository struct {
	Queries *db.Queries
	pool    *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Queries: db.New(pool),
		pool:    pool,
	}
}

func (r *Repository) CreatePayment(ctx context.Context, arg db.CreatePaymentParams) (db.Payment, error) {
	return r.Queries.CreatePayment(ctx, arg)
}

func (r *Repository) GetPaymentByOrderID(ctx context.Context, orderID int64) (db.Payment, error) {
	return r.Queries.GetPaymentByOrderID(ctx, orderID)
}

func (r *Repository) GetPaymentByRazorpayOrderID(ctx context.Context, razorpayOrderID string) (db.Payment, error) {
	return r.Queries.GetPaymentByRazorpayOrderID(ctx, razorpayOrderID)
}

func (r *Repository) UpdatePaymentSuccess(ctx context.Context, razorpayOrderID, razorpayPaymentID, razorpaySignature string) (db.Payment, error) {
	return r.Queries.UpdatePaymentSuccess(ctx, db.UpdatePaymentSuccessParams{
		RazorpayOrderID:   razorpayOrderID,
		RazorpayPaymentID: razorpayPaymentID,
		RazorpaySignature: razorpaySignature,
	})
}

func (r *Repository) UpdatePaymentFailed(ctx context.Context, razorpayOrderID string) error {
	return r.Queries.UpdatePaymentFailed(ctx, razorpayOrderID)
}

func VerifySignature(orderID, paymentID, signature, secret string) bool {
	payload := orderID + "|" + paymentID
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
