package telemetry

import (
	"github.com/gieart87/gohexaclean/internal/port/outbound/service"
	"github.com/gieart87/gohexaclean/pkg/utils"
)

// HashServiceImpl implements HashService interface
type HashServiceImpl struct{}

// NewHashService creates a new hash service
func NewHashService() service.HashService {
	return &HashServiceImpl{}
}

// HashPassword hashes a password
func (s *HashServiceImpl) HashPassword(password string) (string, error) {
	return utils.HashPassword(password)
}

// CheckPasswordHash checks if password matches hash
func (s *HashServiceImpl) CheckPasswordHash(password, hash string) bool {
	return utils.CheckPasswordHash(password, hash)
}
