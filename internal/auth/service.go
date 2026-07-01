package auth

import (
	"context"
	"time"

	"github.com/gumla-hds/gumla-backend/internal/types"
	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	jwtpkg "github.com/gumla-hds/gumla-backend/pkg/jwt"
	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type FirebaseIDTokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (phone string, err error)
}

type Service struct {
	repo        *Repository
	jwtManager  *jwtpkg.Manager
	firebaseVer FirebaseIDTokenVerifier
}

func NewService(repo *Repository, jwtManager *jwtpkg.Manager, firebaseVer FirebaseIDTokenVerifier) *Service {
	return &Service{
		repo:        repo,
		jwtManager:  jwtManager,
		firebaseVer: firebaseVer,
	}
}

type LoginRequest struct {
	Phone   string `json:"phone" binding:"required"`
	IDToken string `json:"id_token" binding:"required"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Phone string `json:"phone"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterDeviceRequest struct {
	FCMToken   string `json:"fcm_token" binding:"required"`
	DeviceInfo string `json:"device_info"`
}

func userToResponse(u db.User) UserResponse {
	return UserResponse{
		ID:    u.ID,
		Phone: u.Phone,
		Name:  u.Name,
		Role:  string(u.Role),
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, *apperrors.AppError) {
	if s.firebaseVer == nil {
		return nil, apperrors.Unauthorized("firebase not configured", "सर्वर कॉन्फ़िगरेशन त्रुटि, कृपया बाद में पुनः प्रयास करें")
	}

	phone, err := s.firebaseVer.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid firebase token", "OTP सत्यापन विफल हुआ, कृपया पुनः प्रयास करें")
	}

	if req.Phone != phone {
		return nil, apperrors.BadRequest("phone mismatch", "फ़ोन नंबर मेल नहीं खाता")
	}

	user, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil {
		user, err = s.repo.CreateUser(ctx, phone, "", types.RoleCustomer)
		if err != nil {
			return nil, apperrors.InternalWrap("failed to create user", "खाता बनाने में विफल, कृपया पुनः प्रयास करें", err)
		}
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Phone, string(user.Role))
	if err != nil {
		return nil, apperrors.InternalWrap("failed to generate tokens", "कुछ गलत हो गया, कृपया पुनः प्रयास करें", err)
	}

	return &LoginResponse{
		User:         userToResponse(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

type DevLoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Name  string `json:"name"`
}

func (s *Service) DevLogin(ctx context.Context, req DevLoginRequest) (*LoginResponse, *apperrors.AppError) {
	user, err := s.repo.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		user, err = s.repo.CreateUser(ctx, req.Phone, req.Name, types.RoleCustomer)
		if err != nil {
			return nil, apperrors.InternalWrap("failed to create user", "खाता बनाने में विफल", err)
		}
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Phone, string(user.Role))
	if err != nil {
		return nil, apperrors.InternalWrap("failed to generate tokens", "कुछ गलत हो गया", err)
	}

	return &LoginResponse{
		User:         userToResponse(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req RefreshRequest) (*RefreshResponse, *apperrors.AppError) {
	claims, err := s.jwtManager.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid refresh token", "आपका सत्र समाप्त हो गया है, कृपया फिर से लॉगिन करें")
	}

	blacklisted, err := s.repo.IsTokenBlacklisted(ctx, claims.ID)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to check token blacklist", "कुछ गलत हो गया, कृपया पुनः प्रयास करें", err)
	}
	if blacklisted {
		return nil, apperrors.Unauthorized("refresh token has been revoked", "आपका सत्र समाप्त हो गया है, कृपया फिर से लॉगिन करें")
	}

	if err := s.repo.RevokeToken(ctx, claims.ID, claims.ExpiresAt.Time); err != nil {
		return nil, apperrors.InternalWrap("failed to revoke old token", "कुछ गलत हो गया, कृपया पुनः प्रयास करें", err)
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(claims.UserID, claims.Phone, claims.Role)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to generate new tokens", "कुछ गलत हो गया, कृपया पुनः प्रयास करें", err)
	}

	return &RefreshResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, accessJTI, refreshJTI string, accessExp, refreshExp time.Time) *apperrors.AppError {
	if accessJTI != "" {
		if err := s.repo.RevokeToken(ctx, accessJTI, accessExp); err != nil {
			return apperrors.InternalWrap("failed to revoke access token", "कुछ गलत हो गया", err)
		}
	}
	if refreshJTI != "" {
		if err := s.repo.RevokeToken(ctx, refreshJTI, refreshExp); err != nil {
			return apperrors.InternalWrap("failed to revoke refresh token", "कुछ गलत हो गया", err)
		}
	}
	return nil
}

func (s *Service) RegisterDevice(ctx context.Context, userID int64, req RegisterDeviceRequest) *apperrors.AppError {
	info := []byte("{}")
	if req.DeviceInfo != "" {
		info = []byte(req.DeviceInfo)
	}
	if err := s.repo.UpsertDevice(ctx, userID, req.FCMToken, info); err != nil {
		return apperrors.InternalWrap("failed to register device", "डिवाइस रजिस्टर करने में विफल", err)
	}
	return nil
}

func (s *Service) RevokeToken(ctx context.Context, jti string, expiresAt time.Time) *apperrors.AppError {
	if err := s.repo.RevokeToken(ctx, jti, expiresAt); err != nil {
		return apperrors.InternalWrap("failed to revoke token", "कुछ गलत हो गया", err)
	}
	return nil
}

func (s *Service) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	return s.repo.IsTokenBlacklisted(ctx, jti)
}
