package dto

import "github.com/order_management/user_svc/internal/entities"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  *entities.User `json:"user"`
	Token string         `json:"token"`
}

type UserEvent struct {
	EventType string `json:"event_type"`
	Payload   any    `json:"payload"`
}

type UserCreateEvent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	IsActive    bool   `json:"isActive"`
}

type UserUpdateEvent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	IsActive    bool   `json:"isActive"`
}

type DeleteUserEvent struct {
	ID string `json:"id"`
}
