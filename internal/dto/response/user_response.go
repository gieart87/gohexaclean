package response

import (
	"time"

	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/google/uuid"
)

// UserResponse represents the user response DTO
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserResponse creates a new user response from domain model
func NewUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}
