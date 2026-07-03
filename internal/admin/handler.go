package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gumla-hds/gumla-backend/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Dashboard(c *gin.Context) {
	dash, appErr := h.svc.GetDashboard(c.Request.Context())
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.OK(c, dash)
}
