package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/userapi"
	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/gieart87/gohexaclean/internal/dto/response"
	"github.com/gieart87/gohexaclean/internal/port/inbound/mock"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupHandlerTest(t *testing.T) (*Handler, *mock.MockUserServicePort, *gomock.Controller, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mock.NewMockUserServicePort(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()

	return handler, mockService, ctrl, app
}

func TestHandler_CreateUser(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Post("/users", handler.CreateUser)

	req := userapi.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	userResp := &response.UserResponse{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	loginResp := &response.LoginResponse{
		Token: "jwt-token",
		User:  userResp,
	}

	mockService.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(loginResp, nil)

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "User registered successfully", result["message"])
	assert.NotNil(t, result["data"])
}

func TestHandler_CreateUser_InvalidBody(t *testing.T) {
	handler, _, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Post("/users", handler.CreateUser)

	httpReq, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHandler_CreateUser_ValidationError_ShortPassword(t *testing.T) {
	handler, _, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Post("/users", handler.CreateUser)

	req := userapi.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "123",
	}

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "Validation failed", result["message"])
	assert.Equal(t, "VALIDATION_ERROR", result["error_code"])
}

func TestHandler_CreateUser_ValidationError_ShortName(t *testing.T) {
	handler, _, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Post("/users", handler.CreateUser)

	req := userapi.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "AB",
		Password: "password123",
	}

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "Validation failed", result["message"])
	assert.Equal(t, "VALIDATION_ERROR", result["error_code"])
}

func TestHandler_CreateUser_ServiceError(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Post("/users", handler.CreateUser)

	req := userapi.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockService.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(nil, domain.ErrUserAlreadyExists)

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHandler_Login(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Post("/auth/login", handler.Login)

	req := userapi.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userResp := &response.UserResponse{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	loginResp := &response.LoginResponse{
		Token: "jwt-token",
		User:  userResp,
	}

	mockService.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(loginResp, nil)

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "Login successful", result["message"])
	assert.NotNil(t, result["data"])
}

func TestHandler_Login_InvalidBody(t *testing.T) {
	handler, _, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Post("/auth/login", handler.Login)

	httpReq, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHandler_Login_ValidationError_EmptyPassword(t *testing.T) {
	handler, _, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Post("/auth/login", handler.Login)

	req := userapi.LoginRequest{
		Email:    "test@example.com",
		Password: "",
	}

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "Validation failed", result["message"])
	assert.Equal(t, "VALIDATION_ERROR", result["error_code"])
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Post("/auth/login", handler.Login)

	req := userapi.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockService.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(nil, domain.ErrInvalidCredentials)

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestHandler_GetUserById(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return handler.GetUserById(c, openapi_types.UUID(userID))
	})

	userResp := &response.UserResponse{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(userResp, nil)

	httpReq, _ := http.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "User retrieved successfully", result["message"])
	assert.NotNil(t, result["data"])
}

func TestHandler_GetUserById_NotFound(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return handler.GetUserById(c, openapi_types.UUID(userID))
	})

	mockService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(nil, domain.ErrUserNotFound)

	httpReq, _ := http.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestHandler_UpdateUser(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	app.Put("/users/:id", func(c *fiber.Ctx) error {
		return handler.UpdateUser(c, openapi_types.UUID(userID))
	})

	req := userapi.UpdateUserRequest{
		Name: "Updated Name",
	}

	userResp := &response.UserResponse{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Updated Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.EXPECT().
		UpdateUser(gomock.Any(), userID, gomock.Any()).
		Return(userResp, nil)

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "User updated successfully", result["message"])
	assert.NotNil(t, result["data"])
}

