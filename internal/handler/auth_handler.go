package handler

import (
	"notes-api/internal/dto"
	"notes-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(us service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: us,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} BaseResponse
// @Failure 400 {object} BaseResponse
// @Failure 409 {object} BaseResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Get validated DTO from middleware
	body := c.Locals("validated").(*dto.RegisterRequest)

	user, err := h.userService.Register(c.Context(), *body)
	if err != nil {
		if err.Error() == "email already exists" {
			return JSONError(c, fiber.StatusConflict, err.Error())
		}
		return err
	}

	return JSONCreated(c, "user registered successfully", user)
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} BaseResponse
// @Failure 400 {object} BaseResponse
// @Failure 401 {object} BaseResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Get validated DTO from middleware
	body := c.Locals("validated").(*dto.LoginRequest)

	authResp, err := h.userService.Login(c.Context(), *body)
	if err != nil {
		// Return 401 for authentication errors
		if err.Error() == "invalid email or password" {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return err
	}

	return JSONSuccess(c, "login successful", authResp)
}
