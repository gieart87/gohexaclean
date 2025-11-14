package inbound

import (
	"context"

	"github.com/gieart87/gohexaclean/internal/dto/request"
	"github.com/gieart87/gohexaclean/internal/dto/response"
	"github.com/google/uuid"
)

// UserServicePort defines the inbound port for user service (use case interface)
// This is what the adapters (HTTP, gRPC) will call
type UserServicePort interface {
	CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.LoginResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*response.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *request.UpdateUserRequest) (*response.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error)
	ListUsers(ctx context.Context, page, limit int) ([]*response.UserResponse, int64, error)
}
