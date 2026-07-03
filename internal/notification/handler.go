package notification

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gumla-hds/gumla-backend/internal/middleware"
	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	"github.com/gumla-hds/gumla-backend/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.Error(c, apperrors.Validation("invalid limit", "गलत सीमा"))
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.Error(c, apperrors.Validation("invalid offset", "गलत ऑफ़सेट"))
		return
	}

	listResp, err := h.svc.GetUserNotifications(c.Request.Context(), userID, int32(limit), int32(offset))
	if err != nil {
		response.Error(c, apperrors.InternalWrap("failed to list notifications", "सूचनाएं लोड करने में विफल", err))
		return
	}

	response.Paginated(c, listResp.Notifications, offset, limit, int(listResp.Total))
}
