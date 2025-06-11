package database

import (
	"fmt"

	"github.com/order_management/iventory_service/configs"
	"github.com/order_management/iventory_service/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg configs.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	conn.AutoMigrate(&entities.Inventory{})
	DB = conn
	fmt.Println("Database connected successfully!")
	return conn, nil
}