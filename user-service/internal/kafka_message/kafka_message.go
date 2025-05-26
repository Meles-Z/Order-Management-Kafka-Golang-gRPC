package kafkamessage

import (
	"encoding/json"
	"log"
)

func HandleUserEvents(msg KafkaMessage) {
	switch msg.EventType {
	case EventUserCreated:
		// Deserialize payload to specific type
		var payload UserCreatedPayload
		mapPayload(msg.Payload, &payload)

		// Example action:
		log.Printf("[Kafka] New user created: %+v\n", payload)
		// → Send email, save audit log, etc.
	default:
		log.Printf("[Kafka] Unknown event type: %s\n", msg.EventType)
	}
}

// mapPayload is a workaround to convert interface{} to struct
func mapPayload(input interface{}, out interface{}) {
	bytes, _ := json.Marshal(input)
	json.Unmarshal(bytes, out)
}
