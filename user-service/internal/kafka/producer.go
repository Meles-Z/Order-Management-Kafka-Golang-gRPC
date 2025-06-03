// // internal/kafka_message/producer.go
// package kafka

// import (
// 	"encoding/json"
// 	"fmt"

//     "github.com/confluentinc/confluent-kafka-go/v2/kafka"
// 	"github.com/order_management/user_service/internal/entities"
// )

// var kafkaBroker = "kafka:9092"
// var topic = "user.events.v1"

// func KafkaProducer(user any) error {
// 	p, err := kafka.NewProducer(&kafka.ConfigMap{
// 		"bootstrap.servers": kafkaBroker,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	defer p.Close()

// 	var userPayload UserCreatedPayload

// 	switch u := user.(type) {
// 	case *UserCreatedPayload:
// 		userPayload = *u
// 	case *entities.User:
// 		userPayload = UserCreatedPayload{
// 			ID:          u.ID,
// 			Name:        u.Name,
// 			Email:       u.Email,
// 			PhoneNumber: u.PhoneNumber,
// 			Address:     u.Address,
// 		}
// 	default:
// 		return fmt.Errorf("unsupported payload type: %T", user)
// 	}

// 	message := KafkaMessage{
// 		EventType: EventUserCreated,
// 		Payload:   userPayload,
// 	}

// 	bytes, err := json.Marshal(message)
// 	if err != nil {
// 		return err
// 	}

// 	return p.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{
// 			Topic:     &topic,
// 			Partition: kafka.PartitionAny,
// 		},
// 		Value: bytes,
// 	}, nil)
// }

package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/order_management/user_service/internal/entities"
)

const broker = "kafka:9092"

var topic = "user-topic"

func KafkaProducer(user *entities.User) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
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
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	msg, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user message: %w", err)
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
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
