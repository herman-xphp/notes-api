package middleware

import (
	"notes-api/internal/constants"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates a unique request ID for each request
// Useful for tracing requests in logs
func RequestIDMiddleware(c *fiber.Ctx) error {
	// Check if request ID already exists in header
	requestID := c.Get("X-Request-ID")
	if requestID == "" {
		// Generate new UUID
		requestID = uuid.New().String()
	}

	// Set request ID in response header
	c.Set("X-Request-ID", requestID)

	// Store in locals using constant key
	c.Locals(constants.ContextKeyRequestID, requestID)

	return c.Next()
}
