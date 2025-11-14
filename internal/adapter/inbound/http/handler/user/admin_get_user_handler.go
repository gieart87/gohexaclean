package user

import (
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// GetUserById handles getting user by ID
// Protected endpoint - requires authentication
// GET /users/{id}
func (h *Handler) GetUserById(c *fiber.Ctx, id openapi_types.UUID) error {
	user, err := h.userService.GetUserByID(c.Context(), uuid.UUID(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse("User not found", err),
		)
	}

	return c.JSON(
		response.NewSuccessResponse("User retrieved successfully", user),
	)
}
