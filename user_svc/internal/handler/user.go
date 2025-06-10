package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/order_management/user_svc/internal/dto"
	"github.com/order_management/user_svc/internal/entities"
	"github.com/order_management/user_svc/internal/kafka"
	"github.com/order_management/user_svc/internal/services"
)

type Handler struct {
	service  *services.Serices
	producer *kafka.Producer
}

func NewHandler(s *services.Serices, producer *kafka.Producer) *Handler {
	return &Handler{
		service:  s,
		producer: producer,
	}
}

func RegisterRoutes(e *echo.Echo, h *Handler) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "User Service is running.")
	})

	user := e.Group("users")
	user.POST("", h.CreateUser())
	user.POST("/login", h.Login())
	user.PUT("", h.UpdateUser())
	user.DELETE("", h.DeleteUser())
}

func (h *Handler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(entities.User)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
		if err := c.Validate(req); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": err.Error()})
		}

		user, err := h.service.CreateUser(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		userCreateEvent := &dto.UserCreateEvent{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
			IsActive:    user.IsActive,
		}

		event := &dto.UserEvent{
			EventType: "create",
			Payload:   userCreateEvent,
		}
		if err := h.producer.Enqueue(event); err != nil {
			c.Logger().Errorf("failed to enqueue user to Kafka: %v", err)
		}

		return c.JSON(http.StatusCreated, user)
	}
}

func (h *Handler) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(entities.User)
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}

		if err := c.Validate(req); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"error:": err.Error(),
			})
		}
		user, err := h.service.UpdateUser(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error to update user": err.Error(),
			})
		}
		updateUserEvent := &dto.UserUpdateEvent{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
			IsActive:    user.IsActive,
		}

		event := &dto.UserEvent{
			EventType: "update",
			Payload:   updateUserEvent,
		}
		if err := h.producer.Enqueue(event); err != nil {
			c.Logger().Errorf("failed to enqueue user to Kafka: %v", err)
		}
		return c.JSON(http.StatusOK, user)
	}
}

func (h *Handler) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		// enqueue delete event id for deleting or for checking isdeleted is order local copy
		userDeleteEvent := &dto.DeleteUserEvent{
			ID: id,
		}
		event := &dto.UserEvent{
			EventType: "delete",
			Payload:   userDeleteEvent,
		}

		if err := h.producer.Enqueue(event); err != nil {
			c.Logger().Errorf("❌ Failed to enqueue user to Kafka: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to enqueue delete event to Kafka",
			})
		}

		// Check if user exists locally
		_, err := h.service.FindUserById(id)
		if err != nil {
			c.Logger().Warnf("⚠️ User with ID %s does not exist locally: %v", id, err)
			return c.JSON(http.StatusOK, echo.Map{
				"message": "User does not exist locally, but delete event sent to Kafka",
			})
		}
		err = h.service.DeleteUser(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": fmt.Sprintf("Error deleting local user: %s", err),
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "User deleted locally and delete event sent to Kafka",
		})
	}
}
