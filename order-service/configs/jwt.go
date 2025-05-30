package configs

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomeClaim struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(secretKey, id, name, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomeClaim{
		ID:   id,
		Name: name,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user-authentication",
			Subject:   id,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	return token.SignedString([]byte(secretKey))
}
