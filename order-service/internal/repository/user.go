package repository

import (
	"errors"
	"fmt"

	"github.com/order_management/order_service/internal/entities"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) CreateUser(user *entities.User) (*entities.User, error) {
	err := repo.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepo) FindUserById(id string) (*entities.User, error) {
	var user entities.User
	err := repo.DB.Where("id=?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) UpdateUser(user *entities.User) (*entities.User, error) {
	err := repo.DB.Model(&entities.User{}).Where("id=?", user.ID).Updates(
		entities.User{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			IsActive:    user.IsActive,
		},
	).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepo) DeleteUser(id string) error {
	var user entities.User
	// First, check if the user exists
	err := repo.DB.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %s does not exist", id)
		}
		return err
	}

	// Now delete the existing user
	err = repo.DB.Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}
