package user

import (
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// DeleteUser handles user deletion
// Protected endpoint - requires authentication
// DELETE /users/{id}
func (h *Handler) DeleteUser(c *fiber.Ctx, id openapi_types.UUID) error {
	if err := h.userService.DeleteUser(c.Context(), uuid.UUID(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse("Failed to delete user", err),
		)
	}

	return c.JSON(
		response.NewSuccessResponse("User deleted successfully", nil),
	)
}
