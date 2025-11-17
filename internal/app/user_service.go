package app

import (
	"context"
	"fmt"
	"log"

	"github.com/gieart87/gohexaclean/internal/adapter/outbound/event"
	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/gieart87/gohexaclean/internal/dto/request"
	"github.com/gieart87/gohexaclean/internal/dto/response"
	"github.com/gieart87/gohexaclean/internal/infra/config"
	"github.com/gieart87/gohexaclean/internal/infrastructure/asynq/tasks"
	"github.com/gieart87/gohexaclean/internal/port/inbound"
	"github.com/gieart87/gohexaclean/internal/port/outbound/repository"
	"github.com/gieart87/gohexaclean/internal/port/outbound/service"
	"github.com/gieart87/gohexaclean/pkg/auth"
	"github.com/gieart87/gohexaclean/pkg/crypto"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// UserService implements the UserServicePort interface
type UserService struct {
	userRepo       repository.UserRepository
	cacheService   service.CacheService
	jwtConfig      *config.JWTConfig
	eventPublisher *event.UserEventPublisher
	taskClient     *asynq.Client
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	cacheService service.CacheService,
	jwtConfig *config.JWTConfig,
	eventPublisher *event.UserEventPublisher,
	taskClient *asynq.Client,
) inbound.UserServicePort {
	return &UserService{
		userRepo:       userRepo,
		cacheService:   cacheService,
		jwtConfig:      jwtConfig,
		eventPublisher: eventPublisher,
		taskClient:     taskClient,
	}
}

// CreateUser creates a new user and returns a token
func (s *UserService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.LoginResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(req.Password)
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

	// Generate token for the newly registered user
	token, err := auth.GenerateJWT(user.ID, user.Email, s.jwtConfig.Secret, s.jwtConfig.Expired)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Publish user created event
	if s.eventPublisher != nil {
		event := domain.NewUserCreatedEvent(user.ID, user.Email, user.Name)
		if err := s.eventPublisher.PublishUserCreated(ctx, event); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to publish user created event: %v\n", err)
		}
	}

	// Enqueue welcome email task asynchronously
	if s.taskClient != nil {
		task, err := tasks.NewEmailWelcomeTask(user.ID.String(), user.Email, user.Name)
		if err != nil {
			log.Printf("failed to create welcome email task: %v", err)
		} else {
			info, err := s.taskClient.Enqueue(task)
			if err != nil {
				log.Printf("failed to enqueue welcome email task: %v", err)
			} else {
				log.Printf("enqueued welcome email task: id=%s queue=%s", info.ID, info.Queue)
			}
		}
	}

	return &response.LoginResponse{
		Token: token,
		User:  response.NewUserResponse(user),
	}, nil
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

	// Publish user updated event
	if s.eventPublisher != nil {
		event := domain.NewUserUpdatedEvent(user.ID, user.Name)
		if err := s.eventPublisher.PublishUserUpdated(ctx, event); err != nil {
			fmt.Printf("failed to publish user updated event: %v\n", err)
		}
	}

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

	// Publish user deleted event
	if s.eventPublisher != nil {
		event := domain.NewUserDeletedEvent(id)
		if err := s.eventPublisher.PublishUserDeleted(ctx, event); err != nil {
			fmt.Printf("failed to publish user deleted event: %v\n", err)
		}
	}

	return nil
}

// Login authenticates a user and returns a token
func (s *UserService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check password
	if !crypto.CheckPasswordHash(req.Password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate token
	token, err := auth.GenerateJWT(user.ID, user.Email, s.jwtConfig.Secret, s.jwtConfig.Expired)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Publish user logged in event
	if s.eventPublisher != nil {
		event := domain.NewUserLoggedInEvent(user.ID, user.Email)
		if err := s.eventPublisher.PublishUserLoggedIn(ctx, event); err != nil {
			fmt.Printf("failed to publish user logged in event: %v\n", err)
		}
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
