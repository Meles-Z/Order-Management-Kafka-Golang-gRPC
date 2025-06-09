package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

func main() {
	// Load environment variables with fallback
	kafkaBootstrap := getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092")
	userTopic := getEnv("KAFKA_TOPIC", "user-topic")
	userGroupID := getEnv("KAFKA_GROUP_ID", "order-service-group")
	productTopic := getEnv("PRODUCT_TOPIC", "product-topic")
	productGroupID := getEnv("PRODUCT_GROUP_ID", "product-service-group")
	port := getEnv("PORT", "8081")

	// Initialize DB
	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("❌ Failed to initialize DB: %v", err)
	}

	// Setup repositories and services
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, kafkaBootstrap, userTopic)

	orderRepo := repository.NewOrderRepo(db)
	orderService := services.NewService(orderRepo)

	productRepo := repository.NewProductRepository(db)
	productService := services.NewProductService(productRepo)

	// Start Kafka consumers
	if err := kafka.StartUserConsumer(kafkaBootstrap, userGroupID, userTopic, userService); err != nil {
		log.Fatalf("❌ Failed to start user consumer: %v", err)
	}
	if err := kafka.StartProductConsumer(kafkaBootstrap, productGroupID, productTopic, productService); err != nil {
		log.Fatalf("❌ Failed to start product consumer: %v", err)
	}

	// Set up Echo HTTP server
	e := echo.New()
	e.Validator = &pkg.CustomValidator{Validator: validator.New()}

	orderGroup := e.Group("/order", configs.VerifyToken)
	orderGroup.POST("/create", handler.CreateOrder(*orderService, *userService))

	// Health routes
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "✅ Order service running")
	})
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Start server in goroutine
	go func() {
		log.Printf("🚀 HTTP server running on port %s", port)
		if err := e.Start(":" + port); err != nil {
			log.Fatalf("❌ Echo server stopped: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT or SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("🛑 Shutdown signal received, stopping server...")
	if err := e.Shutdown(nil); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}
	log.Println("✅ Server exited properly")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
