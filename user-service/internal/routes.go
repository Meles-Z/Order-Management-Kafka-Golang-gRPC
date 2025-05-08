package internal

import "github.com/labstack/echo/v4"

func Routes(e *echo.Echo, server Server) {
	user:=e.Group("/user")
	user.POST("/create", server.userHandler.CreateUser())
}
