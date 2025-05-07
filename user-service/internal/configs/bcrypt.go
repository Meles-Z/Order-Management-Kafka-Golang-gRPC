package configs

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Sprintf("Error hashing password: %v", err)
	}
	return string(hash)
}

func CompareAndPassword(hashedPassword, password string) bool {
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return false
	}
	return true
}
