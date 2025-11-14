package health

import (
	"github.com/gofiber/fiber/v2"
)

// HealthCheck handles health check endpoint
// Public endpoint - no authentication required
// GET /health
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Service is running",
	})
}
