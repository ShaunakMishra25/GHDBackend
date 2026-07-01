package razorpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	razorpay "github.com/razorpay/razorpay-go"
)

type Client struct {
	client        *razorpay.Client
	keySecret     string
	webhookSecret string
}

func NewClient(keyID, keySecret, webhookSecret string) *Client {
	return &Client{
		client:        razorpay.NewClient(keyID, keySecret),
		keySecret:     keySecret,
		webhookSecret: webhookSecret,
	}
}

type OrderRequest struct {
	Amount   int64  // in paise
	Currency string // INR
	Receipt  string
}

type OrderResponse struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
	Status string `json:"status"`
}

func (c *Client) CreateOrder(req OrderRequest) (*OrderResponse, error) {
	data := map[string]interface{}{
		"amount":   req.Amount,
		"currency": req.Currency,
		"receipt":  req.Receipt,
	}

	body, err := c.client.Order.Create(data, nil)
	if err != nil {
		return nil, fmt.Errorf("create razorpay order: %w", err)
	}

	id, _ := body["id"].(string)
	amount, _ := body["amount"].(float64)
	status, _ := body["status"].(string)

	return &OrderResponse{
		ID:     id,
		Amount: int64(amount),
		Status: status,
	}, nil
}

func (c *Client) VerifyPaymentSignature(orderID, paymentID, signature string) bool {
	payload := orderID + "|" + paymentID
	expected := computeHMACSHA256(payload, c.keySecret)
	return hmac.Equal([]byte(expected), []byte(signature))
}

func (c *Client) VerifyWebhookSignature(payload []byte, signature string) bool {
	expected := computeHMACSHA256(string(payload), c.webhookSecret)
	return hmac.Equal([]byte(expected), []byte(signature))
}

func (c *Client) GetPaymentByID(paymentID string) (map[string]interface{}, error) {
	body, err := c.client.Payment.Fetch(paymentID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch payment %s: %w", paymentID, err)
	}
	return body, nil
}

func computeHMACSHA256(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
