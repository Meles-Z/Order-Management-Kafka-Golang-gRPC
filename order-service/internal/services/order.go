package services

import (
	"github.com/order_management/order_service/internal/entities"
	"github.com/order_management/order_service/internal/repository"
)

type Service struct {
	repo *repository.OrderRepo
}

func NewService(repo *repository.OrderRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateOrder(order *entities.Order) (*entities.Order, error) {
	// check user here?
	newOrder, err := s.repo.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	return newOrder, nil
}
