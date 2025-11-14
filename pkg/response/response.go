package response

import "time"

// SuccessResponse represents a successful response
type SuccessResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Error     string      `json:"error,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Meta      MetaData    `json:"meta"`
	Timestamp time.Time   `json:"timestamp"`
}

// MetaData represents pagination metadata
type MetaData struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, err error) *ErrorResponse {
	resp := &ErrorResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	}

	if err != nil {
		resp.Error = err.Error()
	}

	return resp
}

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(message string, errors interface{}) *ErrorResponse {
	return &ErrorResponse{
		Success:   false,
		Message:   message,
		Errors:    errors,
		Timestamp: time.Now(),
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(message string, data interface{}, page, limit int, total int64) *PaginatedResponse {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: MetaData{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Timestamp: time.Now(),
	}
}
