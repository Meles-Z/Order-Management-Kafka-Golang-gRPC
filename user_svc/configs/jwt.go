package configs

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomeClaim struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	jwt.RegisteredClaims
}

func GenerateToken(secretKey, id, name, email, phoneNumber string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomeClaim{
		ID:          id,
		Name:        name,
		Email:       email,
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user-authentication",
			Subject:   id,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	return token.SignedString([]byte(secretKey))
}
