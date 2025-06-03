// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	kafkaClient "github.com/confluentinc/confluent-kafka-go/v2/kafka"
// 	"github.com/labstack/echo/v4"
// 	"github.com/order_management/order_service/internal/database"
// 	"github.com/order_management/order_service/internal/entities"
// 	"github.com/order_management/order_service/internal/kafka"
// 	"github.com/order_management/order_service/internal/repository"
// 	"github.com/order_management/order_service/internal/services"
// )
// const (
// 	broker = "localhost:9092"
// 	topic  = "test-topic"
// )

// func main() {
// 	fmt.Println("Starting Order Service...")
// 	db, err := database.InitDb()
// 	if err != nil {
// 		log.Fatalf("Database initialization failed: %v", err)
// 	}

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	kafkaBootstrap := os.Getenv("KAFKA_BOOTSTRAP")
// 	if kafkaBootstrap == "" {
// 		kafkaBootstrap = "kafka-order:9092" // Using container name in Docker network
// 	}
// 	topic := "user-topic"
// 	groupID := "user_service_group"

// 	// Enhanced Kafka configuration
// 	err = kafka.EnsureTopicExists(kafkaBootstrap, topic, 1, 1)
// 	if err != nil {
// 		log.Fatalf("Failed to ensure Kafka topic exists: %v", err)
// 	}

// 	userRepo := repository.NewUserRepository(db)
// 	userSvc := services.NewUserService(userRepo, kafkaBootstrap, topic)

// 	consumerConfig := &kafkaClient.ConfigMap{
// 		"bootstrap.servers":               kafkaBootstrap,
// 		"group.id":                        groupID,
// 		"auto.offset.reset":               "earliest",
// 		"enable.auto.commit":              false,
// 		"go.application.rebalance.enable": true,
// 		"session.timeout.ms":              60000, // Increased session timeout
// 		"heartbeat.interval.ms":           15000, // Increased heartbeat interval
// 	}

// 	consumer, err := kafkaClient.NewConsumer(consumerConfig)
// 	if err != nil {
// 		log.Fatalf("Failed to create Kafka consumer: %v", err)
// 	}
// 	defer consumer.Close()

// 	// Start consumer in a separate goroutine
// 	go func() {
// 		log.Printf("Subscribing to topic: %s", topic)
// 		if err := consumer.Subscribe(topic, nil); err != nil {
// 			log.Fatalf("Failed to subscribe to topic: %v", err)
// 		}

// 		for {
// 			msg, err := consumer.ReadMessage(100 * time.Millisecond)
// 			if err != nil {
// 				if err.(kafkaClient.Error).Code() == kafkaClient.ErrTimedOut {
// 					continue
// 				}
// 				log.Printf("Consumer error: %v", err)
// 				continue
// 			}

// 			log.Printf("Received message from partition %d at offset %d: %s\n",
// 				msg.TopicPartition.Partition, msg.TopicPartition.Offset, string(msg.Value))

// 			var user entities.User
// 			if err := json.Unmarshal(msg.Value, &user); err != nil {
// 				log.Printf("Failed to unmarshal message: %v (Raw: %s)", err, string(msg.Value))
// 				continue
// 			}

// 			log.Printf("Processing user: %+v", user)
// 			if _, err := userSvc.CreateUser(&user); err != nil {
// 				log.Printf("Failed to create user: %v", err)
// 			} else {
// 				log.Printf("Successfully processed user: %+v", user)
// 			}

// 			// Manually commit offset
// 			if _, err := consumer.CommitMessage(msg); err != nil {
// 				log.Printf("Failed to commit offset: %v", err)
// 			}
// 		}
// 	}()

// 	e := echo.New()
// 	e.GET("/", func(c echo.Context) error {
// 		return c.String(200, "Order Service is running.")
// 	})

// 	e.GET("/health", func(c echo.Context) error {
// 		// Simple health check
// 		return c.JSON(200, map[string]string{"status": "healthy"})
// 	})

// 	log.Printf("Starting Order Service on port %s", port)
// 	e.Logger.Fatal(e.Start(":" + port))
// }

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/labstack/echo/v4"
	"github.com/order_management/order_service/internal/entities"
)

const (
	defaultTopic   = "user-topic"
	defaultGroupID = "order-service-group"
	defaultPort    = "8080"
)

func main() {
	fmt.Println("🚀 Starting Order Service...")

	kafkaBootstrap := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBootstrap == "" {
		kafkaBootstrap = "kafka:9092"
		// Default to centralized Kafka
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

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBootstrap,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("❌ Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBootstrap,
	})
	if err != nil {
		log.Fatalf("❌ Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Subscribe to topic
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("❌ Failed to subscribe to topic: %v", err)
	}

	// Consume messages
	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("⚠️ Error reading message: %v", err)
				continue
			}

			log.Printf("✅ Received message from topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))

			var user entities.User // Replace with your actual user struct
			if err := json.Unmarshal(msg.Value, &user); err != nil {
				log.Printf("⚠️ Failed to unmarshal user: %v", err)
				continue
			}

			log.Printf("🔄 Processing user: %+v", user)
			// Add your business logic here
		}
	}()

	// Setup Echo server
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Order Service is running.")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})

	// Graceful shutdown
	go func() {
		log.Printf("🌐 Starting HTTP server on port %s...", port)
		if err := e.Start(":" + port); err != nil && err != echo.ErrServiceUnavailable {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
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
