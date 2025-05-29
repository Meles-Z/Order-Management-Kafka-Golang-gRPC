// internal/kafka_message/producer.go
package kafkamessage

import (
	"encoding/json"
	"fmt"

    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/order_management/user_service/internal/entities"
)

var kafkaBroker = "kafka:9092"
var topic = "user.events.v1"

func KafkaProducer(user any) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
	})
	if err != nil {
		return err
	}
	defer p.Close()

	var userPayload UserCreatedPayload

	switch u := user.(type) {
	case *UserCreatedPayload:
		userPayload = *u
	case *entities.User:
		userPayload = UserCreatedPayload{
			ID:          u.ID,
			Name:        u.Name,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			Address:     u.Address,
		}
	default:
		return fmt.Errorf("unsupported payload type: %T", user)
	}

	message := KafkaMessage{
		EventType: EventUserCreated,
		Payload:   userPayload,
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: bytes,
	}, nil)
}
