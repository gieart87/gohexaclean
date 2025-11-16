package response

import (
	"time"

	"github.com/google/uuid"
)

// Meta represents response metadata
type Meta struct {
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// MetaWithPagination represents metadata with pagination
type MetaWithPagination struct {
	RequestID  string         `json:"request_id"`
	Timestamp  time.Time      `json:"timestamp"`
	Pagination PaginationMeta `json:"pagination"`
}

// SuccessResponse represents a successful response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    Meta        `json:"meta"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool                `json:"success"`
	Message   string              `json:"message"`
	ErrorCode string              `json:"error_code,omitempty"`
	Errors    map[string][]string `json:"errors,omitempty"`
	Meta      Meta                `json:"meta"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    interface{}        `json:"data"`
	Meta    MetaWithPagination `json:"meta"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Meta{
			RequestID: uuid.New().String(),
			Timestamp: time.Now(),
		},
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, err error) *ErrorResponse {
	resp := &ErrorResponse{
		Success: false,
		Message: message,
		Meta: Meta{
			RequestID: uuid.New().String(),
			Timestamp: time.Now(),
		},
	}

	if err != nil {
		// For general errors, set error code as BAD_REQUEST
		resp.ErrorCode = "BAD_REQUEST"
		resp.Errors = map[string][]string{
			"detail": {err.Error()},
		}
	}

	return resp
}

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(message string, errors map[string][]string) *ErrorResponse {
	return &ErrorResponse{
		Success:   false,
		Message:   message,
		ErrorCode: "VALIDATION_ERROR",
		Errors:    errors,
		Meta: Meta{
			RequestID: uuid.New().String(),
			Timestamp: time.Now(),
		},
	}
}

// NewErrorResponseWithCode creates a new error response with custom error code
func NewErrorResponseWithCode(message string, errorCode string, err error) *ErrorResponse {
	resp := &ErrorResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
		Meta: Meta{
			RequestID: uuid.New().String(),
			Timestamp: time.Now(),
		},
	}

	if err != nil {
		resp.Errors = map[string][]string{
			"detail": {err.Error()},
		}
	}

	return resp
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(message string, data interface{}, page, perPage int, total int64) *PaginatedResponse {
	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: MetaWithPagination{
			RequestID: uuid.New().String(),
			Timestamp: time.Now(),
			Pagination: PaginationMeta{
				Page:       page,
				PerPage:    perPage,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	}
}
