package middleware

import (
	"time"

	"github.com/gieart87/gohexaclean/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// LoggerMiddleware creates a logging middleware
func LoggerMiddleware(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details
		log.Info("HTTP Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		if err != nil {
			log.Error("Request error",
				zap.Error(err),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
		}

		return err
	}
}
