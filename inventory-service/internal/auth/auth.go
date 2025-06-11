package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/order_management/iventory_service/pkg/logger"
)

type CustomUserClaim struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	jwt.RegisteredClaims
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			logger.Error("Authorization header missing")
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Error("Invalid authorization header format")
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid authorization format"})
		}

		tokenString := parts[1]
		secret := os.Getenv("SECRET")

		token, err := jwt.ParseWithClaims(tokenString, &CustomUserClaim{}, func(token *jwt.Token) (any, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil {
			logger.Error("Invalid token", "error", err)
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(*CustomUserClaim)
		if !ok || !token.Valid {
			logger.Error("Invalid JWT claims")
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token claims"})
		}

		// Set claims to context so you can access it in handlers
		c.Set("claims", claims)
		return next(c)
	}
}