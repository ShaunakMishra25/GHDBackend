package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gumla-hds/gumla-backend/internal/types"
	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		queries: db.New(pool),
	}
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (db.User, error) {
	return r.queries.GetUserByPhone(ctx, phone)
}

func (r *Repository) GetUserByID(ctx context.Context, id int64) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *Repository) CreateUser(ctx context.Context, phone, name string, role types.UserRole) (db.User, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		Phone: phone,
		Name:  name,
		Role:  role,
	})
}

func (r *Repository) UpsertDevice(ctx context.Context, userID int64, fcmToken string, deviceInfo []byte) error {
	return r.queries.UpsertDevice(ctx, db.UpsertDeviceParams{
		UserID:     userID,
		FcmToken:   fcmToken,
		DeviceInfo: deviceInfo,
	})
}

func (r *Repository) GetUserDevices(ctx context.Context, userID int64) ([]db.Device, error) {
	return r.queries.GetUserDevices(ctx, userID)
}

func (r *Repository) RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error {
	var ts pgtype.Timestamptz
	if err := ts.Scan(expiresAt); err != nil {
		return err
	}
	return r.queries.RevokeToken(ctx, db.RevokeTokenParams{
		TokenJti:  jti,
		ExpiresAt: ts,
	})
}

func (r *Repository) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	return r.queries.IsTokenBlacklisted(ctx, jti)
}

func (r *Repository) UpdateUser(ctx context.Context, id int64, name string) error {
	return r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:   id,
		Name: name,
	})
}
