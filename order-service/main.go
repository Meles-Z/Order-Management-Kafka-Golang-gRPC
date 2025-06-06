package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	broker := "localhost:9093" // Use localhost for external access
	topic := "user-topic"
	groupID := "user_service_group"

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	// Handle OS signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				log.Printf("Received message: %s", string(msg.Value))
			} else {
				log.Printf("Error reading message: %v", err)
			}
		}
	}()

	// Wait for shutdown signal
	<-sigchan
	log.Println("Shutting down consumer...")
}

/*


| Step                         | Protocol/Tech | Why It’s Smart                                    |
| ---------------------------- | ------------- | ------------------------------------------------- |
| 1. REST Gateway              | REST          | Easy for frontend, auth/security handled here     |
| 2. gRPC Order                | gRPC          | Fast, typed, ideal for internal Go services       |
| 3. gRPC to Product/Inventory | gRPC          | Needed sync response, keep services modular       |
| 4. Kafka OrderPlaced         | Kafka         | Decouples services, supports multiple consumers   |
| 5. Kafka for Inventory       | Kafka         | Async, retryable, no tight link to Order          |
| 6. Kafka for Payment         | Kafka         | Non-blocking, handles failure gracefully          |
| 7. Kafka for Notification    | Kafka         | Clean separation of concerns, flexible channels   |
| 8. Kafka for Analytics       | Kafka         | Real-time data stream, replayability, zero impact |



✅ Microservice Communication Steps (Execution Flow)
Step 1: User Auth + REST Request
User logs in via API Gateway (REST)
→ Authenticated via JWT or session

User places an order via POST /orders
→ Gateway receives the request and forwards to Order Service via gRPC

Step 2: Order Service (gRPC)
Order Service (via gRPC):

Validates User (if necessary, via gRPC to User Service)

Checks Product Availability (via gRPC to Product Service)

Checks Inventory Availability (via gRPC to Inventory Service)

If all OK → creates the order (in DB)

Step 3: Emit OrderPlaced Event
Order Service emits OrderPlaced event to Kafka
(e.g. to topic order.events)
→ This marks the point where order creation is done ✅

Step 4: Event-Driven Services React to Kafka
Inventory Service subscribes to OrderPlaced:

Reserves the stock

Optionally emits StockReserved event

Payment Service subscribes to OrderPlaced:

Processes payment (e.g. Stripe, PayPal)

Emits PaymentConfirmed or PaymentFailed

Notification Service subscribes to PaymentConfirmed:

Sends email/SMS/Push to user

Analytics Service subscribes to all events (OrderPlaced, PaymentConfirmed, etc.)

Stores data in dashboard/warehouse for reporting

🔁 Optional Additional Events
If PaymentFailed, Notification may still inform the user

Inventory could release the stock if payment fails (you may want a compensating OrderFailed event)

🧠 Summary of What You Build (In Order)
Step	What You Build
1.	RESTful API Gateway (authentication, validation, REST → gRPC call)
2.	gRPC-based Order Service
3.	Product and Inventory Service with gRPC endpoints
4.	Kafka setup with topics (e.g., order.events, payment.events, notifications)
5.	Event consumers: Inventory, Payment, Notification, Analytics

🎯 Final Flow Recap (Short Bullet Style)
✅ REST API Gateway handles login + create order

⚡ Gateway calls Order Service using gRPC

📦 Order Service checks product + inventory (gRPC), creates order

📣 Order Service emits OrderPlaced event (Kafka)

📦 Inventory Service listens → reserves stock

💳 Payment Service listens → processes payment → emits PaymentConfirmed

📬 Notification Service listens → sends confirmation email/SMS

📊 Analytics Service listens to all events → dashboards

*/
