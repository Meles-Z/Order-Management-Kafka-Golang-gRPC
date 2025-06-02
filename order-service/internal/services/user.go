package services

import (
	"github.com/order_management/order_service/internal/entities"
	"github.com/order_management/order_service/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{userRepo: &userRepo}
}

func (svc *UserService) CreateUser(user *entities.User) (*entities.User, error) {
	user, err := svc.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *UserService) UpdateUser(user *entities.User) (*entities.User, error) {
	user, err := svc.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *UserService) DeleteUser(id string) error {
	err := svc.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
