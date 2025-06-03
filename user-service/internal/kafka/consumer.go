package kafka

import (
	"context"
	"encoding/json"
	"log"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func StartKafkaConsumer(ctx context.Context, broker, topic, groupID string, handler func(KafkaMessage)) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":        broker,
		"group.id":                 groupID,
		"auto.offset.reset":        "earliest",
		"go.events.channel.enable": true,
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	log.Println("[Kafka] Consumer started on topic:", topic)

runLoop:
	for {
		select {
		case <-ctx.Done():
			log.Println("[Kafka] Consumer shutting down...")
			break runLoop
		default:
			event := consumer.Poll(100)
			switch e := event.(type) {
			case *kafka.Message:
				var msg KafkaMessage
				if err := json.Unmarshal(e.Value, &msg); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					continue
				}
				handler(msg)

			case kafka.Error:
				log.Printf("Kafka error: %v", e)
			}
		}
	}
}
