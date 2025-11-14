package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
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

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error constructors
func BadRequest(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

func Unauthorized(message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, message, err)
}

func Forbidden(message string, err error) *AppError {
	return NewAppError(http.StatusForbidden, message, err)
}

func NotFound(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

func Conflict(message string, err error) *AppError {
	return NewAppError(http.StatusConflict, message, err)
}

func InternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

// GetHTTPStatusCode gets HTTP status code from error
func GetHTTPStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return http.StatusInternalServerError
}
