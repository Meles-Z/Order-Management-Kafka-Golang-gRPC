package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(broker []string, groupID, topic string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker[0],
		"group.id":              groupID,
		"auto.offset.reset":     "earliest",
		"enable.auto.commit":    false, // Manually commit after processing
		"session.timeout.ms":    60000, // Adjust as necessary
		"heartbeat.interval.ms": 20000, // Should be less than session timeout
	})
	if err != nil {
		return nil, err
	}

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		consumer.Close()
		return nil, err
	}

	return consumer, nil
}

// ConsumeMessages continuously reads messages from the Kafka topic
func ConsumeMessages(consumer *kafka.Consumer, processMessage func(msg *kafka.Message)) {
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
				continue
			}
			log.Printf("❌ Error reading message: %s", err)
			continue
		}
		processMessage(msg)
	}
}
