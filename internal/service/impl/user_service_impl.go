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
	"strings"

	"gorm.io/gorm"
)

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserServiceImpl(userRepo repository.UserRepository) service.UserService {
	return &userServiceImpl{userRepo: userRepo}
}

// Register user
func (s *userServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Normalize email: lowercase and trim whitespace
	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))
	
	// Check if email already exists
	existing, _ := s.userRepo.FindByEmail(ctx, normalizedEmail)
	if existing != nil {
		// Generic error message to prevent email enumeration
		return nil, apperrors.ErrInvalidCredentials
	}

	// Validate password strength
	plainPassword := strings.TrimSpace(req.Password)
	if !utils.ValidatePasswordStrength(plainPassword) {
		return nil, apperrors.ErrWeakPassword
	}

	// Hash password
	hash, err := utils.HashPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    normalizedEmail,
		Password: hash,
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}

func (s *userServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Normalize email: lowercase and trim whitespace
	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))
	
	user, err := s.userRepo.FindByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	// Compare password
	// Parameter: (hashedPassword, plainPassword)
	// Trim whitespace from password input to avoid issues
	plainPassword := strings.TrimSpace(req.Password)
	if !utils.CheckPasswordHash(user.Password, plainPassword) {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}
