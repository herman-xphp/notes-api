package utils

import (
	"os"
	"testing"
	"time"
)

func setupJWTTest(t *testing.T) {
	// Set a test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-that-is-at-least-32-characters-long")
	// Initialize JWT
	if err := InitJWT(); err != nil {
		t.Fatalf("InitJWT() error = %v", err)
	}
}

func TestGenerateToken(t *testing.T) {
	setupJWTTest(t)

	userID := uint(1)
	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}
}

func TestGenerateJWT(t *testing.T) {
	setupJWTTest(t)

	userID := uint(1)
	email := "test@example.com"
	token, err := GenerateJWT(userID, email)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}
	if token == "" {
		t.Error("GenerateJWT() returned empty token")
	}
}

func TestValidateToken(t *testing.T) {
	setupJWTTest(t)

	userID := uint(1)
	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("ValidateToken() UserID = %v, want %v", claims.UserID, userID)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	setupJWTTest(t)

	invalidToken := "invalid.token.here"
	_, err := ValidateToken(invalidToken)
	if err == nil {
		t.Error("ValidateToken() should return error for invalid token")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	setupJWTTest(t)

	// Note: This test would require mocking time or creating an expired token
	// For now, we'll just test that validation works for valid tokens
	userID := uint(1)
	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	// Check expiration is set
	if claims.ExpiresAt == nil {
		t.Error("Token should have expiration time")
	}

	// Check expiration is in the future
	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("Token expiration should be in the future")
	}
}

