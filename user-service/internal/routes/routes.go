package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/order_management/user_service/internal/handlers"
)

// SetupRoutes initializes all the routes for the server.
func SetupRoutes(e *echo.Echo, userHandler handlers.UserHandler) {
	userGroup := e.Group("/user")
	userGroup.POST("/create", userHandler.CreateUser())
}
