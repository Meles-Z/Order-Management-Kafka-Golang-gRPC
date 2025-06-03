package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/order_management/user_svc/internal/entities"
	"github.com/order_management/user_svc/internal/kafka"
	"github.com/order_management/user_svc/internal/services"
)

type Handler struct {
	service *services.Serices
}

func NewHandler(s *services.Serices) *Handler {
	return &Handler{service: s}
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

		go func() {
			if err := kafka.KafkaProducer(user); err != nil {
				log.Printf("Failed to send user to Kafka: %v", err)
			}
		}()

		return c.JSON(http.StatusCreated, user)
	}
}
