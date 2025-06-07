package service

import (
	"github.com/order_management/product_service/internal/entities"
	"github.com/order_management/product_service/internal/repository"
)

type Services struct {
	repo *repository.ProdRepo
}

func NewServices(repo *repository.ProdRepo) *Services {
	return &Services{
		repo: repo,
	}
}

func (s *Services) CreateProduct(prod *entities.Product) (*entities.Product, error) {
	prod, err := s.repo.CreateProduct(prod)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *Services) FindProductById(id string) (*entities.Product, error) {
	prod, err := s.repo.FindProductById(id)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *Services) UpdateProduct(product *entities.Product) (*entities.Product, error) {
	newProduct, err := s.repo.UpdateProduct(product)
	if err != nil {
		return nil, err
	}
	return newProduct, nil
}

func (s *Services) DeleteProduct(id string) error {
	err := s.repo.DeleleProduct(id)
	if err != nil {
		return err
	}
	return nil
}
