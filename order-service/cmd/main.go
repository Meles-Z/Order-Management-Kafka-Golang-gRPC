package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/order_management/order_service/configs"
	"github.com/order_management/order_service/internal/database"
	"github.com/order_management/order_service/internal/handler"
	"github.com/order_management/order_service/internal/kafka"
	"github.com/order_management/order_service/internal/repository"
	"github.com/order_management/order_service/internal/services"
	"github.com/order_management/order_service/pkg"
)

const (
	defaultTopic         = "user-topic"
	defaultGroupID       = "order-service-group"
	defaultPort          = "8081"
	productTopic         = "product-topic"
	productReaderGroupId = "product-service-group"
)

func main() {
	fmt.Println("\U0001F680 Starting Order Service...")

	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("❌ Database initialization failed: %v", err)
	}

	kafkaBootstrap := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBootstrap == "" {
		kafkaBootstrap = "kafka:9092"
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = defaultTopic
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = defaultGroupID
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	displayPort := os.Getenv("DISPLAY_PORT")
	if displayPort == "" {
		displayPort = port
	}

	err = kafka.EnsureTopicExists(kafkaBootstrap, topic, 1, 1)
	if err != nil {
		log.Fatalf("❌ Failed to ensure Kafka topic exists: %v", err)
	}

	err = kafka.EnsureTopicExists(kafkaBootstrap, productTopic, 1, 1)
	if err != nil {
		log.Fatalf("❌ Failed to ensure Kafka product topic exists: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo, kafkaBootstrap, topic)
	orderRepo := repository.NewOrderRepo(db)
	orderSvc := services.NewService(orderRepo)
	productRepo := repository.NewProductRepository(db)
	productSvc := services.NewProductService(productRepo)

	// ✅ Start consuming users via Kafka
	if err := kafka.StartUserConsumer(kafkaBootstrap, groupID, topic, userSvc); err != nil {
		log.Fatalf("❌ Failed to start user Kafka consumer: %v", err)
	}

	// start consumer vial kafka for product
	if err := kafka.StartProductConsumer(kafkaBootstrap, productReaderGroupId, productTopic, productSvc); err != nil {
		log.Fatalf("Failed to start product consumer:%v", err)
	}

	e := echo.New()
	e.Validator = &pkg.CustomValidator{Validator: validator.New()}
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Order Service is running.")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})
	order := e.Group("/order")
	order.Use(configs.VerifyToken)
	order.POST("/create", handler.CreateUser(*orderSvc, *userSvc))

	go func() {
		log.Printf("🌐 Starting HTTP server listen on port %s...", displayPort)
		if err := e.Start(":" + port); err != nil && err != echo.ErrConflict {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("📴 Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Echo server shutdown error: %v", err)
	}
	log.Println("✅ Order Service stopped gracefully")
}
