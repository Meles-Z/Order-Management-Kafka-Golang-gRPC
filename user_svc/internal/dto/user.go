package dto

import "github.com/order_management/user_svc/internal/entities"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  entities.User `json:"user"`
	Token string        `json:"token"`
}
