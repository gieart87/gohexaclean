package errors

import (
	stderrors "errors"
	"net/http"

	"github.com/gieart87/gohexaclean/internal/domain"
	asynqerr "github.com/gieart87/gohexaclean/internal/infra/asynq"
	brokererr "github.com/gieart87/gohexaclean/internal/infra/broker"
	cacheerr "github.com/gieart87/gohexaclean/internal/infra/cache"
	dberr "github.com/gieart87/gohexaclean/internal/infra/db"
)

// MapDomainError maps domain errors to HTTP errors with appropriate status codes
// This function provides a centralized way to convert domain-level errors
// into HTTP-friendly error responses
func MapDomainError(err error) *AppError {
	switch {
	// Domain/Business Logic Errors
	case stderrors.Is(err, domain.ErrUserNotFound):
		return NotFound("User not found", err)
	case stderrors.Is(err, domain.ErrUserAlreadyExists):
		return Conflict("User already exists", err)
	case stderrors.Is(err, domain.ErrInvalidCredentials):
		return Unauthorized("Invalid credentials", err)
	case stderrors.Is(err, domain.ErrUnauthorized):
		return Unauthorized("Unauthorized access", err)
	case stderrors.Is(err, domain.ErrForbidden):
		return Forbidden("Access forbidden", err)
	case stderrors.Is(err, domain.ErrInvalidInput):
		return BadRequest("Invalid input provided", err)

	// Database Infrastructure Errors
	case stderrors.Is(err, dberr.ErrDBConnection):
		return InternalServerError("Database connection failed", err)
	case stderrors.Is(err, dberr.ErrDBTimeout):
		return InternalServerError("Database operation timeout", err)
	case stderrors.Is(err, dberr.ErrDBTransaction):
		return InternalServerError("Database transaction failed", err)
	case stderrors.Is(err, dberr.ErrDBMigration):
		return InternalServerError("Database migration failed", err)
	case stderrors.Is(err, dberr.ErrDBRecordNotFound):
		return NotFound("Record not found", err)
	case stderrors.Is(err, dberr.ErrDBDuplicateKey):
		return Conflict("Duplicate entry", err)
	case stderrors.Is(err, dberr.ErrDBConstraint):
		return BadRequest("Database constraint violation", err)

	// Cache Infrastructure Errors
	case stderrors.Is(err, cacheerr.ErrCacheConnection):
		return InternalServerError("Cache connection failed", err)
	case stderrors.Is(err, cacheerr.ErrCacheTimeout):
		return InternalServerError("Cache operation timeout", err)
	case stderrors.Is(err, cacheerr.ErrCacheKeyNotFound):
		return NotFound("Cache entry not found", err)
	case stderrors.Is(err, cacheerr.ErrCacheMarshal):
		return InternalServerError("Failed to serialize data", err)
	case stderrors.Is(err, cacheerr.ErrCacheUnmarshal):
		return InternalServerError("Failed to deserialize data", err)
	case stderrors.Is(err, cacheerr.ErrCacheExpired):
		return NotFound("Cache entry expired", err)

	// Message Broker Infrastructure Errors
	case stderrors.Is(err, brokererr.ErrBrokerConnection):
		return InternalServerError("Message broker connection failed", err)
	case stderrors.Is(err, brokererr.ErrBrokerPublish):
		return InternalServerError("Failed to publish message", err)
	case stderrors.Is(err, brokererr.ErrBrokerSubscribe):
		return InternalServerError("Failed to subscribe to topic", err)
	case stderrors.Is(err, brokererr.ErrBrokerTimeout):
		return InternalServerError("Message broker timeout", err)
	case stderrors.Is(err, brokererr.ErrBrokerChannelClosed):
		return InternalServerError("Message broker channel closed", err)
	case stderrors.Is(err, brokererr.ErrBrokerAck):
		return InternalServerError("Failed to acknowledge message", err)
	case stderrors.Is(err, brokererr.ErrBrokerNack):
		return InternalServerError("Failed to reject message", err)

	// Asynq Task Queue Infrastructure Errors
	case stderrors.Is(err, asynqerr.ErrTaskEnqueue):
		return InternalServerError("Failed to enqueue task", err)
	case stderrors.Is(err, asynqerr.ErrTaskProcess):
		return InternalServerError("Failed to process task", err)
	case stderrors.Is(err, asynqerr.ErrTaskTimeout):
		return InternalServerError("Task processing timeout", err)
	case stderrors.Is(err, asynqerr.ErrTaskRetry):
		return InternalServerError("Task retry limit exceeded", err)
	case stderrors.Is(err, asynqerr.ErrTaskDuplicate):
		return Conflict("Duplicate task", err)
	case stderrors.Is(err, asynqerr.ErrWorkerStart):
		return InternalServerError("Failed to start worker", err)
	case stderrors.Is(err, asynqerr.ErrWorkerStop):
		return InternalServerError("Failed to stop worker", err)

	default:
		return InternalServerError("Internal server error", err)
	}
}

// MapDomainErrorWithCustomMessage maps domain errors to HTTP errors with custom message
// Use this when you want to provide a more specific error message to the client
func MapDomainErrorWithCustomMessage(err error, customMessage string) *AppError {
	appErr := MapDomainError(err)
	appErr.Message = customMessage
	return appErr
}

// GetHTTPStatusFromDomainError returns the HTTP status code for a domain error
// without creating an AppError instance
func GetHTTPStatusFromDomainError(err error) int {
	switch {
	case stderrors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound
	case stderrors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict
	case stderrors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case stderrors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized
	case stderrors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden
	case stderrors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
