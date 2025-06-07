package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/order_management/product_service/configs"
	"github.com/order_management/product_service/internal/api"
	"github.com/order_management/product_service/internal/db"
	"github.com/order_management/product_service/internal/kafka"
	"github.com/order_management/product_service/internal/repository"
	"github.com/order_management/product_service/internal/service"
	pkg "github.com/order_management/product_service/pkg/validator"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error to load configuration:%s", err)
	}
	db, err := db.InitDB(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Database initialization error:%+v", err)
	}
	kafkaProducer, err := kafka.NewProducer("product-topic", 1024)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	repo := repository.NewProductRepository(db)
	svc := service.NewServices(repo)
	api := api.NewAPiService(svc, kafkaProducer)

	e := echo.New()
	e.Validator = &pkg.CustomValidator{Validator: validator.New()}
	product := e.Group("/product")
	product.POST("/create", api.CreateProduct())
	//  Start server using config values
	serverAddress := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", serverAddress)
	if err := e.Start(serverAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
