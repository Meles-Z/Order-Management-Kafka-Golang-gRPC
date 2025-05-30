package repository

import (
	"errors"

	"github.com/order_management/order_service/internal/entities"
	"gorm.io/gorm"
)

type OrderRepo struct {
	DB *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{DB: db}
}

// / CreateOrder creates a new order
func (repo *OrderRepo) CreateOrder(order *entities.Order) (*entities.Order, error) {
	err := repo.DB.Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

// GetAllUsers retrieves all orders
func (repo *OrderRepo) GetAllOrders() ([]entities.Order, error) {
	var orders []entities.Order
	err := repo.DB.Find(&orders).Error
	return orders, err
}

// GetUserByID retrieves a order by their ID
func (repo *OrderRepo) GetOrderByID(id uint) (*entities.Order, error) {
	var order entities.Order
	err := repo.DB.First(&order, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

// UpdateUser updates the existing order
func (repo *OrderRepo) UpdateOrder(order *entities.Order) error {
	return repo.DB.Save(order).Error
}

// DeleteUser deletes a order by ID
func (repo *OrderRepo) DeleteOrder(id uint) error {
	return repo.DB.Delete(&entities.Order{}, id).Error
}
