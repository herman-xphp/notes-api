package impl

import (
	"context"
	"notes-api/internal/domain"
	"notes-api/internal/dto"
	"notes-api/internal/repository"
	"notes-api/internal/service"
)

type noteServiceImpl struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) service.NoteService {
	return &noteServiceImpl{noteRepo: noteRepo}
}

func (s *noteServiceImpl) Create(ctx context.Context, userID uint, req dto.CreateNoteRequest) (*domain.Note, error) {
	note := domain.Note{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := s.noteRepo.Create(ctx, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (s *noteServiceImpl) Update(ctx context.Context, userID uint, noteID uint, req dto.UpdateNoteRequest) (*domain.Note, error) {
	note, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		note.Title = *req.Title
	}
	if req.Content != nil {
		note.Content = *req.Content
	}

	if err := s.noteRepo.Update(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *noteServiceImpl) Delete(ctx context.Context, userID uint, noteID uint) error {
	return s.noteRepo.Delete(ctx, noteID, userID)
}

func (s *noteServiceImpl) GetByID(ctx context.Context, userID uint, noteID uint) (*domain.Note, error) {
	return s.noteRepo.FindByID(ctx, noteID, userID)
}

func (s *noteServiceImpl) GetAll(ctx context.Context, userID uint) ([]domain.Note, error) {
	return s.noteRepo.FindAllByUser(ctx, userID)
}
