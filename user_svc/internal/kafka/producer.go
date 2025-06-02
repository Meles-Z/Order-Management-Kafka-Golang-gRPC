package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka" // This import is correct
	"github.com/order_management/user_svc/internal/entities"
)

const broker = "kafka:9093"

var kafkaTopic = "user-topic"

func KafkaProducer(message *entities.User) error {
	// kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{ // Correct: kafka.NewProducer
		"bootstrap.servers":   broker,
		"acks":                "all",
		"go.delivery.reports": true,
	})
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	defer producer.Close()

	// Start a goroutine to listen for delivery reports
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message: // Correct: kafka.Message
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					// Use a pointer to TopicPartition for clarity, as per Confluent's examples
					// And ensure the topic name is dereferenced for logging
					log.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal user message: %w", err)
	}

	err = producer.Produce(&kafka.Message{ // Correct: kafka.Message
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny}, // Correct: kafka.TopicPartition, kafka.PartitionAny
		Value:          msg,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce Kafka message: %w", err)
	}
	// Wait for delivery report
	producer.Flush(5000)
	log.Printf("Produced Kafka message: %s\n", string(msg))

	return nil
}