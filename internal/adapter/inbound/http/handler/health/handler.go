package health

import (
	"time"

	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/healthapi"
)

// Handler implements healthapi.ServerInterface for health check endpoint
type Handler struct {
	serviceName string
	version     string
	startedAt   time.Time
}

// NewHandler creates a new health handler that implements healthapi.ServerInterface
func NewHandler(serviceName, version string) *Handler {
	return &Handler{
		serviceName: serviceName,
		version:     version,
		startedAt:   time.Now(),
	}
}

// Ensure Handler implements ServerInterface at compile time
var _ healthapi.ServerInterface = (*Handler)(nil)
