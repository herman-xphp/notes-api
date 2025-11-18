package handler

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"notes-api/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	registerResp *dto.AuthResponse
	registerErr  error
	loginResp    *dto.AuthResponse
	loginErr     error
}

func (m *mockUserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	return m.registerResp, m.registerErr
}

func (m *mockUserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	return m.loginResp, m.loginErr
}

// Test Register - Success Case
func TestRegisterSuccess(t *testing.T) {
	mock := &mockUserService{
		registerResp: &dto.AuthResponse{ID: 1, Name: "Test User", Email: "test@example.com", Token: "valid.jwt.token"},
		registerErr:  nil,
	}

	h := NewAuthHandler(mock)
	app := fiber.New()

	app.Post("/register", func(c *fiber.Ctx) error {
		c.Locals("validated", &dto.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "SecurePassword123",
		})
		return h.Register(c)
	})

	req := httptest.NewRequest("POST", "/register", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)
}

// Test Register - Email Already Exists
func TestRegisterEmailExists(t *testing.T) {
	mock := &mockUserService{
		registerResp: nil,
		registerErr:  errors.New("email already exists"),
	}

	h := NewAuthHandler(mock)
	app := fiber.New()

	app.Post("/register", func(c *fiber.Ctx) error {
		c.Locals("validated", &dto.RegisterRequest{
			Name:     "Test User",
			Email:    "existing@example.com",
			Password: "SecurePassword123",
		})
		return h.Register(c)
	})

	req := httptest.NewRequest("POST", "/register", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 409, resp.StatusCode) // Conflict
}

// Test Register - Internal Server Error
func TestRegisterServerError(t *testing.T) {
	mock := &mockUserService{
		registerResp: nil,
		registerErr:  errors.New("database error"),
	}

	h := NewAuthHandler(mock)
	app := fiber.New()

	app.Post("/register", func(c *fiber.Ctx) error {
		c.Locals("validated", &dto.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "SecurePassword123",
		})
		return h.Register(c)
	})

	req := httptest.NewRequest("POST", "/register", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 500, resp.StatusCode)
}

// Test Login - Success Case
func TestLoginSuccess(t *testing.T) {
	mock := &mockUserService{
		loginResp: &dto.AuthResponse{ID: 1, Name: "Test User", Email: "test@example.com", Token: "valid.jwt.token"},
		loginErr:  nil,
	}

	h := NewAuthHandler(mock)
	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		c.Locals("validated", &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "SecurePassword123",
		})
		return h.Login(c)
	})

	req := httptest.NewRequest("POST", "/login", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// Test Login - Invalid Credentials
func TestLoginUnauthorized(t *testing.T) {
	mock := &mockUserService{
		loginResp: nil,
		loginErr:  errors.New("invalid email or password"),
	}

	h := NewAuthHandler(mock)
	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		c.Locals("validated", &dto.LoginRequest{
			Email:    "wrong@example.com",
			Password: "wrongpassword",
		})
		return h.Login(c)
	})

	req := httptest.NewRequest("POST", "/login", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

// Test Login - Server Error
func TestLoginServerError(t *testing.T) {
	mock := &mockUserService{
		loginResp: nil,
		loginErr:  errors.New("database error"),
	}

	h := NewAuthHandler(mock)
	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		c.Locals("validated", &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "SecurePassword123",
		})
		return h.Login(c)
	})

	req := httptest.NewRequest("POST", "/login", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 500, resp.StatusCode)
}

// Test NewAuthHandler - Handler Initialization
func TestNewAuthHandler(t *testing.T) {
	mock := &mockUserService{}
	h := NewAuthHandler(mock)
	require.NotNil(t, h)
}
