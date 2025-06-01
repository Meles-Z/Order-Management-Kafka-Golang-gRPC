package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/order_management/user_svc/internal/db"
	"github.com/order_management/user_svc/internal/handler"
	"github.com/order_management/user_svc/internal/repository"
	"github.com/order_management/user_svc/internal/services"
)

// CustomValidator wraps go-playground/validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements the echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	fmt.Println("Starting User Service...")

	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("DB error: ", err)
	}

	repo := repository.NewUserRepository(dbConn)
	srv := services.NewUserService(repo)
	h := handler.NewHandler(srv)
	producer, err := kafka.NewProducer(&kafka.ConfigMap{})
	if err != nil {
		log.Fatalf("error:%s", err)
	}
	defer producer.Close()
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.POST("/users", h.CreateUser())

	// ✅ Start on port 8080 to match container config
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
