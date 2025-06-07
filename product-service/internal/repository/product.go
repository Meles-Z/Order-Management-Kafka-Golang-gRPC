package repository

import (
	"github.com/order_management/product_service/internal/entities"
	"gorm.io/gorm"
)

type ProdRepo struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProdRepo {
	return &ProdRepo{DB: db}
}

func (r *ProdRepo) CreateProduct(product *entities.Product) (*entities.Product, error) {
	err := r.DB.Create(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProdRepo) FindProductById(id string) (*entities.Product, error) {
	var prod entities.Product
	err := r.DB.Where("id=?", id).Take(&prod).Error
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

func (r *ProdRepo) UpdateProduct(product *entities.Product) (*entities.Product, error) {
	var existingProduct entities.Product
	err := r.DB.Model(&entities.Product{}).Where("id=?", product.ID).Updates(map[string]interface{}{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"stock":       product.Stock,
		"is_active":   product.IsActive,
	}).Scan(&existingProduct).Error
	if err != nil {
		return nil, err
	}
	return &existingProduct, nil
}

func (r *ProdRepo) DeleleProduct(id string) error {
	var prod entities.Product
	err := r.DB.Where("id=?", id).Delete(&prod).Error
	if err != nil {
		return err
	}
	return nil
}
