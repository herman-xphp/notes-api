package handler

import (
	"notes-api/internal/constants"
	"notes-api/internal/dto"
	"notes-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(ns service.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: ns,
	}
}

// Create godoc
// @Summary Create a new note
// @Description Create a new note with title and content
// @Tags Notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateNoteRequest true "Create note request"
// @Success 201 {object} BaseResponse
// @Failure 400 {object} BaseResponse
// @Failure 401 {object} BaseResponse
// @Router /notes [post]
func (h *NoteHandler) Create(c *fiber.Ctx) error {
	// Get validated DTO from middleware
	body := c.Locals("validated").(*dto.CreateNoteRequest)

	// userID extracted from auth middleware
	userID, ok := c.Locals(constants.ContextKeyUserID).(uint)
	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	note, err := h.noteService.Create(c.Context(), userID, *body)
	if err != nil {
		return err
	}

	// Convert to response DTO
	response := dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	return JSONCreated(c, "note created", response)
}

// GetAll godoc
// @Summary Get all notes
// @Description Get all notes for the authenticated user with pagination
// @Tags Notes
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} BaseResponse
// @Failure 400 {object} BaseResponse
// @Failure 401 {object} BaseResponse
// @Router /notes [get]
func (h *NoteHandler) GetAll(c *fiber.Ctx) error {
	userID, ok := c.Locals(constants.ContextKeyUserID).(uint)
	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	// Parse pagination params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	// Validation moved to service layer for consistency
	notes, total, err := h.noteService.GetAll(c.Context(), userID, page, limit)
	if err != nil {
		return err
	}

	// Convert to response DTOs
	responses := make([]dto.NoteResponse, len(notes))
	for i, note := range notes {
		responses[i] = dto.NoteResponse{
			ID:        note.ID,
			Title:     note.Title,
			Content:   note.Content,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		}
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	// Create paginated response
	paginatedResponse := dto.PaginatedResponse{
		Data: responses,
		Meta: dto.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}

	return JSONSuccess(c, "notes fetched", paginatedResponse)
}

// GetByID godoc
// @Summary Get a note by ID
// @Description Get a specific note by ID for the authenticated user
// @Tags Notes
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Success 200 {object} BaseResponse
// @Failure 400 {object} BaseResponse
// @Failure 401 {object} BaseResponse
// @Failure 404 {object} BaseResponse
// @Router /notes/{id} [get]
func (h *NoteHandler) GetByID(c *fiber.Ctx) error {
	userID, ok := c.Locals(constants.ContextKeyUserID).(uint)
	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "invalid note id")
	}

	note, err := h.noteService.GetByID(c.Context(), userID, uint(id))
	if err != nil {
		return err
	}

	// Convert to response DTO
	response := dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	return JSONSuccess(c, "note fetched", response)
}

// Update godoc
// @Summary Update a note
// @Description Update a note by ID for the authenticated user
// @Tags Notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Param request body dto.UpdateNoteRequest true "Update note request"
// @Success 200 {object} BaseResponse
// @Failure 400 {object} BaseResponse
// @Failure 401 {object} BaseResponse
// @Failure 404 {object} BaseResponse
// @Router /notes/{id} [put]
func (h *NoteHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals(constants.ContextKeyUserID).(uint)
	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "invalid note id")
	}

	// Get validated DTO from middleware
	body := c.Locals("validated").(*dto.UpdateNoteRequest)

	// Validate at least one field is provided
	if body.Title == nil && body.Content == nil {
		return JSONError(c, fiber.StatusBadRequest, "at least one field (title or content) must be provided")
	}

	note, err := h.noteService.Update(c.Context(), userID, uint(id), *body)
	if err != nil {
		return err
	}

	// Convert to response DTO
	response := dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	return JSONSuccess(c, "note updated", response)
}

// Delete godoc
// @Summary Delete a note
// @Description Delete a note by ID for the authenticated user
// @Tags Notes
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Success 200 {object} BaseResponse
// @Failure 400 {object} BaseResponse
// @Failure 401 {object} BaseResponse
// @Failure 404 {object} BaseResponse
// @Router /notes/{id} [delete]
func (h *NoteHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals(constants.ContextKeyUserID).(uint)
	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "invalid note id")
	}

	if err := h.noteService.Delete(c.Context(), userID, uint(id)); err != nil {
		return err
	}

	return JSONSuccess(c, "note deleted", nil)
}
