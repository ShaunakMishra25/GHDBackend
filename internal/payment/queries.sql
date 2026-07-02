-- name: CreatePayment :one
INSERT INTO payments (order_id, razorpay_order_id, amount, currency, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPaymentByOrderID :one
SELECT * FROM payments WHERE order_id = $1 LIMIT 1;

-- name: GetPaymentByRazorpayOrderID :one
SELECT * FROM payments WHERE razorpay_order_id = $1 LIMIT 1;

-- name: UpdatePaymentSuccess :one
UPDATE payments
SET razorpay_payment_id = $2, razorpay_signature = $3, status = 'captured', updated_at = NOW()
WHERE razorpay_order_id = $1
RETURNING *;

-- name: UpdatePaymentFailed :exec
UPDATE payments
SET status = 'failed', updated_at = NOW()
WHERE razorpay_order_id = $1;
