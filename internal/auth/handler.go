package auth

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"

	"github.com/gumla-hds/gumla-backend/internal/middleware"
	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	jwtpkg "github.com/gumla-hds/gumla-backend/pkg/jwt"
	"github.com/gumla-hds/gumla-backend/pkg/response"
)

type FirebaseVerifier struct {
	client *auth.Client
}

func NewFirebaseVerifier(client *auth.Client) *FirebaseVerifier {
	return &FirebaseVerifier{client: client}
}

func (fv *FirebaseVerifier) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token, err := fv.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	phone, ok := token.Claims["phone_number"].(string)
	if !ok || phone == "" {
		return "", fmt.Errorf("phone number not found in token")
	}
	return phone, nil
}

type Handler struct {
	service    *Service
	jwtManager *jwtpkg.Manager
}

func NewHandler(service *Service, jwtManager *jwtpkg.Manager) *Handler {
	return &Handler{
		service:    service,
		jwtManager: jwtManager,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	resp, appErr := h.service.Login(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) DevLogin(c *gin.Context) {
	var req DevLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	resp, appErr := h.service.DevLogin(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	resp, appErr := h.service.RefreshToken(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) Logout(c *gin.Context) {
	accessJTI := middleware.GetTokenJTI(c)

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		claims, err := h.jwtManager.VerifyToken(req.RefreshToken)
		if err == nil {
			_ = h.service.RevokeToken(c.Request.Context(), claims.ID, claims.ExpiresAt.Time)
		}
	}

	appErr := h.service.Logout(c.Request.Context(), accessJTI, "", time.Now(), time.Now().Add(24*time.Hour))
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.NoContent(c)
}

func (h *Handler) RegisterDevice(c *gin.Context) {
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	userID := middleware.GetUserID(c)
	appErr := h.service.RegisterDevice(c.Request.Context(), userID, req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.NoContent(c)
}
