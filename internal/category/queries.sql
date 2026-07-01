-- name: ListCategories :many
SELECT * FROM categories WHERE is_active = true ORDER BY sort_order ASC, id ASC;

-- name: ListAllCategories :many
SELECT * FROM categories ORDER BY sort_order ASC, id ASC;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1 LIMIT 1;

-- name: CreateCategory :one
INSERT INTO categories (name_hi, name_en, image_url, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET name_hi = $2, name_en = $3, image_url = $4, sort_order = $5
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
UPDATE categories SET is_active = false WHERE id = $1;
