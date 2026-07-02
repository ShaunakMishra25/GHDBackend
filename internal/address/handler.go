package address

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gumla-hds/gumla-backend/internal/middleware"
	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	"github.com/gumla-hds/gumla-backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	addrs, appErr := h.service.List(c.Request.Context(), userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, addrs)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid address id", "गलत पता आईडी"))
		return
	}

	addr, appErr := h.service.GetByID(c.Request.Context(), id, userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, addr)
}

func (h *Handler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	addr, appErr := h.service.Create(c.Request.Context(), userID, req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.Created(c, addr)
}

func (h *Handler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid address id", "गलत पता आईडी"))
		return
	}

	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	addr, appErr := h.service.Update(c.Request.Context(), id, userID, req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, addr)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid address id", "गलत पता आईडी"))
		return
	}

	if appErr := h.service.Delete(c.Request.Context(), id, userID); appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.NoContent(c)
}
