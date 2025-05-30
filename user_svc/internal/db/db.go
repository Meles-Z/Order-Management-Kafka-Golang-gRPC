package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, name, password)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	conn.AutoMigrate()
	DB = conn
	fmt.Println("database connected successfully")
	return conn, nil
}
