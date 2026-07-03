package notification

import (
	"context"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gumla-hds/gumla-backend/internal/auth"
	"github.com/gumla-hds/gumla-backend/pkg/firebase"
	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Service struct {
	repo     *Repository
	firebase *firebase.Client
	authRepo *auth.Repository
}

func NewService(repo *Repository, fb *firebase.Client, ar *auth.Repository) *Service {
	return &Service{
		repo:     repo,
		firebase: fb,
		authRepo: ar,
	}
}

type NotificationResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	OrderID   *int64 `json:"order_id,omitempty"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
	Limit         int32                  `json:"limit"`
	Offset        int32                  `json:"offset"`
}

func notificationToResponse(n db.NotificationsLog) NotificationResponse {
	resp := NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Body:      n.Body,
		Status:    n.Status,
		CreatedAt: n.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
	if n.OrderID.Valid {
		resp.OrderID = &n.OrderID.Int64
	}
	return resp
}

func (s *Service) SendOrderNotification(ctx context.Context, userID, orderID int64, title, body string) {
	devices, err := s.authRepo.GetUserDevices(ctx, userID)
	if err != nil || len(devices) == 0 {
		log.Printf("no devices for user %d: %v", userID, err)
		return
	}

	payload := firebase.NotificationPayload{
		Title: title,
		Body:  body,
		Data: map[string]string{
			"type":     "order_update",
			"order_id": strconv.FormatInt(orderID, 10),
		},
	}

	tokens := make([]string, len(devices))
	for i, d := range devices {
		tokens[i] = d.FcmToken
	}

	successCount, sendErrors := s.firebase.SendToMultipleDevices(ctx, tokens, payload)

	status := "sent"
	fcmError := ""
	if len(sendErrors) > 0 {
		status = "failed"
		fcmError = sendErrors[0].Error()
	}
	if successCount > 0 && len(sendErrors) > 0 {
		status = "sent"
	}

	if _, err := s.repo.CreateNotificationLog(ctx, db.CreateNotificationLogParams{
		UserID:     userID,
		OrderID:    pgtype.Int8{Int64: orderID, Valid: true},
		Title:      title,
		Body:       body,
		Status:     status,
		FcmError:   fcmError,
		RetryCount: 0,
	}); err != nil {
		log.Printf("failed to log notification: %v", err)
	}
}

func (s *Service) GetUserNotifications(ctx context.Context, userID int64, limit, offset int32) (*ListResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	notifications, err := s.repo.GetUserNotifications(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountUserNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		resp[i] = notificationToResponse(n)
	}

	return &ListResponse{
		Notifications: resp,
		Total:         total,
		Limit:         limit,
		Offset:        offset,
	}, nil
}
