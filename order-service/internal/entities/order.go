package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	Model
	UserID      string  `json:"userId"`
	ProductID   string  `json:"productId"`
	TotalAmount float64 `json:"totalAmount"`
	Status      string  `json:"status"`
}

type Model struct {
	ID         string         `json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletetdAt gorm.DeletedAt `json:"deletedAt"`
}

// Hook: GORM will call this before insert
func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.NewString()
	return
}
