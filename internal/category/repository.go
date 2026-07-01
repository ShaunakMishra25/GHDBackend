package category

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{queries: db.New(pool)}
}

func (r *Repository) ListCategories(ctx context.Context) ([]db.Category, error) {
	return r.queries.ListCategories(ctx)
}

func (r *Repository) GetCategoryByID(ctx context.Context, id int64) (db.Category, error) {
	return r.queries.GetCategoryByID(ctx, id)
}

func (r *Repository) CreateCategory(ctx context.Context, arg db.CreateCategoryParams) (db.Category, error) {
	return r.queries.CreateCategory(ctx, arg)
}

func (r *Repository) UpdateCategory(ctx context.Context, arg db.UpdateCategoryParams) (db.Category, error) {
	return r.queries.UpdateCategory(ctx, arg)
}

func (r *Repository) DeleteCategory(ctx context.Context, id int64) error {
	return r.queries.DeleteCategory(ctx, id)
}
