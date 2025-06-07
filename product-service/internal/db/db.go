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
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode",
		cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	conn.AutoMigrate(&entities.Product{})
	DB = conn
	return conn, nil
}
