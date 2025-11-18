package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}
	if hash == password {
		t.Error("HashPassword() returned plain password instead of hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// Test correct password
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash() should return true for correct password")
	}

	// Test incorrect password
	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("CheckPasswordHash() should return false for incorrect password")
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "testpassword123"
	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Each hash should be different (due to salt)
	if hash1 == hash2 {
		t.Error("HashPassword() should generate different hashes for same password")
	}

	// But both should verify correctly
	if !CheckPasswordHash(password, hash1) {
		t.Error("First hash should verify correctly")
	}
	if !CheckPasswordHash(password, hash2) {
		t.Error("Second hash should verify correctly")
	}
}

