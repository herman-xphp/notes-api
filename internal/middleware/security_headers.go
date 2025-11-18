package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware(c *fiber.Ctx) error {
	// Prevent MIME type sniffing
	c.Set("X-Content-Type-Options", "nosniff")
	
	// Prevent clickjacking
	c.Set("X-Frame-Options", "DENY")
	
	// XSS Protection (legacy, but still useful)
	c.Set("X-XSS-Protection", "1; mode=block")
	
	// Strict Transport Security (HSTS) - only set if using HTTPS
	// Uncomment if deploying with HTTPS
	// c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	
	// Content Security Policy
	c.Set("Content-Security-Policy", "default-src 'self'")
	
	// Referrer Policy
	c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	
	// Permissions Policy (formerly Feature Policy)
	c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	
	return c.Next()
}

