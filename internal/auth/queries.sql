-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (phone, name, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpsertDevice :exec
INSERT INTO devices (user_id, fcm_token, device_info)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, fcm_token) 
DO UPDATE SET device_info = $3, last_active_at = NOW();

-- name: GetUserDevices :many
SELECT * FROM devices WHERE user_id = $1;

-- name: RevokeToken :exec
INSERT INTO token_blacklist (token_jti, expires_at)
VALUES ($1, $2)
ON CONFLICT (token_jti) DO NOTHING;

-- name: IsTokenBlacklisted :one
SELECT EXISTS(
    SELECT 1 FROM token_blacklist 
    WHERE token_jti = $1 AND expires_at > NOW()
) AS is_blacklisted;

-- name: UpdateUser :exec
UPDATE users SET name = $2, updated_at = NOW() WHERE id = $1;
