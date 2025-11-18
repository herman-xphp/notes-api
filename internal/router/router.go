package router

import (
	"notes-api/internal/dto"
	"notes-api/internal/handler"
	"notes-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, authHandler *handler.AuthHandler, noteHandler *handler.NoteHandler) {
	// Health check handler
	healthHandler := handler.NewHealthCheckHandler(db)
	
	// Health check (no prefix)
	app.Get("/health", healthHandler.HealthCheck)

	// API routes with versioning
	api := app.Group("/api/v1")

	// Auth routes with strict rate limiting
	auth := api.Group("/auth")
	auth.Post("/register", middleware.StrictRateLimitMiddleware(), middleware.ValidateRequest(&dto.RegisterRequest{}), authHandler.Register)
	auth.Post("/login", middleware.StrictRateLimitMiddleware(), middleware.ValidateRequest(&dto.LoginRequest{}), authHandler.Login)

	// Notes routes (protected) with rate limiting
	notes := api.Group("/notes", middleware.AuthMiddleware, middleware.RateLimitMiddleware())
	notes.Post("/", middleware.ValidateRequest(&dto.CreateNoteRequest{}), noteHandler.Create)
	notes.Get("/", noteHandler.GetAll)
	notes.Get("/:id", noteHandler.GetByID)
	notes.Put("/:id", middleware.ValidateRequest(&dto.UpdateNoteRequest{}), noteHandler.Update)
	notes.Delete("/:id", noteHandler.Delete)
}
