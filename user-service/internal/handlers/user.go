package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/order_management/user_service/internal/dto"
	"github.com/order_management/user_service/internal/services"
)

type UserHandler interface {
	CreateUser(userHandler) echo.HandlerFunc
	// GetAllUser(userHandler) echo.HandlerFunc
	// FindUserById(userHandler) echo.HandlerFunc
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (c *userHandler) CreateUser(userHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		createUser := new(dto.CreateUserDTO)
		if 
	}
}
