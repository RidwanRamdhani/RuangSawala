package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes password using bcrypt with cost factor of 12
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hashedPassword), err
}

// Comparing string hashedPassword and string password, then return true on match, or false on not match
func CompareHashAndPasswordString(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
