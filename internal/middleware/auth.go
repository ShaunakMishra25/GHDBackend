package middleware

import (
	"context"
	"strings"

	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	jwtpkg "github.com/gumla-hds/gumla-backend/pkg/jwt"
	"github.com/gumla-hds/gumla-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	ContextUserID  contextKey = "user_id"
	ContextPhone   contextKey = "phone"
	ContextRole    contextKey = "role"
	ContextTokenJTI contextKey = "token_jti"
)

func RequireAuth(jwtManager *jwtpkg.Manager, isTokenBlacklisted func(ctx context.Context, jti string) (bool, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperrors.Unauthorized("missing authorization header", "कृपया लॉगिन करें"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			response.Error(c, apperrors.Unauthorized("invalid authorization format", "कृपया लॉगिन करें"))
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := jwtManager.VerifyToken(tokenStr)
		if err != nil {
			response.Error(c, apperrors.Unauthorized("invalid or expired token", "आपका सत्र समाप्त हो गया है, कृपया फिर से लॉगिन करें"))
			c.Abort()
			return
		}

		blacklisted, err := isTokenBlacklisted(c.Request.Context(), claims.ID)
		if err != nil {
			response.Error(c, apperrors.Internal("token verification failed", "कुछ गलत हो गया, कृपया पुनः प्रयास करें"))
			c.Abort()
			return
		}
		if blacklisted {
			response.Error(c, apperrors.Unauthorized("token has been revoked", "आपका सत्र समाप्त हो गया है, कृपया फिर से लॉगिन करें"))
			c.Abort()
			return
		}

		c.Set(string(ContextUserID), claims.UserID)
		c.Set(string(ContextPhone), claims.Phone)
		c.Set(string(ContextRole), claims.Role)
		c.Set(string(ContextTokenJTI), claims.ID)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(string(ContextRole))
		if !exists {
			response.Error(c, apperrors.Unauthorized("not authenticated", "कृपया लॉगिन करें"))
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			response.Error(c, apperrors.Internal("invalid role type", "कुछ गलत हो गया"))
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		response.Error(c, apperrors.Forbidden("insufficient permissions", "आपके पास इस कार्य के लिए अनुमति नहीं है"))
		c.Abort()
	}
}

func GetUserID(c *gin.Context) int64 {
	id, _ := c.Get(string(ContextUserID))
	userID, _ := id.(int64)
	return userID
}

func GetRole(c *gin.Context) string {
	role, _ := c.Get(string(ContextRole))
	roleStr, _ := role.(string)
	return roleStr
}

func GetTokenJTI(c *gin.Context) string {
	jti, _ := c.Get(string(ContextTokenJTI))
	jtiStr, _ := jti.(string)
	return jtiStr
}
