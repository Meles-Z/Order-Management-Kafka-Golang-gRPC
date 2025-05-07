package configs

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims defines the payload for the JWT token.
type CustomClaims struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a signed JWT token for a user.
func GenerateJWTToken(secretKey, userID, name, email, phone, role string, expiration time.Duration) (string, error) {
	claims := CustomClaims{
		ID:          userID,
		Name:        name,
		Email:       email,
		PhoneNumber: phone,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			Issuer:    "user-authentication",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
