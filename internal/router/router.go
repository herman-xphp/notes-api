package router

import (
	"notes-api/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Health check
	app.Get("/health", handler.HealthCheck)
}
