package configs

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPasswod(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password not hashed correcly:%s", err)
	}
	return string(hashed), nil
}

func VerifyPassord(hashedPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
