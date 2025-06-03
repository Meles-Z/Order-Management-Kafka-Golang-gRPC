package handler

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/order_management/user_svc/configs"
	"github.com/order_management/user_svc/internal/dto"
)

func (h *Handler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(dto.LoginRequest)

		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.ErrBadRequest)
		}
		if err := c.Validate(req); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.ErrUnprocessableEntity)
		}

		user, err := h.service.FindUserByEmail(req.Email)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.ErrNotFound)
		}

		if !configs.VerifyPassord(user.Password, req.Password) {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Invalid password",
			})
		}

		key := os.Getenv("SECRET_KEY")

		token, err := configs.GenerateToken(key, user.ID, user.Name, user.Email, user.PhoneNumber)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.ErrInternalServerError.Internal)
		}

		return c.JSON(http.StatusOK, &dto.LoginResponse{
			User:  user,
			Token: token,
		})
	}
}
