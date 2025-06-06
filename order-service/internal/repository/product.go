package repository

import (
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
	var exsistingUser entities.Product
	err := repo.DB.Model(&entities.Product{}).Where("id=?", product.ID).Updates(entities.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		IsActive:    product.IsActive,
	}).Scan(&exsistingUser).Error
	if err != nil {
		return nil, err
	}
	return &exsistingUser, nil
}

func (repo *ProductRepo) DeleteProduct(id string) error {
	var product entities.Product
	err := repo.DB.Take("id=?", id).Delete(&product).Error
	if err != nil {
		return err
	}
	return nil
}
