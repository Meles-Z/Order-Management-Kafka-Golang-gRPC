package repository

import (
	"github.com/order_management/iventory_service/internal/entities"
	"github.com/order_management/iventory_service/pkg/logger"
	"gorm.io/gorm"
)

type Repostory struct {
	DB *gorm.DB
}

func NewRepostiory(db *gorm.DB) *Repostory {
	return &Repostory{
		DB: db,
	}
}

func (r *Repostory) CreateEventory(invetory *entities.Inventory) (*entities.Inventory, error) {
	err := r.DB.Create(invetory).Error
	if err != nil {
		logger.Error("error to create invetory", "error", err)
		return nil, err
	}
	return invetory, nil
}

func (r *Repostory) GetInventory() ([]entities.Inventory, error) {
	var invetory []entities.Inventory
	err := r.DB.Find(&invetory).Error
	if err != nil {
		logger.Error("error to create invetory", "error", err)
		return nil, err
	}
	return invetory, nil
}

func (r *Repostory) FindInvetoryById(id string) (*entities.Inventory, error) {
	var invetory entities.Inventory
	err := r.DB.Where("id=?", id).Take(&invetory).Error
	if err != nil {
		logger.Error("invetory not found", "error", err)
	}
	return &invetory, nil
}

func (r *Repostory) UpdateInvitories(inventory *entities.Inventory) (*entities.Inventory, error) {
	var existingInventory entities.Inventory

	err := r.DB.
		Model(&entities.Inventory{}).
		Where("id = ?", inventory.ID).
		Updates(inventory).
		Scan(&existingInventory).Error

	if err != nil {
		logger.Error("error to update inventory", "error", err)
		return nil, err
	}

	return &existingInventory, nil
}

func (r *Repostory) DeleteInventory(id string) error {
	var inventory entities.Inventory
	result := r.DB.Where("id=?", id).Delete(&inventory)
	if result.Error != nil {
		logger.Error("inventory deleted successfully", "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		logger.Warn("no inventory found to delete", "id", id)
		return gorm.ErrRecordNotFound
	}

	logger.Info("inventory deleted successfully", "id", id)
	return nil
}
