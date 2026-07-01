package product

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
	categoryIDStr := c.Query("category_id")
	var categoryID *int64
	if categoryIDStr != "" {
		id, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			response.Error(c, apperrors.BadRequest("invalid category_id", "गलत श्रेणी आईडी"))
			return
		}
		categoryID = &id
	}

	searchStr := c.Query("search")
	var search *string
	if searchStr != "" {
		search = &searchStr
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		limit = 20
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		offset = 0
	}

	resp, appErr := h.service.List(c.Request.Context(), categoryID, search, int32(limit), int32(offset))
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid product id", "गलत उत्पाद आईडी"))
		return
	}

	p, appErr := h.service.GetByID(c.Request.Context(), id)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Create(c *gin.Context) {
	userRole := middleware.GetRole(c)
	if userRole != "admin" {
		response.Error(c, apperrors.Forbidden("admin only", "केवल एडमिन के लिए"))
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	p, appErr := h.service.Create(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.Created(c, p)
}

func (h *Handler) Update(c *gin.Context) {
	userRole := middleware.GetRole(c)
	if userRole != "admin" {
		response.Error(c, apperrors.Forbidden("admin only", "केवल एडमिन के लिए"))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid product id", "गलत उत्पाद आईडी"))
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", "कृपया सही जानकारी दर्ज करें"))
		return
	}

	p, appErr := h.service.Update(c.Request.Context(), id, req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Delete(c *gin.Context) {
	userRole := middleware.GetRole(c)
	if userRole != "admin" {
		response.Error(c, apperrors.Forbidden("admin only", "केवल एडमिन के लिए"))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid product id", "गलत उत्पाद आईडी"))
		return
	}

	if appErr := h.service.Delete(c.Request.Context(), id); appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.NoContent(c)
}
