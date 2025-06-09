package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/order_management/order_service/configs"
	"github.com/order_management/order_service/internal/entities"
	"github.com/order_management/order_service/internal/services"
	"github.com/order_management/order_service/pkg"
)

func CreateOrder(orderSvc services.Service, userSvc services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("claims").(*configs.CustomClaim)
		if claims == nil {
			err := pkg.New("Claim is empty", http.StatusUnauthorized)
			return c.JSON(err.StatusCode, err)
		}

		req := new(entities.Order)
		if err := c.Bind(req); err != nil {
			err := pkg.New("Error to bind json data", http.StatusBadRequest)
			return c.JSON(err.StatusCode, err.Error())
		}

		if err := c.Validate(req); err != nil {
			err := pkg.New("Error to validate data", http.StatusUnprocessableEntity)
			return c.JSON(err.StatusCode, err.Error())
		}
		// check user is exist
		_, err := userSvc.FindUserById(claims.ID)
		if err != nil {
			err := pkg.New("Error to get user", http.StatusNotFound)
			return c.JSON(err.StatusCode, err)
		}
		req.UserID = claims.ID

		order, err := orderSvc.CreateOrder(req)
		if err != nil {
			err := pkg.New("error to create order:", http.StatusInternalServerError)
			return c.JSON(err.StatusCode, err.Error())
		}
		return c.JSON(http.StatusOK, order)
	}
}
