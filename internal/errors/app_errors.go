package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with status code and message
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorWithErr creates a new AppError with underlying error
func NewAppErrorWithErr(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors
var (
	ErrNotFound          = NewAppError(http.StatusNotFound, "resource not found")
	ErrUnauthorized      = NewAppError(http.StatusUnauthorized, "unauthorized")
	ErrForbidden         = NewAppError(http.StatusForbidden, "forbidden")
	ErrBadRequest        = NewAppError(http.StatusBadRequest, "bad request")
	ErrConflict          = NewAppError(http.StatusConflict, "resource conflict")
	ErrInternalServer    = NewAppError(http.StatusInternalServerError, "internal server error")
	ErrTooManyRequests   = NewAppError(http.StatusTooManyRequests, "too many requests")
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "invalid email or password")
	ErrEmailExists       = NewAppError(http.StatusConflict, "email already exists")
	ErrWeakPassword      = NewAppError(http.StatusBadRequest, "password must be at least 8 characters long and contain uppercase, lowercase, and number")
)

