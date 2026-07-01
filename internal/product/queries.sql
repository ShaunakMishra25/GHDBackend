-- name: ListProducts :many
SELECT p.* FROM products p
WHERE p.is_active = true
  AND (sqlc.narg('category_id')::bigint IS NULL OR p.category_id = sqlc.narg('category_id'))
  AND (sqlc.narg('search')::text IS NULL OR p.name_hi ILIKE '%' || sqlc.narg('search') || '%' OR p.name_en ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY p.created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountProducts :one
SELECT COUNT(*) FROM products p
WHERE p.is_active = true
  AND (sqlc.narg('category_id')::bigint IS NULL OR p.category_id = sqlc.narg('category_id'))
  AND (sqlc.narg('search')::text IS NULL OR p.name_hi ILIKE '%' || sqlc.narg('search') || '%' OR p.name_en ILIKE '%' || sqlc.narg('search') || '%');

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: GetProductsByIDs :many
SELECT * FROM products WHERE id = ANY(sqlc.arg('ids')::bigint[]);

-- name: CreateProduct :one
INSERT INTO products (category_id, name_hi, name_en, description_hi, description_en, price, unit, image_url, stock_qty)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET category_id = $2, name_hi = $3, name_en = $4, description_hi = $5, description_en = $6,
    price = $7, unit = $8, image_url = $9, stock_qty = $10, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
UPDATE products SET is_active = false, updated_at = NOW() WHERE id = $1;
