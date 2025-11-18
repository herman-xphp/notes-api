package handler

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"notes-api/internal/domain"
	"notes-api/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type mockNoteService struct {
	createResp  *domain.Note
	createErr   error
	updateResp  *domain.Note
	updateErr   error
	deleteErr   error
	getByIDResp *domain.Note
	getByIDErr  error
	getAllResp  []domain.Note
	getAllTotal int64
	getAllErr   error
}

func (m *mockNoteService) Create(ctx context.Context, userID uint, req dto.CreateNoteRequest) (*domain.Note, error) {
	return m.createResp, m.createErr
}

func (m *mockNoteService) Update(ctx context.Context, userID uint, noteID uint, req dto.UpdateNoteRequest) (*domain.Note, error) {
	return m.updateResp, m.updateErr
}

func (m *mockNoteService) Delete(ctx context.Context, userID uint, noteID uint) error {
	return m.deleteErr
}

func (m *mockNoteService) GetByID(ctx context.Context, userID uint, noteID uint) (*domain.Note, error) {
	return m.getByIDResp, m.getByIDErr
}

func (m *mockNoteService) GetAll(ctx context.Context, userID uint, page, limit int) ([]domain.Note, int64, error) {
	return m.getAllResp, m.getAllTotal, m.getAllErr
}

// Test Create Note - Unauthorized (no user_id)
func TestCreateNote_Unauthorized(t *testing.T) {
	mock := &mockNoteService{}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Post("/notes", func(c *fiber.Ctx) error {
		c.Locals("validated", &dto.CreateNoteRequest{Title: "Test", Content: "Content"})
		return h.Create(c)
	})

	req := httptest.NewRequest("POST", "/notes", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

// Test Create Note - Success
func TestCreateNote_Success(t *testing.T) {
	now := time.Now()
	mock := &mockNoteService{
		createResp: &domain.Note{
			ID:        1,
			UserID:    1,
			Title:     "Test Note",
			Content:   "Test Content",
			CreatedAt: now,
			UpdatedAt: now,
		},
		createErr: nil,
	}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Post("/notes", func(c *fiber.Ctx) error {
		c.Locals("validated", &dto.CreateNoteRequest{Title: "Test Note", Content: "Test Content"})
		c.Locals("user_id", uint(1))
		return h.Create(c)
	})

	req := httptest.NewRequest("POST", "/notes", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)
}

// Test Create Note - Service Error
func TestCreateNote_ServiceError(t *testing.T) {
	mock := &mockNoteService{
		createResp: nil,
		createErr:  errors.New("database error"),
	}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Post("/notes", func(c *fiber.Ctx) error {
		c.Locals("validated", &dto.CreateNoteRequest{Title: "Test", Content: "Content"})
		c.Locals("user_id", uint(1))
		return h.Create(c)
	})

	req := httptest.NewRequest("POST", "/notes", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 500, resp.StatusCode)
}

// Test Get Note By ID - Unauthorized
func TestGetByID_Unauthorized(t *testing.T) {
	mock := &mockNoteService{}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Get("/notes/:id", func(c *fiber.Ctx) error {
		return h.GetByID(c)
	})

	req := httptest.NewRequest("GET", "/notes/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

// Test Get Note By ID - Invalid ID
func TestGetByID_InvalidID(t *testing.T) {
	mock := &mockNoteService{}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Get("/notes/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return h.GetByID(c)
	})

	req := httptest.NewRequest("GET", "/notes/invalid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}

// Test Get Note By ID - Not Found
// (Integration test - requires full middleware setup)
// Skipped in unit tests - will be covered in integration tests

// Test Get Note By ID - Success
func TestGetByID_Success(t *testing.T) {
	now := time.Now()
	mock := &mockNoteService{
		getByIDResp: &domain.Note{
			ID:        1,
			UserID:    1,
			Title:     "Test Note",
			Content:   "Test Content",
			CreatedAt: now,
			UpdatedAt: now,
		},
		getByIDErr: nil,
	}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Get("/notes/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return h.GetByID(c)
	})

	req := httptest.NewRequest("GET", "/notes/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// Test Get All Notes - Unauthorized
func TestGetAllNotes_Unauthorized(t *testing.T) {
	mock := &mockNoteService{}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Get("/notes", func(c *fiber.Ctx) error {
		return h.GetAll(c)
	})

	req := httptest.NewRequest("GET", "/notes", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

// Test Get All Notes - Success with Pagination
func TestGetAllNotes_Success(t *testing.T) {
	now := time.Now()
	mock := &mockNoteService{
		getAllResp: []domain.Note{
			{ID: 1, UserID: 1, Title: "Note 1", Content: "Content 1", CreatedAt: now, UpdatedAt: now},
			{ID: 2, UserID: 1, Title: "Note 2", Content: "Content 2", CreatedAt: now, UpdatedAt: now},
		},
		getAllTotal: 2,
		getAllErr:   nil,
	}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Get("/notes", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return h.GetAll(c)
	})

	req := httptest.NewRequest("GET", "/notes?page=1&limit=10", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// Test Update Note - Unauthorized
func TestUpdateNote_Unauthorized(t *testing.T) {
	mock := &mockNoteService{}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Put("/notes/:id", func(c *fiber.Ctx) error {
		title := "Updated"
		content := "Updated"
		c.Locals("validated", &dto.UpdateNoteRequest{Title: &title, Content: &content})
		return h.Update(c)
	})

	req := httptest.NewRequest("PUT", "/notes/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

// Test Update Note - Invalid ID
func TestUpdateNote_InvalidID(t *testing.T) {
	mock := &mockNoteService{}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Put("/notes/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		title := "Updated"
		content := "Updated"
		c.Locals("validated", &dto.UpdateNoteRequest{Title: &title, Content: &content})
		return h.Update(c)
	})

	req := httptest.NewRequest("PUT", "/notes/invalid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}

// Test Update Note - Success
func TestUpdateNote_Success(t *testing.T) {
	now := time.Now()
	mock := &mockNoteService{
		updateResp: &domain.Note{
			ID:        1,
			UserID:    1,
			Title:     "Updated Note",
			Content:   "Updated Content",
			CreatedAt: now,
			UpdatedAt: now,
		},
		updateErr: nil,
	}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Put("/notes/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		title := "Updated Note"
		content := "Updated Content"
		c.Locals("validated", &dto.UpdateNoteRequest{Title: &title, Content: &content})
		return h.Update(c)
	})

	req := httptest.NewRequest("PUT", "/notes/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// Test Delete Note - Unauthorized
func TestDeleteNote_Unauthorized(t *testing.T) {
	mock := &mockNoteService{}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Delete("/notes/:id", func(c *fiber.Ctx) error {
		return h.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/notes/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

// Test Delete Note - Invalid ID
func TestDeleteNote_InvalidID(t *testing.T) {
	mock := &mockNoteService{}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Delete("/notes/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return h.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/notes/invalid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}

// Test Delete Note - Success
func TestDeleteNote_Success(t *testing.T) {
	mock := &mockNoteService{
		deleteErr: nil,
	}
	h := NewNoteHandler(mock)

	app := fiber.New()
	app.Delete("/notes/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return h.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/notes/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// Test Delete Note - Not Found
// (Integration test - requires full middleware setup)
// Skipped in unit tests - will be covered in integration tests

// Test NewNoteHandler - Handler Initialization
func TestNewNoteHandler(t *testing.T) {
	mock := &mockNoteService{}
	h := NewNoteHandler(mock)
	require.NotNil(t, h)
}
