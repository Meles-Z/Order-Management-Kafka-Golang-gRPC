package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/order_management/product_service/internal/entities"
	"github.com/order_management/product_service/internal/service"
)

type Handler struct {
	service *service.Services
}

func NewAPiService(svc *service.Services) *Handler {
	return &Handler{service: svc}
}
func (h *Handler) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(entities.Product)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error:": err.Error(),
			})
		}

		if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{
				"error:": err.Error(),
			})
		}

		product, err := h.service.CreateProduct(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"Error to create product": err.Error(),
			})
		}
		return c.JSON(http.StatusCreated, product)
	}
}
