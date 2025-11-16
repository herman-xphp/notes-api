package main

import (
	"log"
	"notes-api/configs"
	"notes-api/internal/domain"
	"notes-api/internal/router"
	"notes-api/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load config
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("failed load config:", err)
	}

	// Init DB
	db, err := database.NewMySQLDB(cfg)
	if err != nil {
		log.Fatal("failed connect db:", err)
	}

	db.AutoMigrate(&domain.User{}, &domain.Note{})

	log.Println("Database connected:", db != nil)

	// Init Fiber
	app := fiber.New()

	// Setup router
	router.Setup(app)

	// Start server
	log.Println("Server running at port", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