func TestHandler_UpdateUser_InvalidBody(t *testing.T) {
	handler, _, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	app.Put("/users/:id", func(c *fiber.Ctx) error {
		return handler.UpdateUser(c, openapi_types.UUID(userID))
	})

	httpReq, _ := http.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHandler_UpdateUser_ValidationError_ShortName(t *testing.T) {
	handler, _, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	app.Put("/users/:id", func(c *fiber.Ctx) error {
		return handler.UpdateUser(c, openapi_types.UUID(userID))
	})

	req := userapi.UpdateUserRequest{
		Name: "AB",
	}

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "Validation failed", result["message"])
	assert.Equal(t, "VALIDATION_ERROR", result["error_code"])
}

func TestHandler_UpdateUser_ServiceError(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	app.Put("/users/:id", func(c *fiber.Ctx) error {
		return handler.UpdateUser(c, openapi_types.UUID(userID))
	})

	req := userapi.UpdateUserRequest{
		Name: "Updated Name",
	}

	mockService.EXPECT().
		UpdateUser(gomock.Any(), userID, gomock.Any()).
		Return(nil, errors.New("update failed"))

	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHandler_DeleteUser(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		return handler.DeleteUser(c, openapi_types.UUID(userID))
	})

	mockService.EXPECT().
		DeleteUser(gomock.Any(), userID).
		Return(nil)

	httpReq, _ := http.NewRequest(http.MethodDelete, "/users/"+userID.String(), nil)

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "User deleted successfully", result["message"])
}

func TestHandler_DeleteUser_NotFound(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		return handler.DeleteUser(c, openapi_types.UUID(userID))
	})

	mockService.EXPECT().
		DeleteUser(gomock.Any(), userID).
		Return(domain.ErrUserNotFound)

	httpReq, _ := http.NewRequest(http.MethodDelete, "/users/"+userID.String(), nil)

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestHandler_ListUsers(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	page := 1
	limit := 10

	app.Get("/users", func(c *fiber.Ctx) error {
		params := userapi.ListUsersParams{
			Page:  &page,
			Limit: &limit,
		}
		return handler.ListUsers(c, params)
	})

	users := []*response.UserResponse{
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

	mockService.EXPECT().
		ListUsers(gomock.Any(), page, limit).
		Return(users, int64(2), nil)

	httpReq, _ := http.NewRequest(http.MethodGet, "/users?page=1&limit=10", nil)

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.Equal(t, "Users retrieved successfully", result["message"])
	assert.NotNil(t, result["data"])
}

func TestHandler_ListUsers_DefaultPagination(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	app.Get("/users", func(c *fiber.Ctx) error {
		params := userapi.ListUsersParams{}
		return handler.ListUsers(c, params)
	})

	users := []*response.UserResponse{}

	mockService.EXPECT().
		ListUsers(gomock.Any(), 1, 10).
		Return(users, int64(0), nil)

	httpReq, _ := http.NewRequest(http.MethodGet, "/users", nil)

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHandler_ListUsers_InvalidPagination(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	invalidPage := -1
	invalidLimit := 200

	app.Get("/users", func(c *fiber.Ctx) error {
		params := userapi.ListUsersParams{
			Page:  &invalidPage,
			Limit: &invalidLimit,
		}
		return handler.ListUsers(c, params)
	})

	users := []*response.UserResponse{}

	// Should normalize to page=1, limit=10
	mockService.EXPECT().
		ListUsers(gomock.Any(), 1, 10).
		Return(users, int64(0), nil)

	httpReq, _ := http.NewRequest(http.MethodGet, "/users", nil)

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestHandler_ListUsers_ServiceError(t *testing.T) {
	handler, mockService, ctrl, app := setupHandlerTest(t)
	defer ctrl.Finish()

	page := 1
	limit := 10

	app.Get("/users", func(c *fiber.Ctx) error {
		params := userapi.ListUsersParams{
			Page:  &page,
			Limit: &limit,
		}
		return handler.ListUsers(c, params)
	})

	mockService.EXPECT().
		ListUsers(gomock.Any(), page, limit).
		Return(nil, int64(0), errors.New("database error"))

	httpReq, _ := http.NewRequest(http.MethodGet, "/users", nil)

	resp, err := app.Test(httpReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
