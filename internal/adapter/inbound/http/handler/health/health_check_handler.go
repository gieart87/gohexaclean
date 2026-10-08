package health

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck handles health check endpoint
// Public endpoint - no authentication required
// GET /health
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"service":        h.serviceName,
		"status":         "ok",
		"version":        h.version,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
		"uptime_seconds": int(time.Since(h.startedAt).Seconds()),
	})
}
