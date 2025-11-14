package user

import (
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/userapi"
	"github.com/gieart87/gohexaclean/internal/port/inbound"
)

// Handler implements userapi.ServerInterface for user-related endpoints
type Handler struct {
	userService inbound.UserServicePort
}

// NewHandler creates a new user handler that implements userapi.ServerInterface
func NewHandler(userService inbound.UserServicePort) *Handler {
	return &Handler{
		userService: userService,
	}
}

// Ensure Handler implements ServerInterface at compile time
var _ userapi.ServerInterface = (*Handler)(nil)
