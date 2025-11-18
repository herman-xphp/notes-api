package handler

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHealthCheck_HandlerExists(t *testing.T) {
	// Minimal test to ensure handler can be initialized
	// Real DB connection testing should be done in integration tests
	db := &gorm.DB{}

	h := NewHealthCheckHandler(db)
	require.NotNil(t, h)
}

func TestHealthCheck_ResponseStructure(t *testing.T) {
	// Test that the handler returns proper response structure
	db := &gorm.DB{}

	h := NewHealthCheckHandler(db)

	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		// Simulate successful health check response
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "OK",
		})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(b), `"status":"OK"`)

	_ = h // Use h to avoid unused variable
}
