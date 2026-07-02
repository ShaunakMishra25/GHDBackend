package order

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gumla-hds/gumla-backend/internal/types"
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

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *Repository) GetOrderByID(ctx context.Context, id int64) (db.GetOrderByIDRow, error) {
	return r.Queries.GetOrderByID(ctx, id)
}

func (r *Repository) GetOrderItems(ctx context.Context, orderID int64) ([]db.OrderItem, error) {
	return r.Queries.GetOrderItems(ctx, orderID)
}

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID int64, limit, offset int32) ([]db.GetOrdersByUserIDRow, error) {
	return r.Queries.GetOrdersByUserID(ctx, db.GetOrdersByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *Repository) CountOrdersByUserID(ctx context.Context, userID int64) (int64, error) {
	return r.Queries.CountOrdersByUserID(ctx, userID)
}

func (r *Repository) GetAllOrders(ctx context.Context, limit, offset int32) ([]db.GetAllOrdersRow, error) {
	return r.Queries.GetAllOrders(ctx, db.GetAllOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *Repository) CountAllOrders(ctx context.Context) (int64, error) {
	return r.Queries.CountAllOrders(ctx)
}

func (r *Repository) GetOrdersByStatus(ctx context.Context, status types.OrderStatus, limit, offset int32) ([]db.GetOrdersByStatusRow, error) {
	return r.Queries.GetOrdersByStatus(ctx, db.GetOrdersByStatusParams{
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *Repository) CountOrdersByStatus(ctx context.Context, status types.OrderStatus) (int64, error) {
	return r.Queries.CountOrdersByStatus(ctx, status)
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, id int64, status types.OrderStatus) (db.Order, error) {
	return r.Queries.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     id,
		Status: status,
	})
}

func (r *Repository) GetOrdersByUserIDAndStatus(ctx context.Context, userID int64, status types.OrderStatus, limit, offset int32) ([]db.GetOrdersByUserIDAndStatusRow, error) {
	return r.Queries.GetOrdersByUserIDAndStatus(ctx, db.GetOrdersByUserIDAndStatusParams{
		UserID: userID,
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *Repository) CountOrdersByUserIDAndStatus(ctx context.Context, userID int64, status types.OrderStatus) (int64, error) {
	return r.Queries.CountOrdersByUserIDAndStatus(ctx, db.CountOrdersByUserIDAndStatusParams{
		UserID: userID,
		Status: status,
	})
}
