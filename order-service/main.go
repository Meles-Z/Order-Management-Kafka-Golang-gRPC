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


Great — you’re already doing a solid job using Kafka for async event-driven communication and maintaining local copies of user and product data in your Order service for performance. Now, about gRPC:

✅ When to Use gRPC in Your Microservices Architecture
Since you're already using Kafka for event-driven, eventually consistent, pub/sub communication, gRPC should be used where synchronous, fast, and direct request-response calls are needed between services.

Here’s where and why you’d use gRPC in your system:

1. Order Service → Inventory Service (Check product stock before placing order)
✅ Use gRPC

Reason: Needs immediate response to confirm availability.

Alternative to polling Kafka topic or relying on local copies.

go
Copy
Edit
// Example gRPC method
rpc CheckProductStock(CheckStockRequest) returns (CheckStockResponse);
2. Order Service → Payment Service (Process payment after order initiated)
✅ Use gRPC

Reason: Needs confirmation if the payment is successful/failed immediately.

go
Copy
Edit
rpc ProcessPayment(PaymentRequest) returns (PaymentResponse);
3. Order Service → Shipping Service (Initiate shipment after payment)
✅ Use gRPC

Reason: Immediate acknowledgment is needed.

go
Copy
Edit
rpc ShipOrder(ShipOrderRequest) returns (ShipOrderResponse);
4. Order Service → Notification Service (Send email/SMS/Telegram after order)
✅ gRPC is optional — use Kafka or gRPC depending on real-time needs:

If you want guaranteed delivery, async: Kafka

If you want to push instantly, fire-and-forget: gRPC

5. Admin Panels or Dashboards (e.g. getting real-time product, inventory, user info)
✅ Use gRPC between UI Gateway API ↔ microservices for real-time queries.

Reason: Faster than REST, type-safe, and less overhead.

Summary
Communication	Recommended Protocol
User/Product Created/Updated	Kafka (already done ✅)
Order → Inventory (stock check)	gRPC ✅
Order → Payment (charge card)	gRPC ✅
Order → Shipping (dispatch)	gRPC ✅
Order → Notification	Kafka OR gRPC
Admin Dashboard ↔ Services	gRPC ✅


High-Level Flow Summary
css
Copy
Edit
      [User]      [Product]
         │            │
     ┌───┴───┐    ┌───┴───┐
     │ Kafka │◄──▶│ Events│ (User/Product created, updated, deleted)
     └───────┘    └───────┘
         ▼            ▼
     [Order Service] — maintains local cache of users/products

         ▼
[Check inventory, create order]
         ▼
 [Inventory Service] — via gRPC (check & reserve)

         ▼
 [Payment Service] — via gRPC (process payment)

         ▼
 [Inventory Service] — via gRPC (confirm/reduce stock)

         ▼
 [Shipping Service] — via gRPC (schedule shipping)

         ▼
 [Notification Service] — via gRPC or Kafka (send emails, SMS, etc.)
🧭 Detailed Step-by-Step Flow
✅ 1. User & Product Services → Kafka → Order Service
When a user or product is created/updated/deleted:

They publish events to Kafka:

Topic: user-events, product-events

Order service subscribes to those and updates its local DB cache.

✅ 2. User Places an Order
The user sends a PlaceOrder request (e.g., via API Gateway).

Order Service:

Verifies user and product info from its local copy

Calls Inventory Service via gRPC to check if products are available:

go
Copy
Edit
rpc CheckProductStock(productID) returns (CheckStockResponse)
✅ 3. Inventory Service
If enough stock is available:

It reserves the requested quantity (e.g., updates reserved column).

Responds OK.

✅ 4. Payment Service
Order service calls:

go
Copy
Edit
rpc ProcessPayment(orderID, amount) returns (PaymentResponse)
On success:

Status is updated to "Paid"

On failure:

Inventory Service is called again to release reservation

✅ 5. Inventory Update
After successful payment:

Order service calls:

go
Copy
Edit
rpc ConfirmStockReduction(productID, qty)
Inventory updates:

quantity = quantity - qty

reserved = reserved - qty

✅ 6. Shipping
Order service calls:

go
Copy
Edit
rpc ShipOrder(orderID, address)
Shipping service prepares delivery

Updates shipping status

✅ 7. Notification
Order service or shipping service triggers notifications:

Via gRPC:

go
Copy
Edit
rpc SendNotification(userID, message)
Or via Kafka (e.g., to topic notification-events)

*/
