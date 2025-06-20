package message

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/order_management/iventory_service/internal/dto"
	"github.com/order_management/iventory_service/internal/entities"
	"github.com/order_management/iventory_service/internal/service"
	"github.com/order_management/iventory_service/pkg/logger"
)

// Corrected: Match the incoming JSON payload field "id"
type InventoryCreatePayload struct {
	ProductID string `json:"id"`
}

type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewConsumer(bootstrapServers, groupId, topic string) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	}
	consumer, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %v", err)
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
				logger.Error("Failed to read message", "error", err)
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

func StartConsumeProduct(bootstrapServers, groupId, topic string, inventoryService *service.Service) error {
	consumer, err := NewConsumer(bootstrapServers, groupId, topic)
	if err != nil {
		return err
	}

	handle := func(msg *kafka.Message) {
		logger.Info("✅ Received message:", "value:", string(msg.Value))

		var event dto.InventoryEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error("⚠️ Failed to unmarshal message", "error", err)
			return
		}

		switch event.EventType {
		case "create":
			var createPayload InventoryCreatePayload
			if err := json.Unmarshal(event.Payload, &createPayload); err != nil {
				logger.Error("❌ Failed to unmarshal create payload", "error", err)
				return
			}

			logger.Info("Parsed Product ID", "productId", createPayload.ProductID)

			existingInventory, err := inventoryService.FindInventoryByProductID(createPayload.ProductID)
			if err != nil {
				// Inventory does not exist: create it with quantity 1
				newInventory := entities.Inventory{
					ProductID: createPayload.ProductID,
					Quantity:  1,
				}

				logger.Info("Creating new inventory", "inventory", newInventory)

				if _, err := inventoryService.CreateEventory(&newInventory); err != nil {
					logger.Error("❌ Failed to create inventory", "error", err)
				} else {
					logger.Info("✅ Successfully created new inventory")
				}
			} else if existingInventory.ID != "" {
				// Inventory exists: increment quantity
				existingInventory.Quantity += 1

				if _, err := inventoryService.UpdateInvitories(existingInventory); err != nil {
					logger.Error("❌ Failed to update inventory", "error", err)
				} else {
					logger.Info("✅ Successfully updated inventory quantity", "inventory", existingInventory)
				}
			}

		default:
			logger.Info("⚠️ Unknown event type", "info", event.EventType)
		}
	}

	return consumer.Start(handle)
}
