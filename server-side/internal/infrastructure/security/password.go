package security

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher handles password hashing and verification
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher creates a new password hasher
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		cost: bcrypt.DefaultCost, // Cost of 10
	}
}

// Hash hashes a plain text password
func (ph *PasswordHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Verify verifies a plain text password against a hashed password
func (ph *PasswordHasher) Verify(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// IsValidPassword checks if a password is valid (returns true if valid)
func (ph *PasswordHasher) IsValidPassword(hashedPassword, plainPassword string) bool {
	err := ph.Verify(hashedPassword, plainPassword)
	return err == nil
}
