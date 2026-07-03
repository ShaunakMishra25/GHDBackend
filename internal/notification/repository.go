package notification

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Repository struct {
	Queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Queries: db.New(pool),
	}
}

func (r *Repository) CreateNotificationLog(ctx context.Context, arg db.CreateNotificationLogParams) (db.NotificationsLog, error) {
	return r.Queries.CreateNotificationLog(ctx, arg)
}

func (r *Repository) UpdateNotificationStatus(ctx context.Context, id int64, status, fcmError string, retryCount int32) error {
	return r.Queries.UpdateNotificationStatus(ctx, db.UpdateNotificationStatusParams{
		ID:         id,
		Status:     status,
		FcmError:   fcmError,
		RetryCount: retryCount,
	})
}

func (r *Repository) GetUserNotifications(ctx context.Context, userID int64, limit, offset int32) ([]db.NotificationsLog, error) {
	return r.Queries.GetUserNotifications(ctx, db.GetUserNotificationsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *Repository) CountUserNotifications(ctx context.Context, userID int64) (int64, error) {
	return r.Queries.CountUserNotifications(ctx, userID)
}
