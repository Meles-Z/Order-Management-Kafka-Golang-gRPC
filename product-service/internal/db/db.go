package db

import (
	"fmt"

	"github.com/order_management/product_service/configs"
	"github.com/order_management/product_service/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg configs.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password, cfg.SSLMode)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate your models
	err = conn.AutoMigrate(&entities.Product{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	DB = conn
	return conn, nil
}
