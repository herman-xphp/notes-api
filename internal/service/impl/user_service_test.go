package impl

import (
	"context"
	"notes-api/internal/domain"
	"notes-api/internal/dto"
	"testing"

	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	users      map[string]*domain.User
	usersByID  map[uint]*domain.User
	createErr  error
	findErr    error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:     make(map[string]*domain.User),
		usersByID: make(map[uint]*domain.User),
	}
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	if user, ok := m.usersByID[id]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.Email] = user
	m.usersByID[user.ID] = user
	return nil
}

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	service := NewUserServiceImpl(mockRepo)

	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// First registration should succeed
	resp, err := service.Register(ctx, req)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Register() returned nil response")
	}
	if resp.Email != req.Email {
		t.Errorf("Register() Email = %v, want %v", resp.Email, req.Email)
	}
	if resp.Token == "" {
		t.Error("Register() should return a token")
	}

	// Second registration with same email should fail
	_, err = service.Register(ctx, req)
	if err == nil {
		t.Error("Register() should return error for duplicate email")
	}
	if err.Error() != "email already exists" {
		t.Errorf("Register() error = %v, want 'email already exists'", err)
	}
}

func TestUserService_Login(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	service := NewUserServiceImpl(mockRepo)

	// First register a user
	registerReq := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := service.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Test successful login
	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	resp, err := service.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Login() returned nil response")
	}
	if resp.Email != loginReq.Email {
		t.Errorf("Login() Email = %v, want %v", resp.Email, loginReq.Email)
	}
	if resp.Token == "" {
		t.Error("Login() should return a token")
	}

	// Test login with wrong password
	wrongLoginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	_, err = service.Login(ctx, wrongLoginReq)
	if err == nil {
		t.Error("Login() should return error for wrong password")
	}

	// Test login with non-existent email
	nonExistentReq := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	_, err = service.Login(ctx, nonExistentReq)
	if err == nil {
		t.Error("Login() should return error for non-existent email")
	}
}

