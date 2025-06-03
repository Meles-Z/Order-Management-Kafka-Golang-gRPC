package services

import (
	"github.com/order_management/user_svc/configs"
	"github.com/order_management/user_svc/internal/entities"
	"github.com/order_management/user_svc/internal/repository"
)

type Serices struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *Serices {
	return &Serices{
		repo: repo,
	}
}

func (s *Serices) CreateUser(user *entities.User) (*entities.User, error) {
	pass, err := configs.HashPasswod(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = pass
	user, err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Serices) GetUsers() ([]entities.User, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Serices) FindUserById(id string) (*entities.User, error) {
	user, err := s.repo.FindUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *Serices) FindUserByEmail(email string) (*entities.User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *Serices) UpdateUser(user *entities.User) (*entities.User, error) {
	user, err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Serices) DeleteUser(id string) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
