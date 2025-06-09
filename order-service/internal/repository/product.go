package repository

import (
	"errors"
	"fmt"

	"github.com/order_management/order_service/internal/entities"
	"gorm.io/gorm"
)

type ProductRepo struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		DB: db,
	}
}

func (repo *ProductRepo) CreateProduct(product *entities.Product) (*entities.Product, error) {
	err := repo.DB.Create(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *ProductRepo) FindProductById(id string) (*entities.Product, error) {
	var product entities.Product
	err := repo.DB.Where("id=?", id).Find(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *ProductRepo) UpdateProduct(product *entities.Product) (*entities.Product, error) {
	var existingProduct entities.Product

	// err := repo.DB.
	// Model(&entities.Product{}).
	// Where("id = ?", product.ID).
	// Select("*"). // 👈 force all fields to be updated
	// Updates(product).
	// Scan(&existingProduct).Error
	err := repo.DB.
		Model(&entities.Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]interface{}{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"is_active":   product.IsActive,
		}).
		Scan(&existingProduct).Error

	if err != nil {
		return nil, err
	}
	return &existingProduct, nil
}

func (repo *ProductRepo) DeleteProduct(id string) error {
	var product entities.Product
	// First, check if the product exists
	err := repo.DB.First(&product, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("product with ID %s does not exist", id)
		}
		return err
	}

	// Now delete the existing product
	err = repo.DB.Delete(&product).Error
	if err != nil {
		return err
	}
	return nil
}
