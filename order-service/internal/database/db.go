package database

import (
	"fmt"
	"os"

	"github.com/order_management/order_service/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() (*gorm.DB, error) {
	// Read from .env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, name, password,
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate schema
	err = conn.AutoMigrate(&entities.User{}, &entities.Order{})
	if err != nil {
		return nil, err
	}

	DB = conn
	fmt.Println("✅ Connected to Postgres successfully!")
	return conn, nil
}
