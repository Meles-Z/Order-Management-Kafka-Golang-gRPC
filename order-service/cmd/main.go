package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"

	kafkaClient "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/order_management/order_service/internal/database"
	"github.com/order_management/order_service/internal/entities"
	"github.com/order_management/order_service/internal/kafka"
	"github.com/order_management/order_service/internal/repository"
	"github.com/order_management/order_service/internal/services"
)

func main() {
	fmt.Println("Starting Order Service...")
	fmt.Println("New test app is working")
	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	kafkaBootstrap := os.Getenv("KAFKA_BOOTSTRAP")
	if kafkaBootstrap == "" {
		kafkaBootstrap = "kafka:9092"
	}
	topic := "orders"
	groupID := "order_service_group"

	// Ensure Kafka topic exists
	err = kafka.EnsureTopicExists(kafkaBootstrap, topic, 1, 1)
	if err != nil {
		log.Fatalf("Failed to ensure Kafka topic exists: %v", err)
	}

	// Create repository and service
	repo := repository.NewOrderRepo(db)
	srv := services.NewService(repo, kafkaBootstrap, topic)

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(kafkaBootstrap, groupID, topic)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Start Kafka consumer loop
	go func() {
		log.Printf("Starting Kafka consumer...")
		err := consumer.Consume(func(message *kafkaClient.Message) {
			var order entities.Order
			if err := json.Unmarshal(message.Value, &order); err != nil {
				log.Printf("Failed to unmarshal Kafka message: %v", err)
				return
			}

			_, err := srv.CreateOrder(&order)
			if err != nil {
				log.Fatalf("error to create user:%s", err)
			}

		})
		if err != nil {
			log.Fatalf("Kafka consume error: %v", err)
		}
	}()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		log.Println("Health check endpoint hit.")
		return c.String(200, "Order Service is running.")
	})
	e.Logger.Fatal(e.Start(":" + port))
}
