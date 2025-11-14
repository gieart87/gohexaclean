package app

import (
	"context"
	"fmt"

	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/gieart87/gohexaclean/internal/dto/request"
	"github.com/gieart87/gohexaclean/internal/dto/response"
	"github.com/gieart87/gohexaclean/internal/port/inbound"
	"github.com/gieart87/gohexaclean/internal/port/outbound/repository"
	"github.com/gieart87/gohexaclean/internal/port/outbound/service"
	"github.com/google/uuid"
)

// UserService implements the UserServicePort interface
type UserService struct {
	userRepo     repository.UserRepository
	hashService  service.HashService
	tokenService service.TokenService
	cacheService service.CacheService
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	hashService service.HashService,
	tokenService service.TokenService,
	cacheService service.CacheService,
) inbound.UserServicePort {
	return &UserService{
		userRepo:     userRepo,
		hashService:  hashService,
		tokenService: tokenService,
		cacheService: cacheService,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.UserResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.hashService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create domain entity
	user := domain.NewUser(req.Email, req.Name, hashedPassword)

	// Validate domain entity
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return response.NewUserResponse(user), nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.NewUserResponse(user), nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return response.NewUserResponse(user), nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req *request.UpdateUserRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.UpdateProfile(req.Name)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", id.String())
	_ = s.cacheService.Delete(ctx, cacheKey)

	return response.NewUserResponse(user), nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", id.String())
	_ = s.cacheService.Delete(ctx, cacheKey)

	return nil
}

// Login authenticates a user and returns a token
func (s *UserService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check password
	if !s.hashService.CheckPasswordHash(req.Password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate token
	token, err := s.tokenService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &response.LoginResponse{
		Token: token,
		User:  response.NewUserResponse(user),
	}, nil
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, page, limit int) ([]*response.UserResponse, int64, error) {
	offset := (page - 1) * limit

	users, err := s.userRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	userResponses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = response.NewUserResponse(user)
	}

	return userResponses, total, nil
}
