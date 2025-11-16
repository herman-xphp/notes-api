package repository

import (
	"context"
	"notes-api/internal/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}
