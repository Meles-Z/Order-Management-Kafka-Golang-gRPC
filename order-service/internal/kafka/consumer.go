package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/order_management/order_service/internal/entities"
	"github.com/order_management/order_service/internal/services"
)

type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewConsumer(bootstrapServers, groupID, topic string) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (c *Consumer) Start(handler func(msg *kafka.Message)) error {
	if err := c.consumer.Subscribe(c.topic, nil); err != nil {
		return fmt.Errorf("failed to subscribe to topic: %v", err)
	}

	go func() {
		for {
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("⚠️ Failed to read message: %v", err)
				continue
			}
			handler(msg)
		}
	}()

	return nil
}

func (c *Consumer) Stop() {
	c.consumer.Unsubscribe()
	c.consumer.Close()
}

// Specific function to consume user data and create user
func StartUserConsumer(bootstrapServers, groupID, topic string, userService *services.UserService) error {
	consumer, err := NewConsumer(bootstrapServers, groupID, topic)
	if err != nil {
		return err
	}

	handler := func(msg *kafka.Message) {
		log.Printf("✅ Received message: %s", string(msg.Value))

		var user entities.User
		if err := json.Unmarshal(msg.Value, &user); err != nil {
			log.Printf("⚠️ Failed to unmarshal message: %v", err)
			return
		}

		log.Printf("🔄 Processing user: %+v", user)
		if _, err := userService.CreateUser(&user); err != nil {
			log.Printf("❌ Failed to create user: %v", err)
		} else {
			log.Printf("✅ Successfully processed user: %+v", user)
		}
	}

	return consumer.Start(handler)
}
