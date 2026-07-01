package firebase

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/gumla-hds/gumla-backend/pkg/retry"
)

type NotificationPayload struct {
	Title   string
	Body    string
	Data    map[string]string
}

func (c *Client) SendToDevice(ctx context.Context, fcmToken string, payload NotificationPayload) error {
	message := &messaging.Message{
		Token: fcmToken,
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Body,
		},
		Data: payload.Data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: "order_updates",
				Priority:  messaging.PriorityHigh,
				Sound:     "default",
			},
		},
	}

	cfg := retry.DefaultConfig()
	cfg.MaxAttempts = 3

	return retry.Do(ctx, cfg, func() error {
		_, err := c.Messaging.Send(ctx, message)
		if err != nil {
			return fmt.Errorf("send fcm: %w", err)
		}
		return nil
	})
}

func (c *Client) SendToMultipleDevices(ctx context.Context, fcmTokens []string, payload NotificationPayload) (successCount int, errors []error) {
	successCount = 0

	for _, token := range fcmTokens {
		if err := c.SendToDevice(ctx, token, payload); err != nil {
			errors = append(errors, fmt.Errorf("token %s: %w", token, err))
			continue
		}
		successCount++
	}

	return successCount, errors
}

func (c *Client) SendWithTimeout(ctx context.Context, fcmToken string, payload NotificationPayload, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.SendToDevice(ctx, fcmToken, payload)
}
