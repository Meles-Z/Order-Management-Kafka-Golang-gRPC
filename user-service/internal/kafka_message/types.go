// internal/kafka_message/types.go
package kafkamessage

type EventType string

const (
	EventUserCreated EventType = "USER_CREATED"
)

type KafkaMessage struct {
	EventType EventType   `json:"eventType"`
	Payload   interface{} `json:"payload"`
}

type UserCreatedPayload struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}
