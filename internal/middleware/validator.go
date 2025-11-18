package middleware

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	// Register custom tag name function to use json tag names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			name = fld.Name
		}
		return name
	})
}

// ValidateRequest validates request body against struct tags
func ValidateRequest(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create new instance of DTO type
		dtoType := reflect.TypeOf(dto)
		if dtoType.Kind() == reflect.Ptr {
			dtoType = dtoType.Elem()
		}
		newDto := reflect.New(dtoType).Interface()

		// Parse body
		if err := c.BodyParser(newDto); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate
		if err := validate.Struct(newDto); err != nil {
			return err
		}

		// Store validated DTO in context
		c.Locals("validated", newDto)
		return c.Next()
	}
}
