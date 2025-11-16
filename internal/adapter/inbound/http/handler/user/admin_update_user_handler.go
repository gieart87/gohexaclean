package user

import (
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/userapi"
	"github.com/gieart87/gohexaclean/internal/dto/request"
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// UpdateUser handles user update
// Protected endpoint - requires authentication
// PUT /users/{id}
func (h *Handler) UpdateUser(c *fiber.Ctx, id openapi_types.UUID) error {
	var req userapi.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid request body", err),
		)
	}

	// Convert generated type to domain DTO
	updateReq := &request.UpdateUserRequest{
		Name: req.Name,
	}

	// Validate request
	if err := updateReq.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewValidationErrorResponse("Validation failed", response.ParseValidationErrors(err)),
		)
	}

	user, err := h.userService.UpdateUser(c.Context(), uuid.UUID(id), updateReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Failed to update user", err),
		)
	}

	return c.JSON(
		response.NewSuccessResponse("User updated successfully", user),
	)
}
