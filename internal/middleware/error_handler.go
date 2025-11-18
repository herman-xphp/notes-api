package middleware

import (
	"errors"
	apperrors "notes-api/internal/errors"
	"notes-api/internal/handler"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ErrorHandler adalah global error handler untuk semua route
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	// Check for custom AppError
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		code = appErr.Code
		message = appErr.Message
	} else {
		// Check error type
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			message = fiberErr.Message
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			code = fiber.StatusNotFound
			message = "Resource not found"
		} else {
			// Check for validation errors
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				code = fiber.StatusBadRequest
				message = formatValidationErrors(validationErrors)
			} else {
				// Check for specific error messages that should return 401
				errMsg := err.Error()
				if errMsg == "invalid email or password" || errMsg == "unauthorized" {
					code = fiber.StatusUnauthorized
				}
				// Generic error message
				message = errMsg
			}
		}
	}

	// Use BaseResponse structure for consistency
	return c.Status(code).JSON(handler.BaseResponse{
		Status:  "error",
		Message: message,
	})
}

// formatValidationErrors formats validator errors into readable message
func formatValidationErrors(errs validator.ValidationErrors) string {
	if len(errs) == 0 {
		return "Validation failed"
	}

	// Return first error message
	err := errs[0]
	field := err.Field()
	tag := err.Tag()

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + err.Param() + " characters"
	case "max":
		return field + " must be at most " + err.Param() + " characters"
	default:
		return field + " is invalid"
	}
}
