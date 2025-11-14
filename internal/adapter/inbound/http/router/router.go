package router

import (
	"os"

	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/healthapi"
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated/userapi"
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/handler"
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/handler/health"
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/handler/user"
	"github.com/gieart87/gohexaclean/internal/infra/logger"
	"github.com/gieart87/gohexaclean/internal/port/inbound"
	"github.com/gieart87/gohexaclean/internal/port/outbound/telemetry"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all routes for the application using OpenAPI auto-generated routing
func SetupRoutes(
	app *fiber.App,
	userService inbound.UserServicePort,
	jwtSecret string,
	log *logger.Logger,
	metricsService telemetry.MetricsService,
	tracingService telemetry.TracingService,
) {
	// API v1 group
	api := app.Group("/api/v1")

	// Swagger documentation
	swaggerHandler := handler.NewSwaggerHandler()
	api.Get("/swagger", swaggerHandler.ServeSwaggerUI)
	api.Get("/swagger/spec", func(c *fiber.Ctx) error {
		// Read OpenAPI spec file
		spec, err := os.ReadFile("api/openapi/user-api.yaml")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load API specification",
			})
		}
		c.Set("Content-Type", "application/x-yaml")
		return c.Send(spec)
	})

	// Create health handler that implements healthapi.ServerInterface
	healthHandler := health.NewHandler()

	// Create user handler that implements userapi.ServerInterface
	userHandler := user.NewHandler(userService)

	// Auto-register health routes from OpenAPI spec
	// This will create: GET /health (public - health check)
	healthapi.RegisterHandlers(api, healthHandler)

	// Auto-register user routes from OpenAPI spec
	// This will create routes for:
	// - POST /auth/login (public - login)
	// - POST /users (public - register)
	// - GET /users (protected - list users)
	// - GET /users/{id} (protected - get user)
	// - PUT /users/{id} (protected - update user)
	// - DELETE /users/{id} (protected - delete user)
	userapi.RegisterHandlers(api, userHandler)

	// Note: For protected routes, you'll need to add auth middleware
	// This can be done by creating a custom wrapper or using middleware in specific routes
}
