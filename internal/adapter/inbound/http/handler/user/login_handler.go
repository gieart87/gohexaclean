package user

import (
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/userapi"
	"github.com/gieart87/gohexaclean/internal/dto/request"
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// Login handles user login
// Public endpoint - no authentication required
// POST /auth/login
func (h *Handler) Login(c *fiber.Ctx) error {
	var req userapi.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid request body", err),
		)
	}

	// Convert generated type to domain DTO
	loginResp, err := h.userService.Login(c.Context(), &request.LoginRequest{
		Email:    string(req.Email),
		Password: req.Password,
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewErrorResponse("Invalid credentials", err),
		)
	}

	return c.JSON(
		response.NewSuccessResponse("Login successful", loginResp),
	)
}
