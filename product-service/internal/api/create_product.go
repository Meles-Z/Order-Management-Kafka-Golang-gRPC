package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/order_management/product_service/internal/dto"
	"github.com/order_management/product_service/internal/entities"
	"github.com/order_management/product_service/internal/kafka"
	"github.com/order_management/product_service/internal/service"
)

type Handler struct {
	service  *service.Services
	producer *kafka.Producer
}

func NewAPiService(svc *service.Services, producer *kafka.Producer) *Handler {
	return &Handler{
		service:  svc,
		producer: producer,
	}
}

func RegisterRoutes(e *echo.Echo, h *Handler) {
	product := e.Group("/products")
	product.POST("", h.CreateProduct())
}

func (h *Handler) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req entities.Product
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		product, err := h.service.CreateProduct(&req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create product")
		}

		// Convert to DTO and send to Kafka
		productDTO := dto.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			IsActive:    product.IsActive,
		}

		if err := h.producer.Enqueue(&productDTO); err != nil {
			// Log the error but don't fail the request
			c.Logger().Errorf("failed to enqueue product to Kafka: %v", err)
		}

		return c.JSON(http.StatusCreated, product)
	}
}
