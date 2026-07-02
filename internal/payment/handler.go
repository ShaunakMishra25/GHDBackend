package payment

import (
	"net/http"

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

// Initiate POST /payments/initiate
// Verify   POST /payments/verify

func (h *Handler) Initiate(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.Validation("invalid request body", "गलत डेटा"))
		return
	}

	resp, appErr := h.svc.Initiate(c.Request.Context(), userID, req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.Created(c, resp)
}

func (h *Handler) Verify(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.Validation("invalid request body", "गलत डेटा"))
		return
	}

	appErr := h.svc.Verify(c.Request.Context(), userID, req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.OK(c, gin.H{"message": "payment verified successfully"})
}

func (h *Handler) Webhook(c *gin.Context) {
	signature := c.GetHeader("X-Razorpay-Signature")
	if signature == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValid := h.svc.VerifyWebhookSignature(body, signature)
	if !isValid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	event, _ := payload["event"].(string)
	if event != "payment.captured" && event != "order.paid" {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	rpOrderID := ""
	if pay, ok := payload["payload"].(map[string]interface{}); ok {
		if ord, ok := pay["order"].(map[string]interface{}); ok {
			rpOrderID, _ = ord["id"].(string)
		}
	}

	if rpOrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	appErr := h.svc.HandleWebhook(c.Request.Context(), rpOrderID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
