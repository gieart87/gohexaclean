package user

import (
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/userapi"
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// ListUsers handles listing users with pagination
// Protected endpoint - requires authentication
// GET /users
func (h *Handler) ListUsers(c *fiber.Ctx, params userapi.ListUsersParams) error {
	page := 1
	if params.Page != nil {
		page = *params.Page
	}

	limit := 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	// Validate pagination params
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
