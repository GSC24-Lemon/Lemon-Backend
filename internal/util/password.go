package util

import (
	"lemon_be/internal/controller/http/errorWrapper"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword return bcrypt hashed password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// return "", fmt.Errorf("failed to hash password: %w", err)
		return "", errorWrapper.NewHTTPError(err, 500, "failed to hash password")
	}

	return string(hashedPassword), nil
}

// CheckPassword check jika password yang diberikan cocok atau tidak dg hashed password dari database
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
