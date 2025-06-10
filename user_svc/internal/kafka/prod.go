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
	"github.com/order_management/user_svc/internal/dto"
)

type Producer struct {
	producer    *kafka.Producer
	topic       string
	jobQueue    chan *dto.UserEvent
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	workerCount int
}

type Option func(*Producer)

func NewProducer(bufferSize int, options ...Option) (*Producer, error) {
	kafkaBootstarp := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBootstarp == "" {
		kafkaBootstarp = "kafka:9092"
	}
	kafkaTopic := os.Getenv("user-topic")
	if kafkaTopic == "" {
		kafkaTopic = "user-topic"
	}
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   kafkaBootstarp,
		"acks":                "all",
		"go.delivery.reports": true,
		"message.timeout.ms":  5000,
	})
	if err != nil {
		return nil, err
	}
	cxt, cancel := context.WithCancel(context.Background())
	producer := &Producer{
		producer:    p,
		topic:       kafkaTopic,
		jobQueue:    make(chan *dto.UserEvent, bufferSize),
		ctx:         cxt,
		cancel:      cancel,
		workerCount: 4, // this is the default

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

func (p *Producer) publish(event *dto.UserEvent) {
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

func WithWorkerCount(count int) Option {
	return func(p *Producer) {
		p.workerCount = count
	}
}

func (p *Producer) Enqueue(event *dto.UserEvent) error {
	select {
	case p.jobQueue <- event:
		return nil
	default:
		log.Println("⚠️ Kafka job queue full; message dropped")
		return ErrQueueFull
	}

}

func (p *Producer) EnqueueWithTimeout(event *dto.UserEvent, timeout time.Duration) error {
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
