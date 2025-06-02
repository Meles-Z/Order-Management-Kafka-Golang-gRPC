package repository

import (
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

// CreateUser
func (repo *UserRepo) CreateUser(user *entities.User) (*entities.User, error) {
	err := repo.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser
func (repo *UserRepo) UpdateUser(user *entities.User) (*entities.User, error) {
	err := repo.DB.Model(&entities.User{}).Where("id=?", user.ID).Updates(
		entities.User{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
	).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser

func (repo *UserRepo) DeleteUser(id string) error {
	var user entities.User
	err := repo.DB.Where("id=?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}
