package handler

import (
	"strconv"

	"github.com/gieart87/gohexaclean/internal/dto/request"
	"github.com/gieart87/gohexaclean/internal/port/inbound"
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService inbound.UserServicePort
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService inbound.UserServicePort) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles user creation
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "Create user request"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req request.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid request body", err),
		)
	}

	user, err := h.userService.CreateUser(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Failed to create user", err),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.NewSuccessResponse("User created successfully", user),
	)
}

// GetUser handles getting user by ID
// @Summary Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid user ID", err),
		)
	}

	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse("User not found", err),
		)
	}

	return c.JSON(
		response.NewSuccessResponse("User retrieved successfully", user),
	)
}

// UpdateUser handles user update
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body request.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid user ID", err),
		)
	}

	var req request.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid request body", err),
		)
	}

	user, err := h.userService.UpdateUser(c.Context(), id, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Failed to update user", err),
		)
	}

	return c.JSON(
		response.NewSuccessResponse("User updated successfully", user),
	)
}

// DeleteUser handles user deletion
// @Summary Delete user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid user ID", err),
		)
	}

	if err := h.userService.DeleteUser(c.Context(), id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse("Failed to delete user", err),
		)
	}

	return c.JSON(
		response.NewSuccessResponse("User deleted successfully", nil),
	)
}

// ListUsers handles listing users with pagination
// @Summary List users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userService.ListUsers(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse("Failed to list users", err),
		)
	}

	return c.JSON(
		response.NewPaginatedResponse("Users retrieved successfully", users, page, limit, total),
	)
}

// Login handles user login
// @Summary User login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login request"
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid request body", err),
		)
	}

	loginResp, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewErrorResponse("Invalid credentials", err),
		)
	}

	return c.JSON(
		response.NewSuccessResponse("Login successful", loginResp),
	)
}
