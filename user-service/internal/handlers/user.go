package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/order_management/user_service/internal/configs"
	"github.com/order_management/user_service/internal/dto"
	"github.com/order_management/user_service/internal/entities"
	kafkamessage "github.com/order_management/user_service/internal/kafka_message"
	"github.com/order_management/user_service/internal/services"
)

type UserHandler interface {
	CreateUser() echo.HandlerFunc
	// GetAllUser(userHandler) echo.HandlerFunc
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		const maxBodySize = 1_048_576 // 1MB

		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
		defer cancel()

		// Read and validate JSON
		var createUser dto.CreateUserDTO
		err := configs.ReadJSON(c, &createUser, maxBodySize)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error":  "Invalid request payload",
				"detail": err.Error(),
			})
		}

		// Validate DTO
		if err := configs.Validate(createUser); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"error":  "Validation failed",
				"detail": err.Error(),
			})
		}
		user, err := h.userService.CreateUser(ctx, &entities.User{
			Name:        createUser.Name,
			Email:       createUser.Email,
			PhoneNumber: createUser.PhoneNumber,
			Password:    createUser.Password,
			Address:     createUser.Address,
		})

		err = kafkamessage.KafkaProducer(user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error":  "Failed to produce message to Kafka",
				"detail": err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, user)
	}
}
