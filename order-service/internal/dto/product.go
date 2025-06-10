package dto

import "encoding/json"

// dto/product_event.go or similar

type ProductEvent struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"` // keep it generic for now
}

