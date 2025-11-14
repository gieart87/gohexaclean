package middleware

import (
	"strings"

	"github.com/gieart87/gohexaclean/pkg/auth"
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				response.NewErrorResponse("Missing authorization header", nil),
			)
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				response.NewErrorResponse("Invalid authorization header format", nil),
			)
		}

		token := parts[1]

		// Validate token
		claims, err := auth.ValidateJWT(token, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				response.NewErrorResponse("Invalid or expired token", err),
			)
		}

		// Store user ID in context
		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}
