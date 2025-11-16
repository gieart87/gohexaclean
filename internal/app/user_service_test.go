package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/gieart87/gohexaclean/internal/dto/request"
	"github.com/gieart87/gohexaclean/internal/infra/config"
	"github.com/gieart87/gohexaclean/internal/port/outbound/repository/mock"
	servicemock "github.com/gieart87/gohexaclean/internal/port/outbound/service/mock"
	"github.com/gieart87/gohexaclean/pkg/crypto"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserServiceTest(t *testing.T) (*UserService, *mock.MockUserRepository, *servicemock.MockCacheService, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockUserRepository(ctrl)
	mockCache := servicemock.NewMockCacheService(ctrl)

	jwtConfig := &config.JWTConfig{
		Secret:  "test-secret",
		Expired: 24,
	}

	service := &UserService{
		userRepo:       mockRepo,
		cacheService:   mockCache,
		jwtConfig:      jwtConfig,
		eventPublisher: nil, // No event publisher in tests (gracefully handled)
	}

	return service, mockRepo, mockCache, ctrl
}

func TestUserService_CreateUser(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	req := &request.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.EXPECT().
		ExistsByEmail(gomock.Any(), req.Email).
		Return(false, nil)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *domain.User) error {
			assert.Equal(t, req.Email, user.Email)
			assert.Equal(t, req.Name, user.Name)
			assert.NotEqual(t, req.Password, user.Password) // Should be hashed
			return nil
		})

	resp, err := service.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, req.Name, resp.User.Name)
}

func TestUserService_CreateUser_EmailAlreadyExists(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	req := &request.CreateUserRequest{
		Email:    "existing@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.EXPECT().
		ExistsByEmail(gomock.Any(), req.Email).
		Return(true, nil)

	resp, err := service.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserAlreadyExists, err)
	assert.Nil(t, resp)
}

func TestUserService_CreateUser_ExistsCheckError(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	req := &request.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.EXPECT().
		ExistsByEmail(gomock.Any(), req.Email).
		Return(false, errors.New("database error"))

	resp, err := service.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUserService_CreateUser_EmptyEmail(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	req := &request.CreateUserRequest{
		Email:    "",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.EXPECT().
		ExistsByEmail(gomock.Any(), req.Email).
		Return(false, nil)

	resp, err := service.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEmailRequired, err)
	assert.Nil(t, resp)
}

func TestUserService_CreateUser_CreateError(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	req := &request.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.EXPECT().
		ExistsByEmail(gomock.Any(), req.Email).
		Return(false, nil)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errors.New("database error"))

	resp, err := service.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUserService_Login(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	password := "password123"
	hashedPassword, err := crypto.HashPassword(password)
	require.NoError(t, err)

	user := &domain.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Name:     "Test User",
		Password: hashedPassword,
	}

	req := &request.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), req.Email).
		Return(user, nil)

	resp, err := service.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, user.Email, resp.User.Email)
	assert.Equal(t, user.Name, resp.User.Name)
}

func TestUserService_Login_InvalidCredentials_UserNotFound(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	req := &request.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), req.Email).
		Return(nil, domain.ErrUserNotFound)

	resp, err := service.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	assert.Nil(t, resp)
}

func TestUserService_Login_InvalidCredentials_WrongPassword(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	hashedPassword, err := crypto.HashPassword("correctpassword")
	require.NoError(t, err)

	user := &domain.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Name:     "Test User",
		Password: hashedPassword,
	}

	req := &request.LoginRequest{
		Email:    user.Email,
		Password: "wrongpassword",
	}

	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), req.Email).
		Return(user, nil)

	resp, err := service.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	assert.Nil(t, resp)
}

func TestUserService_GetUserByID(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.EXPECT().
		FindByID(gomock.Any(), user.ID).
		Return(user, nil)

	resp, err := service.GetUserByID(context.Background(), user.ID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.Name, resp.Name)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	userID := uuid.New()

	mockRepo.EXPECT().
		FindByID(gomock.Any(), userID).
		Return(nil, domain.ErrUserNotFound)

	resp, err := service.GetUserByID(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, resp)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), user.Email).
		Return(user, nil)

	resp, err := service.GetUserByEmail(context.Background(), user.Email)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.Name, resp.Name)
}

