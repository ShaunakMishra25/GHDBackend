-- name: GetDashboardStats :one
SELECT
  COALESCE(COUNT(*) FILTER (WHERE DATE(created_at) = CURRENT_DATE), 0)::BIGINT AS today_order_count,
  COALESCE(SUM(total) FILTER (WHERE DATE(created_at) = CURRENT_DATE), 0)::NUMERIC AS today_revenue,
  COALESCE(COUNT(*) FILTER (WHERE status = 'pending'), 0)::BIGINT AS pending_order_count,
  COALESCE(COUNT(*) FILTER (WHERE status = 'dispatched'), 0)::BIGINT AS dispatched_order_count,
  (SELECT COUNT(*) FROM users)::BIGINT AS total_users,
  (SELECT COUNT(*) FROM orders)::BIGINT AS total_orders,
  (SELECT COUNT(*) FROM products WHERE is_active = true)::BIGINT AS active_products
FROM orders;

-- name: GetTodayOrders :many
SELECT o.*, a.full_address as address_text, u.phone as user_phone, u.name as user_name
FROM orders o
JOIN addresses a ON a.id = o.address_id
JOIN users u ON u.id = o.user_id
WHERE DATE(o.created_at) = CURRENT_DATE
ORDER BY o.created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountTodayOrders :one
SELECT COUNT(*) FROM orders WHERE DATE(created_at) = CURRENT_DATE;

-- name: GetOrderStatusCounts :many
SELECT status, COUNT(*)::BIGINT as count
FROM orders
GROUP BY status;
