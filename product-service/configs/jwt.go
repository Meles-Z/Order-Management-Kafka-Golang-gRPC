package configs

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	jwt.RegisteredClaims
}

func GenerateToken(secretKey string, id string, name string, email string, phoneNumber string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		Id:          id,
		Name:        name,
		Email:       email,
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user-authentication",
			Subject:   id,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	return token.SignedString([]byte(secretKey))

}
