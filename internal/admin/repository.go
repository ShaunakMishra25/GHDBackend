package admin

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Repository struct {
	Queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Queries: db.New(pool),
	}
}

func (r *Repository) GetDashboardStats(ctx context.Context) (db.GetDashboardStatsRow, error) {
	return r.Queries.GetDashboardStats(ctx)
}

func (r *Repository) GetTodayOrders(ctx context.Context, limit, offset int32) ([]db.GetTodayOrdersRow, error) {
	return r.Queries.GetTodayOrders(ctx, db.GetTodayOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *Repository) CountTodayOrders(ctx context.Context) (int64, error) {
	return r.Queries.CountTodayOrders(ctx)
}

func (r *Repository) GetOrderStatusCounts(ctx context.Context) ([]db.GetOrderStatusCountsRow, error) {
	return r.Queries.GetOrderStatusCounts(ctx)
}
