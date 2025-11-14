package router

import (
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/handler"
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/middleware"
	"github.com/gieart87/gohexaclean/internal/infra/logger"
	"github.com/gieart87/gohexaclean/internal/port/outbound/service"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(
	app *fiber.App,
	userHandler *handler.UserHandler,
	tokenService service.TokenService,
	log *logger.Logger,
) {
	// API v1 group
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.Post("/login", userHandler.Login)
	}

	// User routes
	users := api.Group("/users")
	{
		// Public routes
		users.Post("/", userHandler.CreateUser)

		// Protected routes
		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware(tokenService))
		{
			protected.Get("/", userHandler.ListUsers)
			protected.Get("/:id", userHandler.GetUser)
			protected.Put("/:id", userHandler.UpdateUser)
			protected.Delete("/:id", userHandler.DeleteUser)
		}
	}
}
