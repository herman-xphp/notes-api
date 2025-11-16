package main

import (
	"log"
	"notes-api/configs"
	"notes-api/internal/domain"
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

	// Default health check
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"message": "Notes API Running"})
	})

	log.Println("Server running at port", cfg.AppPort)
	app.Listen(":" + cfg.AppPort)
}
