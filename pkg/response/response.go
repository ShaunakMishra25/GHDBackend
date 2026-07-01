package response

import (
	"net/http"

	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool             `json:"success"`
	Error   ErrorDetail      `json:"error"`
}

type ErrorDetail struct {
	Code    apperrors.Code `json:"code"`
	Message string         `json:"message"`
	UserMsg string         `json:"user_msg"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{Success: true, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{Success: true, Data: data})
}

func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func Paginated(c *gin.Context, data interface{}, offset, limit, total int) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			Offset: offset,
			Limit:  limit,
			Total:  total,
		},
	})
}

func Error(c *gin.Context, appErr *apperrors.AppError) {
	c.JSON(appErr.HTTPCode, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
			UserMsg: appErr.UserMsg,
		},
	})
}

func ErrorWithStatus(c *gin.Context, status int, code apperrors.Code, message, userMsg string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			UserMsg: userMsg,
		},
	})
}
