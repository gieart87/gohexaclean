package health

import (
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/healthapi"
)

// Handler implements healthapi.ServerInterface for health check endpoint
type Handler struct{}

// NewHandler creates a new health handler that implements healthapi.ServerInterface
func NewHandler() *Handler {
	return &Handler{}
}

// Ensure Handler implements ServerInterface at compile time
var _ healthapi.ServerInterface = (*Handler)(nil)
