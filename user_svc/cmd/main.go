package main

import (
	"fmt"
	"log"
	"os"

	// "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/order_management/user_svc/internal/db"
	"github.com/order_management/user_svc/internal/handler"
	"github.com/order_management/user_svc/internal/repository"
	"github.com/order_management/user_svc/internal/services"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	fmt.Println("Starting User Service...")
	file := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	fmt.Println(file)
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("DB error: ", err)
	}

	repo := repository.NewUserRepository(dbConn)
	srv := services.NewUserService(repo)
	h := handler.NewHandler(srv)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "User Service is running.")
	})
	user := e.Group("/user")
	user.POST("/create", h.CreateUser())

	if err := e.Start(":8080"); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
