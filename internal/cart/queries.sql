-- name: GetCartByUserID :many
SELECT ci.* FROM cart_items ci
WHERE ci.user_id = $1
ORDER BY ci.created_at ASC;

-- name: GetCartItem :one
SELECT * FROM cart_items WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetCartItemByProduct :one
SELECT * FROM cart_items WHERE user_id = $1 AND product_id = $2 LIMIT 1;

-- name: UpsertCartItem :exec
INSERT INTO cart_items (user_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, product_id)
DO UPDATE SET quantity = EXCLUDED.quantity, updated_at = NOW();

-- name: RemoveCartItem :exec
DELETE FROM cart_items WHERE id = $1 AND user_id = $2;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE user_id = $1;

-- name: GetCartCountByUserID :one
SELECT COUNT(*) FROM cart_items WHERE user_id = $1;
