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
	// product.Use(middleware.VerifyToken)
	product.POST("", h.CreateProduct())
	product.GET("", h.FindProductById())
	product.PUT("", h.UpdateProduct())

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
		return c.JSON(http.StatusOK, prod)
	}
}
