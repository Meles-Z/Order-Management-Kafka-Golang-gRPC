package db

import (
	"fmt"

	"github.com/order_management/user_service/internal/configs"
	"github.com/order_management/user_service/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUsername, cfg.DBName, cfg.DBPassword)

	orderdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = orderdb

	if err := orderdb.AutoMigrate(&entities.User{}); err != nil {
		return nil, err
	}

	fmt.Println("database is connected successfully")
	return orderdb, nil
}
