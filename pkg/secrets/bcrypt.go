package secrets

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the provided password with bcrypt.DefaultCost.
// Returns the hash as a string, or an error.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPasswordHash compares a password with a bcrypt hash.
// Returns true if the password matches the hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
