package product

import (
	"context"

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

func (r *Repository) ListProducts(ctx context.Context, categoryID *int64, search *string, limit, offset int32) ([]db.Product, error) {
	params := db.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	}
	if categoryID != nil {
		params.CategoryID = pgtype.Int8{Int64: *categoryID, Valid: true}
	}
	if search != nil {
		params.Search = pgtype.Text{String: *search, Valid: true}
	}
	return r.queries.ListProducts(ctx, params)
}

func (r *Repository) CountProducts(ctx context.Context, categoryID *int64, search *string) (int64, error) {
	params := db.CountProductsParams{}
	if categoryID != nil {
		params.CategoryID = pgtype.Int8{Int64: *categoryID, Valid: true}
	}
	if search != nil {
		params.Search = pgtype.Text{String: *search, Valid: true}
	}
	return r.queries.CountProducts(ctx, params)
}

func (r *Repository) GetProductByID(ctx context.Context, id int64) (db.Product, error) {
	return r.queries.GetProductByID(ctx, id)
}

func (r *Repository) GetProductsByIDs(ctx context.Context, ids []int64) ([]db.Product, error) {
	return r.queries.GetProductsByIDs(ctx, ids)
}

func (r *Repository) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
	return r.queries.CreateProduct(ctx, arg)
}

func (r *Repository) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
	return r.queries.UpdateProduct(ctx, arg)
}

func (r *Repository) DeleteProduct(ctx context.Context, id int64) error {
	return r.queries.DeleteProduct(ctx, id)
}
