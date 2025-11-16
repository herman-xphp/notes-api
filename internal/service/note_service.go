package service

import (
	"context"
	"notes-api/internal/domain"
	"notes-api/internal/dto"
)

type NoteService interface {
	Create(ctx context.Context, userID uint, req dto.CreateNoteRequest) (*domain.Note, error)
	Update(ctx context.Context, userID uint, noteID uint, req dto.UpdateNoteRequest) (*domain.Note, error)
	Delete(ctx context.Context, userID uint, noteID uint) error
	GetByID(ctx context.Context, userID uint, noteID uint) (*domain.Note, error)
	GetAll(ctx context.Context, userID uint) ([]domain.Note, error)
}
