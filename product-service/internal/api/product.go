package api

import (
	"fmt"
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
	// product.Use(middleware.VerifyToken)
	product.POST("", h.CreateProduct())
	product.GET("", h.FindProductById())
	product.PUT("", h.UpdateProduct())
	product.DELETE("", h.DeleteProduct())

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
		createPayload := &dto.CreateProductEvent{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			IsActive:    product.IsActive,
		}

		event := &dto.ProductEvent{
			EventType: "create",
			Payload:   createPayload,
		}

		if err := h.producer.Enqueue(event); err != nil {
			// Log the error but don't fail the request
			c.Logger().Errorf("failed to enqueue product to Kafka: %v", err)
		}

		return c.JSON(http.StatusCreated, product)
	}
}

func (h *Handler) FindProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		product, err := h.service.FindProductById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error:": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, product)

	}
}

func (h *Handler) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(entities.Product)
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{
				"echo": err.Error(),
			})
		}
		if err := c.Validate(req); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"error:": err.Error(),
			})
		}

		prod, err := h.service.UpdateProduct(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error to update product": err.Error(),
			})
		}
		updateProductEvent := &dto.UpdateProductEvent{
			ID:          prod.ID,
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
			Stock:       prod.Stock,
			IsActive:    prod.IsActive,
		}
		event := &dto.ProductEvent{
			EventType: "update",
			Payload:   updateProductEvent,
		}
		if err := h.producer.Enqueue(event); err != nil {
			c.Logger().Errorf("failed to enqueue product to Kafka: %v", err)
		}
		return c.JSON(http.StatusOK, prod)
	}
}

func (h *Handler) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")

		// Always enqueue delete event to Kafka
		productDeleteEvent := &dto.DeleteProductEvent{
			ID: id,
		}
		event := &dto.ProductEvent{
			EventType: "delete",
			Payload:   productDeleteEvent,
		}
		if err := h.producer.Enqueue(event); err != nil {
			c.Logger().Errorf("❌ Failed to enqueue product to Kafka: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to enqueue delete event to Kafka",
			})
		}

		// Check if product exists locally
		_, err := h.service.FindProductById(id)
		if err != nil {
			c.Logger().Warnf("⚠️ Product with ID %s does not exist locally: %v", id, err)
			return c.JSON(http.StatusOK, echo.Map{
				"message": "Product does not exist locally, but delete event sent to Kafka",
			})
		}

		// If found, delete it
		err = h.service.DeleteProduct(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": fmt.Sprintf("Error deleting local product: %s", err),
			})
		}

		// Success
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Product deleted locally and delete event sent to Kafka",
		})
	}
}
