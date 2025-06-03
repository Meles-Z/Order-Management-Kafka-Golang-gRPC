package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/order_management/order_service/internal/database"
	"github.com/order_management/order_service/internal/kafka"
	"github.com/order_management/order_service/internal/repository"
	"github.com/order_management/order_service/internal/services"
)

const (
	defaultTopic   = "user-topic"
	defaultGroupID = "order-service-group"
	defaultPort    = "8080"
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

	err = kafka.EnsureTopicExists(kafkaBootstrap, topic, 1, 1)
	if err != nil {
		log.Fatalf("❌ Failed to ensure Kafka topic exists: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo, kafkaBootstrap, topic)

	// ✅ Start consuming users via Kafka
	if err := kafka.StartUserConsumer(kafkaBootstrap, groupID, topic, userSvc); err != nil {
		log.Fatalf("❌ Failed to start user Kafka consumer: %v", err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Order Service is running.")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})

	go func() {
		log.Printf("🌐 Starting HTTP server on port %s...", port)
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
