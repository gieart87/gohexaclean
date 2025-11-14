package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/middleware"
	"github.com/gieart87/gohexaclean/internal/adapter/inbound/http/router"
	"github.com/gieart87/gohexaclean/internal/bootstrap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	configPath := getConfigPath()

	// Initialize container
	container, err := bootstrap.NewContainer(configPath)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      container.Config.App.Name,
		ServerHeader: "GoHexaClean",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.RecoveryMiddleware(container.Logger))
	app.Use(middleware.LoggerMiddleware(container.Logger))
	app.Use(middleware.CORSMiddleware(&container.Config.CORS))

	// Setup routes
	router.SetupRoutes(
		app,
		container.UserService,
		container.Config.JWT.Secret,
		container.Logger,
	)

	// Start server
	port := fmt.Sprintf(":%d", container.Config.Server.HTTP.Port)
	container.Logger.Info(fmt.Sprintf("HTTP Server starting on port %d", container.Config.Server.HTTP.Port))

	// Graceful shutdown
	go func() {
		if err := app.Listen(port); err != nil {
			container.Logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		container.Logger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	container.Logger.Info("Server exited")
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "config/app.yaml"
}

// customErrorHandler handles errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": "An error occurred",
		"error":   err.Error(),
	})
}
