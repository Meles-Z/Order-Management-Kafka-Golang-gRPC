package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserClaim struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	jwt.RegisteredClaims
}

func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error:": "Token is missing",
			})
		}
		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Invalid header format")
		}
		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, UserClaim{}, func(t *jwt.Token) (interface{}, error) {
			return []byte("thisissecretkey"), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Token validation failed",
			})
		}
		claims, ok := token.Claims.(*UserClaim)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid token claim",
			})
		}
		c.Set("claims", claims)
		return next(c)
	}
}
