package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// InitJWT initializes JWT secret from environment variable
// Must be called after godotenv.Load() in main.go
// JWT_SECRET must be at least 32 characters long for security
func InitJWT() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return os.ErrInvalid
	}
	// Minimum 32 characters for HS256 algorithm security
	if len(secret) < 32 {
		return os.ErrInvalid
	}
	jwtSecret = []byte(secret)
	return nil
}

// InitJWTWithSecret initializes JWT secret with explicit parameter
// Used when loading from config instead of environment variable
func InitJWTWithSecret(secret string) error {
	if secret == "" {
		return os.ErrInvalid
	}
	if len(secret) < 32 {
		return os.ErrInvalid
	}
	jwtSecret = []byte(secret)
	return nil
}

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT baru untuk user
func GenerateToken(userID uint) (string, error) {
	return GenerateJWT(userID, "")
}

// GenerateJWT membuat JWT dengan userID dan email
func GenerateJWT(userID uint, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken memvalidasi JWT
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
