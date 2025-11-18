package impl

import (
	"context"
	"errors"
	apperrors "notes-api/internal/errors"
	"notes-api/internal/domain"
	"notes-api/internal/dto"
	"notes-api/internal/repository"
	"notes-api/internal/service"
	"notes-api/internal/utils"

	"gorm.io/gorm"
)

type noteServiceImpl struct {
	noteRepo repository.NoteRepository
}

func NewNoteServiceImpl(noteRepo repository.NoteRepository) service.NoteService {
	return &noteServiceImpl{noteRepo: noteRepo}
}

func (s *noteServiceImpl) Create(ctx context.Context, userID uint, req dto.CreateNoteRequest) (*domain.Note, error) {
	// Sanitize input to prevent XSS
	sanitizedTitle := utils.SanitizeString(req.Title)
	sanitizedContent := utils.SanitizeHTML(req.Content)

	note := domain.Note{
		UserID:  userID,
		Title:   sanitizedTitle,
		Content: sanitizedContent,
	}

	if err := s.noteRepo.Create(ctx, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (s *noteServiceImpl) Update(ctx context.Context, userID uint, noteID uint, req dto.UpdateNoteRequest) (*domain.Note, error) {
	note, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	if req.Title != nil {
		note.Title = utils.SanitizeString(*req.Title)
	}
	if req.Content != nil {
		note.Content = utils.SanitizeHTML(*req.Content)
	}

	if err := s.noteRepo.Update(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *noteServiceImpl) Delete(ctx context.Context, userID uint, noteID uint) error {
	// Check if note exists first
	_, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}
		return err
	}

	return s.noteRepo.Delete(ctx, noteID, userID)
}

func (s *noteServiceImpl) GetByID(ctx context.Context, userID uint, noteID uint) (*domain.Note, error) {
	note, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return note, nil
}

func (s *noteServiceImpl) GetAll(ctx context.Context, userID uint, page, limit int) ([]domain.Note, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max limit to prevent excessive data retrieval
	}

	return s.noteRepo.FindAllByUser(ctx, userID, page, limit)
}
