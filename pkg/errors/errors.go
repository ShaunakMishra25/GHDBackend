package errors

import "net/http"

type Code string

const (
	CodeNotFound       Code = "NOT_FOUND"
	CodeBadRequest     Code = "BAD_REQUEST"
	CodeUnauthorized   Code = "UNAUTHORIZED"
	CodeForbidden      Code = "FORBIDDEN"
	CodeInternal       Code = "INTERNAL_ERROR"
	CodeConflict       Code = "CONFLICT"
	CodeValidation     Code = "VALIDATION_ERROR"
	CodePaymentFailed  Code = "PAYMENT_FAILED"
)

var codeToStatus = map[Code]int{
	CodeNotFound:      http.StatusNotFound,
	CodeBadRequest:    http.StatusBadRequest,
	CodeUnauthorized:  http.StatusUnauthorized,
	CodeForbidden:     http.StatusForbidden,
	CodeInternal:      http.StatusInternalServerError,
	CodeConflict:      http.StatusConflict,
	CodeValidation:    http.StatusUnprocessableEntity,
	CodePaymentFailed: http.StatusPaymentRequired,
}

type AppError struct {
	Code     Code   `json:"code"`
	Message  string `json:"message"`
	UserMsg  string `json:"user_msg"`
	HTTPCode int    `json:"-"`
	Err      error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code Code, message string, userMsg string) *AppError {
	status, ok := codeToStatus[code]
	if !ok {
		status = http.StatusInternalServerError
	}
	return &AppError{
		Code:     code,
		Message:  message,
		UserMsg:  userMsg,
		HTTPCode: status,
	}
}

func Wrap(code Code, message string, userMsg string, err error) *AppError {
	appErr := New(code, message, userMsg)
	appErr.Err = err
	return appErr
}

func NotFound(message string, userMsg string) *AppError {
	return New(CodeNotFound, message, userMsg)
}

func NotFoundWrap(message string, userMsg string, err error) *AppError {
	return Wrap(CodeNotFound, message, userMsg, err)
}

func BadRequest(message string, userMsg string) *AppError {
	return New(CodeBadRequest, message, userMsg)
}

func BadRequestWrap(message string, userMsg string, err error) *AppError {
	return Wrap(CodeBadRequest, message, userMsg, err)
}

func Unauthorized(message string, userMsg string) *AppError {
	return New(CodeUnauthorized, message, userMsg)
}

func Forbidden(message string, userMsg string) *AppError {
	return New(CodeForbidden, message, userMsg)
}

func Internal(message string, userMsg string) *AppError {
	return New(CodeInternal, message, userMsg)
}

func InternalWrap(message string, userMsg string, err error) *AppError {
	return Wrap(CodeInternal, message, userMsg, err)
}

func Conflict(message string, userMsg string) *AppError {
	return New(CodeConflict, message, userMsg)
}

func Validation(message string, userMsg string) *AppError {
	return New(CodeValidation, message, userMsg)
}

func PaymentFailed(message string, userMsg string) *AppError {
	return New(CodePaymentFailed, message, userMsg)
}
