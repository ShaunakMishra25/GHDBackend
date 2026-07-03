package order

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gumla-hds/gumla-backend/internal/middleware"
	"github.com/gumla-hds/gumla-backend/internal/types"
	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	"github.com/gumla-hds/gumla-backend/pkg/response"
)

type NotificationSender interface {
	SendOrderNotification(ctx context.Context, userID, orderID int64, title, body string)
}

type Handler struct {
	svc      *Service
	notifSvc NotificationSender
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) SetNotificationSender(ns NotificationSender) {
	h.notifSvc = ns
}

func (h *Handler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.Validation("invalid request body", "गलत डेटा"))
		return
	}

	order, appErr := h.svc.Create(c.Request.Context(), userID, req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.Created(c, order)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		response.Error(c, apperrors.Validation("invalid order id", "गलत ऑर्डर आईडी"))
		return
	}

	if role == "admin" {
		order, appErr := h.svc.GetByIDAdmin(c.Request.Context(), orderID)
		if appErr != nil {
			response.Error(c, appErr)
			return
		}
		response.OK(c, order)
		return
	}

	order, appErr := h.svc.GetByID(c.Request.Context(), userID, orderID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.OK(c, order)
}

func (h *Handler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

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

	if role == "admin" {
		status := c.Query("status")
		listResp, appErr := h.svc.ListAll(c.Request.Context(), status, int32(limit), int32(offset))
		if appErr != nil {
			response.Error(c, appErr)
			return
		}
		response.Paginated(c, listResp.Orders, offset, limit, int(listResp.Total))
		return
	}

	listResp, appErr := h.svc.ListUserOrders(c.Request.Context(), userID, int32(limit), int32(offset))
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.Paginated(c, listResp.Orders, offset, limit, int(listResp.Total))
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	role := middleware.GetRole(c)
	userID := middleware.GetUserID(c)

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		response.Error(c, apperrors.Validation("invalid order id", "गलत ऑर्डर आईडी"))
		return
	}

	var req StatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.Validation("invalid request body", "गलत डेटा"))
		return
	}

	var order *OrderResponse
	var appErr *apperrors.AppError

	if role == "admin" {
		order, appErr = h.svc.UpdateStatus(c.Request.Context(), orderID, req, "admin")
	} else {
		if req.Status != string(types.OrderCancelled) {
			response.Error(c, apperrors.Forbidden("only admin can update order status", "केवल एडमिन ऑर्डर की स्थिति बदल सकता है"))
			return
		}
		order, appErr = h.svc.GetByID(c.Request.Context(), userID, orderID)
		if appErr != nil {
			response.Error(c, appErr)
			return
		}
		order, appErr = h.svc.UpdateStatus(c.Request.Context(), orderID, req, "user")
	}

	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	if h.notifSvc != nil && order != nil {
		go h.fireStatusNotification(c.Request.Context(), order, role)
	}

	response.OK(c, order)
}

func (h *Handler) fireStatusNotification(ctx context.Context, order *OrderResponse, actorRole string) {
	title := "ऑर्डर अपडेट"
	var body string
	switch types.OrderStatus(order.Status) {
	case types.OrderConfirmed:
		body = "आपका ऑर्डर कन्फर्म हो गया है"
	case types.OrderAccepted:
		body = "आपका ऑर्डर स्वीकार कर लिया गया है"
	case types.OrderPreparing:
		body = "आपका ऑर्डर तैयार किया जा रहा है"
	case types.OrderDispatched:
		body = "आपका ऑर्डर डिस्पैच हो गया है"
	case types.OrderDelivered:
		body = "आपका ऑर्डर डिलीवर कर दिया गया है"
	case types.OrderCancelled:
		title = "ऑर्डर रद्द"
		body = "आपका ऑर्डर रद्द कर दिया गया है"
	default:
		return
	}
	h.notifSvc.SendOrderNotification(ctx, order.UserID, order.ID, title, body)
}
