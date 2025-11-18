package middleware

import (
	"strings"

	"notes-api/internal/constants"
	"notes-api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "missing authorization header",
		})
	}

	// Format: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid authorization format",
		})
	}

	token := parts[1]

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid token",
		})
	}

	// Store userID in context for handlers using constant key
	c.Locals(constants.ContextKeyUserID, claims.UserID)

	return c.Next()
}
