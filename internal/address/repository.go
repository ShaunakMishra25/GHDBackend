package address

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

func (r *Repository) ListAddresses(ctx context.Context, userID int64) ([]db.Address, error) {
	return r.queries.ListAddresses(ctx, userID)
}

func (r *Repository) GetAddressByID(ctx context.Context, id, userID int64) (db.Address, error) {
	return r.queries.GetAddressByID(ctx, db.GetAddressByIDParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *Repository) CreateAddress(ctx context.Context, arg db.CreateAddressParams) (db.Address, error) {
	return r.queries.CreateAddress(ctx, arg)
}

func (r *Repository) UpdateAddress(ctx context.Context, arg db.UpdateAddressParams) (db.Address, error) {
	return r.queries.UpdateAddress(ctx, arg)
}

func (r *Repository) DeleteAddress(ctx context.Context, id, userID int64) error {
	return r.queries.DeleteAddress(ctx, db.DeleteAddressParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *Repository) UnsetDefaultAddress(ctx context.Context, userID int64) error {
	return r.queries.UnsetDefaultAddress(ctx, userID)
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
