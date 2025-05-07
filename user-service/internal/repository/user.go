package repository

import (
	"github.com/order_management/user_service/internal/dto"
	"github.com/order_management/user_service/internal/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(*entities.User) (*entities.User, error)
	GetAllUsers() ([]entities.User, error)
	FindUserById(id string) (*entities.User, error)
	FindUserByEmail(email string) (*entities.User, error)
	UpdateUser(*entities.User) (*entities.User, error)
	UpdatePassword(pass *dto.PasswordUpdateDTO) error
	DeleteUser(id string) error
}

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		conn: db,
	}
}

func (repo *userRepository) CreateUser(user *entities.User) (*entities.User, error) {
	err := repo.conn.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) GetAllUsers() ([]entities.User, error) {
	var users []entities.User
	err := repo.conn.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) FindUserById(id string) (*entities.User, error) {
	var user entities.User
	err := repo.conn.Where("id=?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := repo.conn.Where("email=?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) UpdateUser(user *entities.User) (*entities.User, error) {
	err := repo.conn.Model(&user).Where("id=?", user.ID).Updates(entities.User{
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
	}).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) UpdatePassword(pass *dto.PasswordUpdateDTO) error {
	err := repo.conn.Model(&entities.User{}).Where("id=?", pass.ID).Update("password", pass.Password).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) DeleteUser(id string) error {
	var user entities.User
	err := repo.conn.Where("id=?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}
