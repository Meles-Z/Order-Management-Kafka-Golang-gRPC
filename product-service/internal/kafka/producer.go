package kafka

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/order_management/product_service/internal/dto"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
	jobQueue chan *dto.Product
}

func NewProducer(topic string, bufferSize int) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	})
	if err != nil {
		return nil, err
	}

	producer := &Producer{
		producer: p,
		topic:    topic,
		jobQueue: make(chan *dto.Product, bufferSize),
	}

	// Start worker
	go producer.runWorker()

	return producer, nil
}

func (p *Producer) runWorker() {
	for msg := range p.jobQueue {
		p.publish(msg)
	}
}

func (p *Producer) publish(product *dto.Product) {
	data, err := json.Marshal(product)
	if err != nil {
		log.Printf("Error marshalling product: %v", err)
		return
	}
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)
	if err != nil {
		log.Printf("Kafka produce error: %v", err)
	}
}

func (p *Producer) Enqueue(product *dto.Product) {
	select {
	case p.jobQueue <- product:
	default:
		log.Println("Kafka job queue full; message dropped")
	}
}

func (p *Producer) Close() {
	close(p.jobQueue)
	p.producer.Close()
}
