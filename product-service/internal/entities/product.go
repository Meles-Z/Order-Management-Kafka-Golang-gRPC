package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	IsActive    bool    `json:"isActive"`
}

type Model struct {
	ID        string         `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

func (m *Model) BeforeCreate(tr *gorm.DB) (err error) {
	m.ID = uuid.NewString()
	return
}
