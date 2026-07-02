package cart

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

func (h *Handler) GetCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	cart, appErr := h.service.GetCart(c.Request.Context(), userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, cart)
}

func (h *Handler) AddItem(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	if appErr := h.service.AddItem(c.Request.Context(), userID, req); appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.NoContent(c)
}

func (h *Handler) UpdateItem(c *gin.Context) {
	userID := middleware.GetUserID(c)

	itemID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid item id", "गलत आइटम आईडी"))
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	if appErr := h.service.UpdateItem(c.Request.Context(), userID, itemID, req); appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.NoContent(c)
}

func (h *Handler) RemoveItem(c *gin.Context) {
	userID := middleware.GetUserID(c)

	itemID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid item id", "गलत आइटम आईडी"))
		return
	}

	if appErr := h.service.RemoveItem(c.Request.Context(), userID, itemID); appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.NoContent(c)
}

func (h *Handler) ClearCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if appErr := h.service.ClearCart(c.Request.Context(), userID); appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.NoContent(c)
}
