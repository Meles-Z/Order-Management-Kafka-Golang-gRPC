package services

import (
	"context"
	"fmt"

	"github.com/order_management/user_service/internal/configs"
	"github.com/order_management/user_service/internal/dto"
	"github.com/order_management/user_service/internal/entities"
	"github.com/order_management/user_service/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	FindUserById(ctx context.Context, id string) (*entities.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	UpdatePassword(ctx context.Context, pass *dto.PasswordUpdateDTO) error
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (svc *userService) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	hashedPassword, err := configs.HashAndSalt(user.Password)
	if err != nil {
		return nil, fmt.Errorf("error to hashing password:%s", err)
	}
	user.Password = hashedPassword
	usr, err := svc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (svc *userService) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	users, err := svc.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (svc *userService) FindUserById(ctx context.Context, id string) (*entities.User, error) {
	user, err := svc.userRepo.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *userService) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := svc.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *userService) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	newUser, err := svc.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (svc *userService) UpdatePassword(ctx context.Context, pass *dto.PasswordUpdateDTO) error {
	err := svc.userRepo.UpdatePassword(ctx, pass)
	if err != nil {
		return err
	}
	return nil
}

func (svc *userService) DeleteUser(ctx context.Context, id string) error {
	err := svc.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
