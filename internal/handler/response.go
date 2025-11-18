package handler

import "github.com/gofiber/fiber/v2"

// BaseResponse adalah struktur response standard untuk semua endpoint
type BaseResponse struct {
	Status  string      `json:"status"`  // "success" atau "error"
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSONSuccess mengembalikan response sukses dengan status 200
func JSONSuccess(c *fiber.Ctx, msg string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(BaseResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}

// JSONCreated mengembalikan response sukses dengan status 201
func JSONCreated(c *fiber.Ctx, msg string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(BaseResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}

// JSONError mengembalikan response error
func JSONError(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(BaseResponse{
		Status:  "error",
		Message: msg,
	})
}
