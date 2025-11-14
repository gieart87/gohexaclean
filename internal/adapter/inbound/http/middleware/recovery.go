package middleware

import (
	"github.com/gieart87/gohexaclean/internal/infra/logger"
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RecoveryMiddleware creates a panic recovery middleware
func RecoveryMiddleware(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic recovered",
					zap.Any("panic", r),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
				)

				c.Status(fiber.StatusInternalServerError).JSON(
					response.NewErrorResponse("Internal server error", nil),
				)
			}
		}()

		return c.Next()
	}
}
