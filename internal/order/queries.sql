-- name: CreateOrder :one
INSERT INTO orders (user_id, address_id, status, subtotal, delivery_charge, total, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, product_id, product_name, unit_price, quantity, total_price)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetOrderByID :one
SELECT o.*, a.full_address as address_text
FROM orders o
JOIN addresses a ON a.id = o.address_id
WHERE o.id = $1 LIMIT 1;

-- name: GetOrderItems :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: GetOrdersByUserID :many
SELECT o.*, a.full_address as address_text
FROM orders o
JOIN addresses a ON a.id = o.address_id
WHERE o.user_id = $1
ORDER BY o.created_at DESC
LIMIT $2
OFFSET $3;

-- name: CountOrdersByUserID :one
SELECT COUNT(*) FROM orders WHERE user_id = $1;

-- name: GetAllOrders :many
SELECT o.*, a.full_address as address_text, u.phone as user_phone, u.name as user_name
FROM orders o
JOIN addresses a ON a.id = o.address_id
JOIN users u ON u.id = o.user_id
ORDER BY o.created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountAllOrders :one
SELECT COUNT(*) FROM orders;

-- name: GetOrdersByStatus :many
SELECT o.*, a.full_address as address_text, u.phone as user_phone, u.name as user_name
FROM orders o
JOIN addresses a ON a.id = o.address_id
JOIN users u ON u.id = o.user_id
WHERE o.status = $1
ORDER BY o.created_at DESC
LIMIT $2
OFFSET $3;

-- name: CountOrdersByStatus :one
SELECT COUNT(*) FROM orders WHERE status = $1;

-- name: UpdateOrderStatus :one
UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: GetOrdersByUserIDAndStatus :many
SELECT o.*, a.full_address as address_text
FROM orders o
JOIN addresses a ON a.id = o.address_id
WHERE o.user_id = $1 AND o.status = $2
ORDER BY o.created_at DESC
LIMIT $3
OFFSET $4;

-- name: CountOrdersByUserIDAndStatus :one
SELECT COUNT(*) FROM orders WHERE user_id = $1 AND status = $2;
