package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/order_management/product_service/internal/dto"
)

type Producer struct {
	producer    *kafka.Producer
	topic       string
	jobQueue    chan *dto.ProductEvent
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	workerCount int
}

type Option func(*Producer)

func WithWorkerCount(count int) Option {
	return func(p *Producer) {
		p.workerCount = count
	}
}

func NewProducer(bufferSize int, options ...Option) (*Producer, error) {
	kafkaBootstrap := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBootstrap == "" {
		kafkaBootstrap = "kafka:9092"
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "product-topic"
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   kafkaBootstrap,
		"acks":                "all",
		"go.delivery.reports": true,
		"message.timeout.ms":  5000,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	producer := &Producer{
		producer:    p,
		topic:       kafkaTopic,
		jobQueue:    make(chan *dto.ProductEvent, bufferSize),
		ctx:         ctx,
		cancel:      cancel,
		workerCount: 4, // default
	}

	for _, opt := range options {
		opt(producer)
	}

	for i := 0; i < producer.workerCount; i++ {
		producer.wg.Add(1)
		go producer.runWorker()
	}

	go producer.handleDeliveryReports()

	return producer, nil
}

func (p *Producer) runWorker() {
	defer p.wg.Done()
	for {
		select {
		case event := <-p.jobQueue:
			if event != nil {
				p.publish(event)
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Producer) handleDeliveryReports() {
	for e := range p.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("❌ Delivery failed: %v\n", ev.TopicPartition.Error)
			} else {
				log.Printf("✅ Delivered message to %v [%d] at offset %v\n",
					*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
			}
		}
	}
}

func (p *Producer) publish(event *dto.ProductEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ Failed to marshal product event: %v\n", err)
		return
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value:     data,
		Timestamp: time.Now(),
	}, nil)

	if err != nil {
		log.Printf("❌ Kafka produce error: %v\n", err)
	}
}

func (p *Producer) Enqueue(event *dto.ProductEvent) error {
	select {
	case p.jobQueue <- event:
		return nil
	default:
		log.Println("⚠️ Kafka job queue full; message dropped")
		return ErrQueueFull
	}
}

func (p *Producer) EnqueueWithTimeout(event *dto.ProductEvent, timeout time.Duration) error {
	select {
	case p.jobQueue <- event:
		return nil
	case <-time.After(timeout):
		log.Println("⚠️ Kafka enqueue timeout; message dropped")
		return ErrQueueFull
	}
}

func (p *Producer) Close() {
	p.cancel()
	p.wg.Wait()
	close(p.jobQueue)
	p.producer.Flush(5000)
	p.producer.Close()
}

var ErrQueueFull = errors.New("kafka producer queue is full")
