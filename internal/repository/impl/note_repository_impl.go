package impl

import (
	"context"
	"notes-api/internal/domain"
	"notes-api/internal/repository"

	"gorm.io/gorm"
)

type noteRepositoryImpl struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) repository.NoteRepository {
	return &noteRepositoryImpl{db: db}
}

func (r *noteRepositoryImpl) Create(ctx context.Context, note *domain.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *noteRepositoryImpl) Update(ctx context.Context, note *domain.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

func (r *noteRepositoryImpl) Delete(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Note{}).
		Error
}

func (r *noteRepositoryImpl) FindByID(ctx context.Context, id uint, userID uint) (*domain.Note, error) {
	var note domain.Note
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&note).Error

	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noteRepositoryImpl) FindAllByUser(ctx context.Context, userID uint) ([]domain.Note, error) {
	var notes []domain.Note

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notes).Error

	if err != nil {
		return nil, err
	}

	return notes, nil
}
