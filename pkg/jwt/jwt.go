package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessJTI    string `json:"-"`
	RefreshJTI   string `json:"-"`
}

type Manager struct {
	secret         []byte
	accessExpiry   time.Duration
	refreshExpiry  time.Duration
}

func NewManager(secret string, accessExpiry, refreshExpiry time.Duration) *Manager {
	return &Manager{
		secret:        []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (m *Manager) GenerateTokenPair(userID int64, phone, role string) (*TokenPair, error) {
	accessJTI := uuid.New().String()
	refreshJTI := uuid.New().String()
	now := time.Now()

	accessClaims := Claims{
		UserID: userID,
		Phone:  phone,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpiry)),
			Issuer:    "gumla-hds",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := Claims{
		UserID: userID,
		Phone:  phone,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshExpiry)),
			Issuer:    "gumla-hds",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		AccessJTI:    accessJTI,
		RefreshJTI:   refreshJTI,
	}, nil
}

func (m *Manager) VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
