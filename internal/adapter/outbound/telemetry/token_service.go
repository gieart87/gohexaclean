package telemetry

import (
	"fmt"
	"time"

	"github.com/gieart87/gohexaclean/internal/infra/config"
	"github.com/gieart87/gohexaclean/internal/port/outbound/service"
	"github.com/gieart87/gohexaclean/pkg/utils"
	"github.com/google/uuid"
)

// TokenServiceImpl implements TokenService interface
type TokenServiceImpl struct {
	jwtSecret string
	jwtExpiry string
}

// NewTokenService creates a new token service
func NewTokenService(cfg *config.JWTConfig) service.TokenService {
	return &TokenServiceImpl{
		jwtSecret: cfg.Secret,
		jwtExpiry: cfg.Expired.String(),
	}
}

// GenerateToken generates a JWT token
func (s *TokenServiceImpl) GenerateToken(userID uuid.UUID, email string) (string, error) {
	expiration, err := time.ParseDuration(s.jwtExpiry)
	if err != nil {
		expiration = 24 * time.Hour // default to 24 hours
	}

	return utils.GenerateJWT(userID, email, s.jwtSecret, expiration)
}

// ValidateToken validates a JWT token
func (s *TokenServiceImpl) ValidateToken(token string) (uuid.UUID, error) {
	claims, err := utils.ValidateJWT(token, s.jwtSecret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims.UserID, nil
}
