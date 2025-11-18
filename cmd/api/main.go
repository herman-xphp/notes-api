package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "notes-api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"notes-api/configs"
	"notes-api/internal/handler"
	"notes-api/internal/middleware"
	repoImpl "notes-api/internal/repository/impl"
	"notes-api/internal/router"
	serviceImpl "notes-api/internal/service/impl"
	"notes-api/internal/utils"
	"notes-api/pkg/database"
)

// @title Notes API
// @version 1.0
// @description A simple and secure Notes API with user authentication and CRUD operations
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @basePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load config
	cfg := configs.LoadConfigOrFatal()

	// Initialize JWT (must be after config load)
	if err := utils.InitJWTWithSecret(cfg.JWT.Secret); err != nil {
		utils.GetLogger().Fatal().Err(err).Msg("Failed to initialize JWT")
	}

	// Connect DB
	db := database.ConnectMySQLWithParams(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)
	if db == nil {
		utils.GetLogger().Fatal().Msg("Failed to connect database")
	}

	// === Init Repository ===
	userRepo := repoImpl.NewUserRepositoryImpl(db)
	noteRepo := repoImpl.NewNoteRepositoryImpl(db)

	// === Init Service ===
	userService := serviceImpl.NewUserServiceImpl(userRepo)
	noteService := serviceImpl.NewNoteServiceImpl(noteRepo)

	// === Init Handler ===
	authHandler := handler.NewAuthHandler(userService)
	noteHandler := handler.NewNoteHandler(noteService)

	// === Start Fiber ===
	app := fiber.New(fiber.Config{
		ErrorHandler:     middleware.ErrorHandler,
		ReadTimeout:      10 * time.Second,
		WriteTimeout:     10 * time.Second,
		IdleTimeout:      120 * time.Second,
		BodyLimit:        4 * 1024 * 1024, // 4MB body size limit
		DisableKeepalive: false,
	})

	// Request ID middleware (should be early in the chain)
	app.Use(middleware.RequestIDMiddleware)

	// Security headers middleware
	app.Use(middleware.SecurityHeadersMiddleware)

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}))

	// Setup routes
	router.Setup(app, db, authHandler, noteHandler)

	// Start server
	port := cfg.App.Port

	// Graceful shutdown
	go func() {
		utils.GetLogger().Info().Str("port", port).Msg("Server running")
		if err := app.Listen(":" + port); err != nil {
			utils.GetLogger().Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	utils.GetLogger().Info().Msg("Shutting down server")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown Fiber app
	if err := app.ShutdownWithContext(ctx); err != nil {
		utils.GetLogger().Fatal().Err(err).Msg("Server forced to shutdown")
	}

	// Close database connections
	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			utils.GetLogger().Warn().Err(err).Msg("Error closing database")
		} else {
			utils.GetLogger().Info().Msg("Database connections closed")
		}
	}

	utils.GetLogger().Info().Msg("Server exited gracefully")
}
