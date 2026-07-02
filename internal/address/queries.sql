-- name: ListAddresses :many
SELECT * FROM addresses WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC;

-- name: GetAddressByID :one
SELECT * FROM addresses WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: CreateAddress :one
INSERT INTO addresses (user_id, label, full_address, landmark, latitude, longitude, is_default)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateAddress :one
UPDATE addresses
SET label = $2, full_address = $3, landmark = $4, latitude = $5, longitude = $6, is_default = $7
WHERE id = $1 AND user_id = $8
RETURNING *;

-- name: DeleteAddress :exec
DELETE FROM addresses WHERE id = $1 AND user_id = $2;

-- name: UnsetDefaultAddress :exec
UPDATE addresses SET is_default = false WHERE user_id = $1;

-- name: GetDefaultAddress :one
SELECT * FROM addresses WHERE user_id = $1 AND is_default = true LIMIT 1;
