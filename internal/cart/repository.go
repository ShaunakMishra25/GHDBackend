package cart

import (
	"context"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{queries: db.New(pool)}
}

func (r *Repository) GetCartByUserID(ctx context.Context, userID int64) ([]db.CartItem, error) {
	return r.queries.GetCartByUserID(ctx, userID)
}

func (r *Repository) GetCartItem(ctx context.Context, itemID, userID int64) (db.CartItem, error) {
	return r.queries.GetCartItem(ctx, db.GetCartItemParams{
		ID:     itemID,
		UserID: userID,
	})
}

func (r *Repository) UpsertCartItem(ctx context.Context, userID, productID int64, quantity float64) error {
	return r.queries.UpsertCartItem(ctx, db.UpsertCartItemParams{
		UserID:    userID,
		ProductID: productID,
		Quantity:  pgNumericFromFloat64(quantity),
	})
}

func (r *Repository) RemoveCartItem(ctx context.Context, id, userID int64) error {
	return r.queries.RemoveCartItem(ctx, db.RemoveCartItemParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *Repository) ClearCart(ctx context.Context, userID int64) error {
	return r.queries.ClearCart(ctx, userID)
}

func (r *Repository) GetCartCount(ctx context.Context, userID int64) (int64, error) {
	return r.queries.GetCartCountByUserID(ctx, userID)
}

func pgNumericFromFloat64(f float64) pgtype.Numeric {
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

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	v, _ := n.Float64Value()
	return v.Float64
}
