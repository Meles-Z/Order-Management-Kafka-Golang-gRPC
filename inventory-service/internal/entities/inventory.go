package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Inventory struct {
	Model
	ProductID string `json:"productId"` // Link to Product
	Quantity  int    `json:"quantity"`  // Available stock
	Reserved  int    `json:"reserved"`  // Reserved stock (in cart, waiting for payment)
	Location  string `json:"location"`
}

type Model struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	UpdatedAt time.Time      `json:"updatedAt"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.NewString()
	return nil
}
