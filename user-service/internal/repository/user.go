package repository

import (
	"context"

	"github.com/order_management/user_service/internal/dto"
	"github.com/order_management/user_service/internal/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	FindUserById(ctx context.Context, id string) (*entities.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	UpdatePassword(ctx context.Context, pass *dto.PasswordUpdateDTO) error
	DeleteUser(ctx context.Context, id string) error
}

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		conn: db,
	}
}

func (repo *userRepository) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	err := repo.conn.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	err := repo.conn.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) FindUserById(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	err := repo.conn.WithContext(ctx).Where("id=?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := repo.conn.WithContext(ctx).Where("email=?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	err := repo.conn.WithContext(ctx).Model(&user).Where("id=?", user.ID).Updates(entities.User{
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

func (repo *userRepository) UpdatePassword(ctx context.Context, pass *dto.PasswordUpdateDTO) error {
	err := repo.conn.WithContext(ctx).Model(&entities.User{}).Where("id=?", pass.ID).Update("password", pass.Password).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) DeleteUser(ctx context.Context, id string) error {
	var user entities.User
	err := repo.conn.WithContext(ctx).Where("id=?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}