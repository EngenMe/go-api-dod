package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher provides methods for hashing and checking passwords
type PasswordHasher struct {
	Cost int
}

// NewPasswordHasher creates a new PasswordHasher
func NewPasswordHasher(cost int) *PasswordHasher {
	return &PasswordHasher{
		Cost: cost,
	}
}

// Hash hashes a password
func (h *PasswordHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Check checks if a password matches a hash
func (h *PasswordHasher) Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
