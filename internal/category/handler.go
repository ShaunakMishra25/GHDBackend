package category

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
	categories, appErr := h.service.List(c.Request.Context())
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, categories)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid category id", "गलत श्रेणी आईडी"))
		return
	}

	cat, appErr := h.service.GetByID(c.Request.Context(), id)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, cat)
}

func (h *Handler) Create(c *gin.Context) {
	userRole := middleware.GetRole(c)
	if userRole != "admin" {
		response.Error(c, apperrors.Forbidden("admin only", "केवल एडमिन के लिए"))
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	cat, appErr := h.service.Create(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.Created(c, cat)
}

func (h *Handler) Update(c *gin.Context) {
	userRole := middleware.GetRole(c)
	if userRole != "admin" {
		response.Error(c, apperrors.Forbidden("admin only", "केवल एडमिन के लिए"))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid category id", "गलत श्रेणी आईडी"))
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	cat, appErr := h.service.Update(c.Request.Context(), id, req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, cat)
}

func (h *Handler) Delete(c *gin.Context) {
	userRole := middleware.GetRole(c)
	if userRole != "admin" {
		response.Error(c, apperrors.Forbidden("admin only", "केवल एडमिन के लिए"))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid category id", "गलत श्रेणी आईडी"))
		return
	}

	if appErr := h.service.Delete(c.Request.Context(), id); appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.NoContent(c)
}
