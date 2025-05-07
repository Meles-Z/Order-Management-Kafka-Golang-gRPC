package services

import (
	"github.com/order_management/user_service/internal/configs"
	"github.com/order_management/user_service/internal/dto"
	"github.com/order_management/user_service/internal/entities"
	"github.com/order_management/user_service/internal/repository"
)

type UserService interface {
	CreateUser(*entities.User) (*entities.User, error)
	GetAllUsers() ([]entities.User, error)
	FindUserById(id string) (*entities.User, error)
	FindUserByEmail(email string) (*entities.User, error)
	UpdateUser(*entities.User) (*entities.User, error)
	UpdatePassword(*dto.PasswordUpdateDTO) error
	DeleteUser(id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
func (svc *userService) CreateUser(user *entities.User) (*entities.User, error) {
	user.Password = configs.HashAndSalt(user.Password)
	usr, err := svc.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (svc *userService) GetAllUsers() ([]entities.User, error) {
	users, err := svc.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (svc *userService) FindUserById(id string) (*entities.User, error) {
	user, err := svc.userRepo.FindUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *userService) FindUserByEmail(email string) (*entities.User, error) {
	user, err := svc.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *userService) UpdateUser(user *entities.User) (*entities.User, error) {
	newUser, err := svc.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (svc *userService) UpdatePassword(user *dto.PasswordUpdateDTO) error {
	err := svc.userRepo.UpdatePassword(user)
	if err != nil {
		return err
	}
	return nil
}

func (svc *userService) DeleteUser(id string) error {
	err := svc.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