func TestUserService_GetUserByEmail_NotFound(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	email := "notfound@example.com"

	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), email).
		Return(nil, domain.ErrUserNotFound)

	resp, err := service.GetUserByEmail(context.Background(), email)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, resp)
}

func TestUserService_UpdateUser(t *testing.T) {
	service, mockRepo, mockCache, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	user := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Old Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &request.UpdateUserRequest{
		Name: "New Name",
	}

	mockRepo.EXPECT().
		FindByID(gomock.Any(), userID).
		Return(user, nil)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *domain.User) error {
			assert.Equal(t, req.Name, u.Name)
			return nil
		})

	mockCache.EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(nil)

	resp, err := service.UpdateUser(context.Background(), userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Name, resp.Name)
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	req := &request.UpdateUserRequest{
		Name: "New Name",
	}

	mockRepo.EXPECT().
		FindByID(gomock.Any(), userID).
		Return(nil, domain.ErrUserNotFound)

	resp, err := service.UpdateUser(context.Background(), userID, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, resp)
}

func TestUserService_UpdateUser_UpdateError(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	user := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Old Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &request.UpdateUserRequest{
		Name: "New Name",
	}

	mockRepo.EXPECT().
		FindByID(gomock.Any(), userID).
		Return(user, nil)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(errors.New("database error"))

	resp, err := service.UpdateUser(context.Background(), userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUserService_DeleteUser(t *testing.T) {
	service, mockRepo, mockCache, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	userID := uuid.New()

	mockRepo.EXPECT().
		Delete(gomock.Any(), userID).
		Return(nil)

	mockCache.EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(nil)

	err := service.DeleteUser(context.Background(), userID)

	assert.NoError(t, err)
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	userID := uuid.New()

	mockRepo.EXPECT().
		Delete(gomock.Any(), userID).
		Return(domain.ErrUserNotFound)

	err := service.DeleteUser(context.Background(), userID)

	assert.Error(t, err)
}

func TestUserService_ListUsers(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	users := []*domain.User{
		{
			ID:        uuid.New(),
			Email:     "user1@example.com",
			Name:      "User 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Email:     "user2@example.com",
			Name:      "User 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	page := 1
	limit := 10
	offset := 0
	total := int64(2)

	mockRepo.EXPECT().
		List(gomock.Any(), offset, limit).
		Return(users, nil)

	mockRepo.EXPECT().
		Count(gomock.Any()).
		Return(total, nil)

	resp, totalCount, err := service.ListUsers(context.Background(), page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp, 2)
	assert.Equal(t, total, totalCount)
	assert.Equal(t, users[0].Email, resp[0].Email)
	assert.Equal(t, users[1].Email, resp[1].Email)
}

func TestUserService_ListUsers_ListError(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	page := 1
	limit := 10
	offset := 0

	mockRepo.EXPECT().
		List(gomock.Any(), offset, limit).
		Return(nil, errors.New("database error"))

	resp, totalCount, err := service.ListUsers(context.Background(), page, limit)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, int64(0), totalCount)
}

func TestUserService_ListUsers_CountError(t *testing.T) {
	service, mockRepo, _, ctrl := setupUserServiceTest(t)
	defer ctrl.Finish()

	users := []*domain.User{
		{
			ID:        uuid.New(),
			Email:     "user1@example.com",
			Name:      "User 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	page := 1
	limit := 10
	offset := 0

	mockRepo.EXPECT().
		List(gomock.Any(), offset, limit).
		Return(users, nil)

	mockRepo.EXPECT().
		Count(gomock.Any()).
		Return(int64(0), errors.New("database error"))

	resp, totalCount, err := service.ListUsers(context.Background(), page, limit)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, int64(0), totalCount)
}
