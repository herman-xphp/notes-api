package repository

import (
	"context"
	"notes-api/internal/domain"
)

type NoteRepository interface {
	Create(ctx context.Context, note *domain.Note) error
	Update(ctx context.Context, note *domain.Note) error
	Delete(ctx context.Context, id uint, userID uint) error

	FindByID(ctx context.Context, id uint, userID uint) (*domain.Note, error)
	FindAllByUser(ctx context.Context, userID uint) ([]domain.Note, error)
}
