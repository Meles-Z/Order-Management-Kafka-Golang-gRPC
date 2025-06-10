package kafka

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
// 	"github.com/order_management/user_svc/internal/entities"
// )

// func KafkaProducer(message *entities.User) error {
// 	// Get Kafka configuration from environment
// 	kafkaBootstrap := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
// 	if kafkaBootstrap == "" {
// 		kafkaBootstrap = "kafka:9092"
// 	}
// 	fmt.Println("kafka:", kafkaBootstrap)

// 	kafkaTopic := os.Getenv("KAFKA_TOPIC")
// 	if kafkaTopic == "" {
// 		kafkaTopic = "user-topic"
// 	}
// 	fmt.Println("Kafka topic:", kafkaTopic)

// 	producer, err := kafka.NewProducer(&kafka.ConfigMap{
// 		"bootstrap.servers":   kafkaBootstrap,
// 		"acks":                "all",
// 		"go.delivery.reports": true,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create Kafka producer: %w", err)
// 	}
// 	defer producer.Close()

// 	msg, err := json.Marshal(message)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal user message: %w", err)
// 	}

// 	deliveryChan := make(chan kafka.Event)

// 	err = producer.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{
// 			Topic:     &kafkaTopic,
// 			Partition: kafka.PartitionAny,
// 		},
// 		Value: msg,
// 	}, deliveryChan)

// 	if err != nil {
// 		close(deliveryChan)
// 		return fmt.Errorf("failed to produce Kafka message: %w", err)
// 	}

// 	e := <-deliveryChan
// 	m := e.(*kafka.Message)

// 	if m.TopicPartition.Error != nil {
// 		log.Printf("❌ Delivery failed: %v\n", m.TopicPartition.Error)
// 		return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
// 	}

// 	log.Printf("✅ Delivered message to topic %s at partition %d, offset %v\n",
// 		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

// 	close(deliveryChan)
// 	return nil
// }
