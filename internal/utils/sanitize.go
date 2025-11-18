package utils

import (
	"html"
	"strings"
)

// SanitizeString sanitizes a string by escaping HTML entities
// Prevents XSS attacks
func SanitizeString(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)
	// Escape HTML entities
	return html.EscapeString(trimmed)
}

// SanitizeHTML allows some HTML but removes dangerous tags
// For content that might need formatting, use a whitelist approach
func SanitizeHTML(input string) string {
	// For now, just escape everything
	// In production, consider using a library like bluemonday for HTML sanitization
	return html.EscapeString(strings.TrimSpace(input))
}

