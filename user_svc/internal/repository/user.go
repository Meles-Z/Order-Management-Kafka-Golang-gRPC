package repository

import (
	"github.com/order_management/user_svc/internal/entities"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	err := repo.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) GetUsers() ([]entities.User, error) {
	var users []entities.User
	err := repo.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *UserRepository) FindUserById(id string) (*entities.User, error) {
	var user entities.User
	err := repo.DB.Where("id=?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := repo.DB.Where("email=?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(user *entities.User) (*entities.User, error) {
	var existingUser entities.User
	err := repo.DB.Model(&user).Where("id=?", user.ID).Updates(entities.User{
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
		Address:     user.Address,
		IsActive:    user.IsActive,
	}).Scan(&existingUser).Error
	if err != nil {
		return nil, err
	}
	return &existingUser, err
}

func (repo *UserRepository) DeleteUser(id string) error {
	var user entities.User
	err := repo.DB.Where("id=?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}
