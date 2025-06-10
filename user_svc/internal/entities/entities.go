package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Model
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	IsActive    bool   `json:"isActive"`
}

type Model struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate is a GORM hook that sets a UUID before inserting a record.
func (model *Model) BeforeCreate(tx *gorm.DB) (err error) {
	model.ID = uuid.NewString()
	return
}
