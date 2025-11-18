package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HealthCheckHandler handles health check requests
type HealthCheckHandler struct {
	db *gorm.DB
}

// NewHealthCheckHandler creates a new health check handler
func NewHealthCheckHandler(db *gorm.DB) *HealthCheckHandler {
	return &HealthCheckHandler{db: db}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check API and database health status
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health [get]
func (h *HealthCheckHandler) HealthCheck(ctx *fiber.Ctx) error {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "unhealthy",
			"message": "database connection error",
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "unhealthy",
			"message": "database ping failed",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "OK",
	})
}
