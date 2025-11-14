package service

import (
	"github.com/google/uuid"
)

// TokenService defines the outbound port for JWT token operations
type TokenService interface {
	GenerateToken(userID uuid.UUID, email string) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
}
