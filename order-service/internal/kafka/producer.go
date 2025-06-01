package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

// NewKafkaProducer initializes a new Kafka producer.
func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokersToString(brokers),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// Publish sends a message to Kafka.
func (p *KafkaProducer) Publish(ctx context.Context, payload interface{}) error {
	message, err := p.encodeMessage(payload)
	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event)
	err = p.producer.Produce(&message, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Optionally wait for delivery report (with timeout)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
		}
	}
	close(deliveryChan)

	return nil
}

// encodeMessage serializes the payload and constructs a Kafka message.
func (p *KafkaProducer) encodeMessage(payload interface{}) (kafka.Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return kafka.Message{}, err
	}

	return kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          data,
		Timestamp:      time.Now(),
	}, nil
}

// brokersToString converts []string to a comma-separated string for Kafka config.
func brokersToString(brokers []string) string {
	return kafka.ConfigMap{"bootstrap.servers": kafka.ConfigValue(brokers)}["bootstrap.servers"].(string)
}
