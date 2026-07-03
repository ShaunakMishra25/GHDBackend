-- name: CreateNotificationLog :one
INSERT INTO notifications_log (user_id, order_id, title, body, status, fcm_error, retry_count)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateNotificationStatus :exec
UPDATE notifications_log
SET status = $2, fcm_error = $3, retry_count = $4
WHERE id = $1;

-- name: GetUserNotifications :many
SELECT * FROM notifications_log
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: CountUserNotifications :one
SELECT COUNT(*) FROM notifications_log WHERE user_id = $1;
