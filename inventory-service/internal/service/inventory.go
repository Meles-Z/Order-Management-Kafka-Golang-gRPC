package service

import (
	"github.com/order_management/iventory_service/internal/entities"
	"github.com/order_management/iventory_service/internal/repository"
	"github.com/order_management/iventory_service/pkg/logger"
)

type Service struct {
	repo *repository.Repostory
}

func NewServices(repo *repository.Repostory) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateEventory(inventory *entities.Inventory) (*entities.Inventory, error) {
	inv, err := s.repo.CreateEventory(inventory)
	if err != nil {
		logger.Error("error to create invetory", "error", err)
		return nil, err
	}
	return inv, nil
}

func (r *Service) GetInventory() ([]entities.Inventory, error) {
	inv, err := r.repo.GetInventory()
	if err != nil {
		logger.Error("error to find invetory", "error", err)
		return nil, err
	}
	return inv, nil
}

func (r *Service) FindInvetoryById(id string) (*entities.Inventory, error) {
	inv, err := r.repo.FindInvetoryById(id)
	if err != nil {
		logger.Error("error to find invetory", "error", err)
		return nil, err
	}
	return inv, nil
}

func (r *Service) UpdateInvitories(env *entities.Inventory) (*entities.Inventory, error) {
	inv, err := r.repo.UpdateInvitories(env)
	if err != nil {
		logger.Error("error to update invetory", "error", err)
		return nil, err
	}
	return inv, nil
}
func (s *Service) FindInventoryByProductID(productID string) (*entities.Inventory, error) {
	return s.repo.GetInventoryByProductID(productID) // Implement this in your repo
}


func (r *Service) DeleteInvitory(id string) error {
	err := r.repo.DeleteInventory(id)
	if err != nil {
		logger.Error("error to delete invetory", "error", err)
		return err
	}
	return nil
}
