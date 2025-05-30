package entities

import "time"

type Order struct {
	ID           string       `json:"id"`
	UserID       string       `json:"userId"`
	TotalAmount  float64      `json:"totalAmount"`
	Status       string       `json:"status"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}
