package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

// RateLimitMiddleware creates a rate limiting middleware
// Limits requests to prevent brute force attacks
func RateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,                // Maximum 10 requests
		Expiration: 1 * time.Minute,   // Per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "too many requests, please try again later",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

// StrictRateLimitMiddleware creates a stricter rate limit for auth endpoints
// More restrictive to prevent brute force attacks on login/register
func StrictRateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,                 // Maximum 5 requests
		Expiration: 15 * time.Minute,  // Per 15 minutes
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "too many requests, please try again later",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

