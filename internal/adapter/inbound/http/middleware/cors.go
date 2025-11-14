package middleware

import (
	"github.com/gieart87/gohexaclean/internal/infra/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSMiddleware creates a CORS middleware
func CORSMiddleware(cfg *config.CORSConfig) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     joinStrings(cfg.AllowOrigins, ","),
		AllowMethods:     joinStrings(cfg.AllowMethods, ","),
		AllowHeaders:     joinStrings(cfg.AllowHeaders, ","),
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           300,
	})
}

// joinStrings joins string slice with separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
