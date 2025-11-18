package utils

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var logger *zerolog.Logger

func init() {
	// Initialize logger with console output
	// Pretty print in development, JSON in production
	var output io.Writer = os.Stderr

	// Check if running in development mode
	isDev := os.Getenv("ENV") != "production"

	if isDev {
		output = zerolog.ConsoleWriter{Out: os.Stderr}
	}

	l := zerolog.New(output).
		With().
		Timestamp().
		Logger()

	logger = &l

	// Set global level
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

// GetLogger returns the global logger instance
func GetLogger() *zerolog.Logger {
	return logger
}

// GetLoggerWithContext returns a logger with request context (request_id, user_id)
// This allows tracing requests through log files
func GetLoggerWithContext(ctx context.Context) *zerolog.Logger {
	l := logger.With()

	// Add request ID if available
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		l = l.Str("request_id", requestID)
	}

	// Add user ID if available
	if userID, ok := ctx.Value("user_id").(uint); ok && userID > 0 {
		l = l.Uint("user_id", userID)
	}

	contextLogger := l.Logger()
	return &contextLogger
}
