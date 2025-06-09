package services

import (
	"github.com/order_management/order_service/internal/entities"
	"github.com/order_management/order_service/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepo
}

func NewProductService(productRepo *repository.ProductRepo) *ProductService {
	return &ProductService{repo: productRepo}
}

func (s *ProductService) CreateProduct(product *entities.Product) (*entities.Product, error) {
	prod, err := s.repo.CreateProduct(product)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *ProductService) FindProductById(id string) (*entities.Product, error) {
	prod, err := s.repo.FindProductById(id)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *ProductService) UpdateProduct(product *entities.Product) (*entities.Product, error) {
	newProd, err := s.repo.UpdateProduct(product)
	if err != nil {
		return nil, err
	}
	return newProd, nil
}

func (s *ProductService) DeleteProduct(id string) error {
	err := s.repo.DeleteProduct(id)
	if err != nil {
		return err
	}
	return nil
}
