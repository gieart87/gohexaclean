package user

import (
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/userapi"
	"github.com/gieart87/gohexaclean/internal/dto/request"
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// CreateUser handles user registration
// Public endpoint - no authentication required
// POST /users
func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var req userapi.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid request body", err),
		)
	}

	// Convert generated type to domain DTO
	registerResp, err := h.userService.CreateUser(c.Context(), &request.CreateUserRequest{
		Email:    string(req.Email),
		Name:     req.Name,
		Password: req.Password,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Failed to create user", err),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.NewSuccessResponse("User registered successfully", registerResp),
	)
}
