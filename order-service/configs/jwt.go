package configs

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CustomClaim struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	jwt.RegisteredClaims
}

func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Token is missing",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid Authorization header format",
			})
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("thisissecretkey"), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Token validation failed",
			})
		}

		claims, ok := token.Claims.(*CustomClaim)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid token claims",
			})
		}

		// Set claims to context for use in next handler
		c.Set("claims", claims)

		return next(c)
	}
}
