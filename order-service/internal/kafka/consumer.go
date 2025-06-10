package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/order_management/order_service/internal/dto"
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

		var event dto.UserEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("⚠️ Failed to unmarshal message: %v", err)
			return
		}

		switch event.EventType {
		case "create":
			var user entities.User
			if err := json.Unmarshal(event.Payload, &user); err != nil {
				log.Printf("❌ Failed to unmarshal create payload: %v", err)
				return
			}
			log.Printf("🔄 Processing user: %+v", user)
			if _, err := userService.CreateUser(&user); err != nil {
				log.Printf("❌ Failed to create user: %v", err)
			} else {
				log.Printf("✅ Successfully processed user: %+v", user)
			}
		case "update":
			var user entities.User
			if err := json.Unmarshal(event.Payload, &user); err != nil {
				log.Printf("❌ Failed to unmarshal create payload: %v", err)
				return
			}

			log.Printf("🔄 Updating User: %+v", user)
			if _, err := userService.UpdateUser(&user); err != nil {
				log.Printf("❌ Failed to update user: %v", err)
			} else {
				log.Printf("✅ User updated successfully: %+v", user)
			}

		case "delete":
			var deletePayload struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(event.Payload, &deletePayload); err != nil {
				log.Printf("❌ Failed to unmarshal delete payload: %v", err)
				return
			}

			log.Printf("🗑️ Delete event received for User ID: %v", deletePayload.ID)

			if err := userService.DeleteUser(deletePayload.ID); err != nil {
				log.Printf("⚠️ User with ID %s not found or delete failed: %v", deletePayload.ID, err)
			} else {
				log.Printf("✅ User deleted successfully: %v", deletePayload.ID)
			}

		default:
			log.Printf("⚠️ Unknown event type: %s", event.EventType)

		}

	}

	return consumer.Start(handler)
}

func StartProductConsumer(bootstrapServers, groupID, topic string, productService *services.ProductService) error {
	consumer, err := NewConsumer(bootstrapServers, groupID, topic)
	if err != nil {
		return err
	}

	handler := func(msg *kafka.Message) {
		log.Printf("✅ Received message:%s", string(msg.Value))

		var event dto.ProductEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("❌ Failed to unmarshal ProductEvent: %v", err)
			return
		}

		switch event.EventType {
		case "create":
			var product entities.Product
			if err := json.Unmarshal(event.Payload, &product); err != nil {
				log.Printf("❌ Failed to unmarshal create payload: %v", err)
				return
			}

			log.Printf("🚀 Creating Product: %+v", product)
			if _, err := productService.CreateProduct(&product); err != nil {
				log.Printf("❌ Failed to create product: %v", err)
			} else {
				log.Printf("✅ Product created successfully: %+v", product)
			}

		case "update":
			var product entities.Product
			if err := json.Unmarshal(event.Payload, &product); err != nil {
				log.Printf("❌ Failed to unmarshal update payload: %v", err)
				return
			}

			log.Printf("🔄 Updating Product: %+v", product)
			if _, err := productService.UpdateProduct(&product); err != nil {
				log.Printf("❌ Failed to update product: %v", err)
			} else {
				log.Printf("✅ Product updated successfully: %+v", product)
			}

		case "delete":
			var deletePayload struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(event.Payload, &deletePayload); err != nil {
				log.Printf("❌ Failed to unmarshal delete payload: %v", err)
				return
			}

			log.Printf("🗑️ Delete event received for Product ID: %v", deletePayload.ID)

			if err := productService.DeleteProduct(deletePayload.ID); err != nil {
				log.Printf("⚠️ Product with ID %s not found or delete failed: %v", deletePayload.ID, err)
			} else {
				log.Printf("✅ Product deleted successfully: %v", deletePayload.ID)
			}

		default:
			log.Printf("⚠️ Unknown event type: %s", event.EventType)
		}
	}

	return consumer.Start(handler)
}
